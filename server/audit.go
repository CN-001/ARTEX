package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Aurox/artex/db"
)

// 审计日志：记录「谁、从哪、做了什么、结果如何」。
//
// 覆盖范围与取舍：
//   - 所有写操作（非 GET）自动记录，含路径/状态码/来源 IP —— 试图改配置、建任务、
//     执行自定义工具都会留痕；
//   - 登录成功/失败、登出、改密、范围门禁拦截等安全事件额外用 auditEvent 记录；
//   - **不记录请求/响应体**：里面常有密码、API key、Cookie，抄进审计表等于制造
//     第二个泄露面。detail 只放短摘要。
//   - 审计写入失败只打日志、绝不影响业务请求（审计不能成为可用性单点）。

const (
	auditDetailMax = 500
	auditPathMax   = 300
	auditTargetMax = 200
)

// auditTrail 包装 API Handler，把写操作与关键读操作写入 audit_log。
func (s *Server) auditTrail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(rec, r)
		if s.m == nil || s.m.pg == nil || !shouldAudit(r) {
			return
		}
		actor := auditActor(r, s.jwtKey)
		ip := clientIP(r)
		path := truncateStr(r.URL.Path, auditPathMax)
		action := auditActionFor(r.Method, path)
		entry := db.AuditEntry{
			Actor:  actor,
			IP:     ip,
			Method: r.Method,
			Path:   path,
			Action: action,
			Target: truncateStr(r.URL.Query().Get("id"), auditTargetMax),
			Status: rec.status,
		}
		s.insertAudit(entry)
	})
}

// statusRecorder 记录响应状态码（审计需要「结果如何」）。
type statusRecorder struct {
	http.ResponseWriter
	status  int
	written bool
}

func (s *statusRecorder) WriteHeader(code int) {
	if !s.written {
		s.status = code
		s.written = true
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	s.written = true
	return s.ResponseWriter.Write(b)
}

// Flush 让 SSE 等流式响应在包装后仍可用。
func (s *statusRecorder) Flush() {
	if f, ok := s.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// shouldAudit 决定哪些请求进审计表：
//   - 所有非 GET 的 /api 写操作；
//   - 少数安全相关 GET（审计查询本身不记，避免自我循环）。
func shouldAudit(r *http.Request) bool {
	p := r.URL.Path
	if p == "/api/health" || p == "/api/security/audit" {
		return false
	}
	if !strings.HasPrefix(p, "/api/") {
		return false
	}
	// /api/auth/login 与 /api/auth/logout 由 auth.go 里的显式 auditEvent 记录
	// （带失败原因/锁定状态等上下文），这里跳过以免一条动作两条记录。
	switch p {
	case "/api/auth/login", "/api/auth/logout":
		return false
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
		return true
	}
	switch p {
	case "/api/auth/init":
		return true
	}
	return false
}

// auditActionFor 把「方法 + 路径」压成便于检索的 action（如 settings.update）。
func auditActionFor(method, path string) string {
	p := strings.TrimPrefix(path, "/api/")
	p = strings.TrimSuffix(p, "/")
	// 去掉纯数字 id 段，让同一类操作归并（/tasks/12/delete → tasks.delete）。
	parts := strings.Split(p, "/")
	kept := make([]string, 0, len(parts))
	for _, seg := range parts {
		if seg == "" {
			continue
		}
		if _, err := strconv.ParseInt(seg, 10, 64); err == nil {
			continue
		}
		kept = append(kept, seg)
	}
	base := strings.Join(kept, ".")
	if base == "" {
		base = "api"
	}
	switch method {
	case http.MethodPost:
		return base + ".create"
	case http.MethodPut, http.MethodPatch:
		return base + ".update"
	case http.MethodDelete:
		return base + ".delete"
	}
	return base + "." + strings.ToLower(method)
}

// auditActor 从请求里解析操作者；未认证/无法解析时记 anon。
func auditActor(r *http.Request, key []byte) string {
	if tok := extractToken(r); tok != "" {
		if ok, _ := verifyJWTWithVer(tok, key, curTokenVer()); ok {
			return "ARTEX"
		}
	}
	return "anon"
}

// insertAudit 异步写库，失败只记日志；每 200 条顺手做一次保留期清理。
func (s *Server) insertAudit(e db.AuditEntry) {
	pg := s.m.pg
	if pg == nil {
		return
	}
	e.Detail = truncateStr(e.Detail, auditDetailMax)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := pg.InsertAuditCtx(ctx, e); err != nil {
			log.Printf("[audit] 写入失败: %v", err)
			return
		}
		if atomic.AddInt64(&s.auditWrites, 1)%200 == 0 {
			if n, err := pg.PruneAudit(180 * 24 * time.Hour); err == nil && n > 0 {
				log.Printf("[audit] 清理过期审计记录 %d 条", n)
			}
		}
	}()
}

// auditEvent 记录一条安全事件（登录/登出/改密/门禁拦截等）。
func (s *Server) auditEvent(r *http.Request, action, target, detail string, status int) {
	if s.m == nil || s.m.pg == nil {
		return
	}
	s.insertAudit(db.AuditEntry{
		Actor:  auditActor(r, s.jwtKey),
		IP:     clientIP(r),
		Method: r.Method,
		Path:   truncateStr(r.URL.Path, auditPathMax),
		Action: action,
		Target: truncateStr(target, auditTargetMax),
		Status: status,
		Detail: detail,
	})
	// 关键安全事件同时推送到 IM/邮件渠道（异步 + 节流，见 security_alert.go）。
	s.securityAlert(r, action, target, detail)
}

// GET /api/audit?limit=100&before_id=… —— 审计查询（受 requireAuth 保护）。
func (s *Server) auditList(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	before, _ := strconv.ParseInt(r.URL.Query().Get("before_id"), 10, 64)
	entries, err := pg.ListAudit(limit, before)
	if err != nil {
		writeErr(w, 500, "读取审计日志失败")
		return
	}
	writeJSON(w, 200, map[string]any{"entries": entries})
}

func truncateStr(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 || len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
