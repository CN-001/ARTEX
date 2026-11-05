package server

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/Aurox/artex/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtKeyFilename = "jwt.key"
	authPassKey    = "auth.password_hash"
	tokenVerKey    = "auth.token_ver"
	jwtTTL         = 7 * 24 * time.Hour
	keyChars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// SSE 短时一次性流凭证：60s 过期，清理器同频清扫。
	streamTokenTTL        = 60 * time.Second
	streamTokenCleanupTick = 60 * time.Second

	// 下限与 setup 页的前端校验一致——校验只放在前端等于没放，直接打 API 就能
	// 绕过。上限是 bcrypt 的硬限制：超过 72 字节 GenerateFromPassword 会返回
	// ErrPasswordTooLong，提前挡掉好过让用户收到一句含义不明的「密码加密失败」。
	minPasswordRunes = 8
	maxPasswordBytes = 72
)

// ---------- token 版本号（服务端可撤销） ----------
//
// tokenVer 是进程内 token 版本号。bumpTokenVersion 让 DB 与内存同时 +1，所有已签发
// token 的 ver 声明与之不符即被 requireAuth 判为「会话已失效」，实现全量吊销（含
// 已建立的 SSE 旧凭证）。0 表示尚未从 DB 载入或 DB 缺失，按 1 处理（与 DB 默认一致）。
var tokenVer atomic.Int64

// tokenVerMu 串行化版本号自增，避免并发 bump 时丢步（单管理员场景下仅防御性）。
var tokenVerMu sync.Mutex

// errTokenVerMismatch 表示 token 签名/过期都通过，但版本号已过期（全局吊销后旧 token）。
var errTokenVerMismatch = errors.New("token version mismatch")

// curTokenVer 返回当前生效的 token 版本号；未载入/缺失时按 1 处理。
func curTokenVer() int64 {
	if v := tokenVer.Load(); v > 0 {
		return v
	}
	return 1
}

// loadTokenVer 启动时从 DB 载入持久化的 token 版本号；DB 不可用或键缺失时保持默认 1。
func (s *Server) loadTokenVer() {
	if s == nil || s.m == nil || s.m.pg == nil {
		return
	}
	v, ok, err := s.m.pg.GetSetting(tokenVerKey)
	if err != nil || !ok {
		return
	}
	n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil || n < 1 {
		return
	}
	tokenVer.Store(n)
}

// bumpTokenVersion 使 DB 与内存同时 +1，返回新版本号。内存先 +1 保证本进程立刻吊销；
// DB 持久化失败时返回错误（此时重启后版本可能回退），调用方自行记录。
func (s *Server) bumpTokenVersion() (int64, error) {
	tokenVerMu.Lock()
	defer tokenVerMu.Unlock()
	newVer := curTokenVer() + 1
	tokenVer.Store(newVer)
	if s != nil && s.m != nil && s.m.pg != nil {
		if err := s.m.pg.SetSetting(tokenVerKey, strconv.FormatInt(newVer, 10)); err != nil {
			return newVer, err
		}
	}
	return newVer, nil
}

// errDataSourceUnavailable 是密码相关读操作失败时统一的回复。这些 handler 绝不能
// 把"读不到"当成"没有设置"：authInit 曾因此在数据库报错时放行，让未认证请求覆盖
// 掉已有的管理员密码。
const errDataSourceUnavailable = "数据源暂时不可用，请稍后重试"

// validatePassword 返回空串表示通过，否则返回可直接展示给用户的中文原因。
func validatePassword(pw string) string {
	if utf8.RuneCountInString(pw) < minPasswordRunes {
		return fmt.Sprintf("密码长度至少 %d 位", minPasswordRunes)
	}
	if len(pw) > maxPasswordBytes {
		return fmt.Sprintf("密码长度不能超过 %d 字节", maxPasswordBytes)
	}
	return ""
}

// loadOrCreateJWTKey reads the 32-byte signing key from keyDir/jwt.key. keyDir is
// the project base dir (next to the executable), NOT the browsable workspace root
// (dataDir) — the signing key must never be listable/downloadable via the file
// manager. Legacy installs kept it at dataDir/jwt.key; if present there and not yet
// at the new location, it is migrated (key preserved, so sessions stay valid) and
// the old file removed so it disappears from the workspace. On first run a random
// key is generated and persisted.
func loadOrCreateJWTKey(keyDir, dataDir string) ([]byte, error) {
	path := filepath.Join(keyDir, jwtKeyFilename)
	// one-time migration out of the old in-workspace location.
	if legacy := filepath.Join(dataDir, jwtKeyFilename); legacy != path {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if data, rerr := os.ReadFile(legacy); rerr == nil {
				if werr := os.WriteFile(path, data, 0o600); werr == nil {
					_ = os.Remove(legacy)
					log.Printf("[auth] JWT key 已从 %s 迁移到 %s（移出可浏览工作区）", legacy, path)
				}
			}
		}
	}
	if data, err := os.ReadFile(path); err == nil && len(strings.TrimSpace(string(data))) >= 32 {
		return []byte(strings.TrimSpace(string(data))), nil
	}
	buf := make([]byte, 32)
	for i := range buf {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(keyChars))))
		if err != nil {
			return nil, fmt.Errorf("generate jwt key: %w", err)
		}
		buf[i] = keyChars[n.Int64()]
	}
	if err := os.WriteFile(path, buf, 0600); err != nil {
		return nil, fmt.Errorf("write jwt key: %w", err)
	}
	log.Printf("[auth] 新 JWT key 已写入 %s", path)
	return buf, nil
}

// signJWT issues a 7-day HS256 token for user ARTEX, signed with the current
// token version (see signJWTWithVer).
func signJWT(key []byte) (string, error) {
	return signJWTWithVer(key, curTokenVer())
}

// signJWTWithVer 签名带自定义声明 ver 的 token，供 bumpTokenVersion 之后的全量吊销使用。
func signJWTWithVer(key []byte, ver int64) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "ARTEX",
		"exp": now.Add(jwtTTL).Unix(),
		"iat": now.Unix(),
		"ver": ver,
	}).SignedString(key)
}

// verifyJWT returns true when tokenStr is a valid, non-expired HS256 token.
// 仅校验签名/过期，不校验版本号——保留给既有测试与无需版本判断的调用点。
func verifyJWT(tokenStr string, key []byte) bool {
	ok, _ := verifyJWTWithVer(tokenStr, key, -1)
	return ok
}

// verifyJWTWithVer 校验签名+过期+版本号。版本号不符返回 (false, errTokenVerMismatch)；
// wantVer < 0 时跳过版本号校验（供 verifyJWT 复用）。
func verifyJWTWithVer(tokenStr string, key []byte, wantVer int64) (bool, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return key, nil
	})
	if err != nil || !t.Valid {
		return false, err
	}
	if wantVer < 0 {
		return true, nil
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}
	ver, ok := claims["ver"].(float64)
	if !ok || int64(ver) != wantVer {
		return false, errTokenVerMismatch
	}
	return true, nil
}

// extractToken reads the JWT from Authorization: Bearer header or the artex_token
// HttpOnly cookie.
//
// 注意：不再接受 ?token= 查询参数。此前 SSE 把长期 JWT 拼在 URL 上，会落进访问
// 日志/代理日志/浏览器历史；SSE 现在必须用一次性流凭证（?stream=，见
// consumeStreamToken），前端也已切换为纯 Cookie 认证。
func extractToken(r *http.Request) string {
	if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	if c, err := r.Cookie("artex_token"); err == nil && c.Value != "" {
		return c.Value
	}
	return ""
}

// requireAuth wraps h with JWT validation.
// /api/auth/*（除 logout / stream-token）与 /api/health 豁免；logout 与 stream-token
// 必须经过鉴权，否则任何人都能吊销全部会话或领取流凭证。
func (s *Server) requireAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		exempt := p == "/api/health" ||
			(strings.HasPrefix(p, "/api/auth/") &&
				p != "/api/auth/logout" &&
				p != "/api/auth/stream-token" &&
				p != "/api/auth/me")
		if exempt {
			h.ServeHTTP(w, r)
			return
		}
		// SSE 短时一次性流凭证：EventSource 无法带 Authorization 头，改用 ?stream=
		// 下发的一次性凭证（验签 + 未过期 + 版本匹配 + nonce 一次性消费）。
		if st := r.URL.Query().Get("stream"); st != "" {
			if !s.consumeStreamToken(st) {
				writeErr(w, 401, "流凭证无效或已使用")
				return
			}
			h.ServeHTTP(w, r)
			return
		}
		tok := extractToken(r)
		if tok == "" {
			writeErr(w, 401, "未授权")
			return
		}
		ok, err := verifyJWTWithVer(tok, s.jwtKey, curTokenVer())
		if errors.Is(err, errTokenVerMismatch) {
			writeErr(w, 401, "会话已失效，请重新登录")
			return
		}
		if !ok {
			writeErr(w, 401, "token 无效或已过期")
			return
		}
		// 写操作防 CSRF：仅当认证来自 cookie（无 Authorization 头）时才校验 Origin/Referer。
		// Authorization 头是跨站无法自动携带的显式凭证，天然免疫 CSRF；cookie 会被浏览器
		// 自动带上，恶意站点可利用其发起写请求，故需校验来源。无 Origin/Referer（curl/CLI）
		// 一律放行，避免误伤自动化调用——这类调用本就不在浏览器 CSRF 威胁模型内。
		if isWriteMethod(r.Method) && !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			if err := checkOriginReferer(r); err != nil {
				writeErr(w, 403, "来源校验失败")
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// ---------- 会话 cookie（HttpOnly） ----------

// requestIsTLS 判断请求是否经 TLS（直连 r.TLS 或可信反代 X-Forwarded-Proto: https）。
func requestIsTLS(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

// setSessionCookie 下发 HttpOnly 会话 cookie，内容与返回的 JWT 一致。Secure 仅在请求
// 经 TLS 时设置，避免纯 HTTP 内网部署下因 Secure 导致浏览器根本不回传 cookie。
func setSessionCookie(w http.ResponseWriter, r *http.Request, tok string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "artex_token",
		Value:    tok,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(jwtTTL.Seconds()),
		Secure:   requestIsTLS(r),
	})
}

// clearSessionCookie 清除会话 cookie（登出用）。
func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "artex_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Secure:   requestIsTLS(r),
	})
}

// ---------- 写操作 Origin/Referer 校验（防 CSRF） ----------

// isWriteMethod 判断是否属于可改变状态的写方法（CSRF 防护对象）。
func isWriteMethod(m string) bool {
	switch m {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

// checkOriginReferer 校验 Origin/Referer 的 host 是否等于 r.Host（允许 http/https 协议不同）。
// 无 Origin/Referer 或无法解析 host 时放行——宁可漏检也不误伤 CLI/自动化调用。
func checkOriginReferer(r *http.Request) error {
	src := r.Header.Get("Origin")
	if src == "" {
		src = r.Header.Get("Referer")
	}
	if src == "" {
		return nil
	}
	u, err := url.Parse(src)
	if err != nil || u.Host == "" {
		return nil
	}
	if !strings.EqualFold(u.Host, r.Host) {
		return errors.New("origin/referer host mismatch")
	}
	return nil
}

// ---------- SSE 短时一次性流凭证 ----------
//
// EventSource 不能带 Authorization 头，把长寿命 JWT 直接拼进 URL 会在日志/反代缓存里
// 泄露。改为下发 60s 一次性流凭证：v1.<base64url(payload)>.<hex hmac>，payload =
// <ver>:<expUnix>:<nonceHex16>。nonce 进程内保存、一次性消费（消费即删）。

type streamNonceStore struct {
	mu     sync.Mutex
	nonces map[string]int64 // nonceHex16 -> expUnix
}

var (
	streamNonces          = &streamNonceStore{nonces: make(map[string]int64)}
	streamNonceCleanupOn  sync.Once
)

func (st *streamNonceStore) add(nonce string, exp int64) {
	st.mu.Lock()
	st.nonces[nonce] = exp
	st.mu.Unlock()
}

func (st *streamNonceStore) consume(nonce string) bool {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.nonces[nonce]; !ok {
		return false
	}
	delete(st.nonces, nonce)
	return true
}

func (st *streamNonceStore) ensureCleanup() {
	streamNonceCleanupOn.Do(func() {
		go func() {
			t := time.NewTicker(streamTokenCleanupTick)
			defer t.Stop()
			for range t.C {
				st.sweep(time.Now().Unix())
			}
		}()
	})
}

func (st *streamNonceStore) sweep(now int64) {
	st.mu.Lock()
	defer st.mu.Unlock()
	for n, exp := range st.nonces {
		if exp <= now {
			delete(st.nonces, n)
		}
	}
}

// issueStreamToken 签发一次性流凭证。
func (s *Server) issueStreamToken() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(b)
	exp := time.Now().Add(streamTokenTTL).Unix()
	payload := fmt.Sprintf("%d:%d:%s", curTokenVer(), exp, nonce)
	mac := hmac.New(sha256.New, s.jwtKey)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	streamNonces.add(nonce, exp)
	return "v1." + base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig, nil
}

// consumeStreamToken 校验并一次性消费流凭证；成功仅一次。
func (s *Server) consumeStreamToken(tok string) bool {
	parts := strings.Split(tok, ".")
	if len(parts) != 3 || parts[0] != "v1" {
		return false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	payload := string(raw)
	mac := hmac.New(sha256.New, s.jwtKey)
	mac.Write([]byte(payload))
	if !hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(parts[2])) {
		return false
	}
	fields := strings.Split(payload, ":")
	if len(fields) != 3 {
		return false
	}
	ver, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil || ver != curTokenVer() {
		return false
	}
	exp, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return false
	}
	return streamNonces.consume(fields[2])
}

// GET /api/auth/status — reports whether the admin password has been initialised.
// 读失败必须回 503 而不是 initialized:false：前端在 initialized:false 时会把用户
// 送到 /setup 去设置密码（login/page.tsx），把数据库故障包装成 200 等于把用户往
// 覆盖已有密码的路上推。
func (s *Server) authStatus(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		return
	}
	hash, _, err := pg.GetSetting(authPassKey)
	if err != nil {
		writeErr(w, 503, errDataSourceUnavailable)
		return
	}
	writeJSON(w, 200, map[string]any{"initialized": hash != ""})
}

// GET /api/auth/me —— 返回当前会话对应的用户（受 requireAuth 保护）。
// 纯 Cookie 认证后前端不再能读 JWT 解析用户名，改由该接口获取；这也是前端
// 判断「会话是否仍然有效」的探针（401 即跳登录）。
func (s *Server) authMe(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"username": "ARTEX",
		"name":     "ARTEX",
		"role":     "admin",
	})
}

// POST /api/auth/init — sets the password for the first time; rejected if already set.
func (s *Server) authInit(w http.ResponseWriter, r *http.Request) {
	// 该路由在 requireAuth 里豁免（首次初始化时尚无会话），因此 CSRF 来源校验在此补：
	// 若浏览器从异源站点提交，直接拒绝，避免"未初始化窗口"被跨站抢占初始密码。
	// 无 Origin/Referer 的 CLI 调用放行（与其它写操作一致的取舍）。
	if err := checkOriginReferer(r); err != nil {
		writeErr(w, 403, "来源校验失败")
		return
	}
	pg := s.pg(w)
	if pg == nil {
		return
	}
	existing, _, err := pg.GetSetting(authPassKey)
	if err != nil {
		writeErr(w, 503, errDataSourceUnavailable)
		return
	}
	if existing != "" {
		writeErr(w, 403, "密码已设置")
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := decode(r, &req); err != nil || req.Password == "" {
		writeErr(w, 400, "密码不能为空")
		return
	}
	if msg := validatePassword(req.Password); msg != "" {
		writeErr(w, 400, msg)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeErr(w, 500, "密码加密失败")
		return
	}
	// 用 INSERT ... ON CONFLICT DO NOTHING 而不是 upsert：上面那次 GetSetting 只是
	// 快速失败路径，真正"仅首次可设"的保证落在主键约束上。bcrypt 要跑几十毫秒，
	// 这期间别的请求完全可能先把密码设好，而读检查本身也可能因故障而失效。
	inserted, err := pg.InsertSettingIfAbsent(authPassKey, string(hash))
	if err != nil {
		writeErr(w, 500, "保存失败: "+err.Error())
		return
	}
	if !inserted {
		writeErr(w, 403, "密码已设置")
		return
	}
	tok, err := signJWT(s.jwtKey)
	if err != nil {
		writeErr(w, 500, "token 生成失败")
		return
	}
	// 首次初始化后直接下发 HttpOnly 会话 cookie（纯 Cookie 认证下前端拿不到 token，
	// 靠这次 Set-Cookie 完成"设完密码即登录"）。
	setSessionCookie(w, r, tok)
	s.auditEvent(r, "auth.init", "ARTEX", "首次初始化口令并建立会话", 200)
	writeJSON(w, 200, map[string]any{"token": tok})
}

// POST /api/auth/change-password — changes the admin password. Requires a valid
// token (this route is under /api/auth/* which requireAuth exempts, so the token
// is validated here) AND the current password.
func (s *Server) authChangePassword(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		return
	}
	if ok, _ := verifyJWTWithVer(extractToken(r), s.jwtKey, curTokenVer()); !ok {
		writeErr(w, 401, "未授权")
		return
	}
	// 该路由同样在 requireAuth 里被豁免（需要在未登录时也能改密），CSRF 来源校验在此补：
	// cookie 认证 + 异源 Origin/Referer → 拒绝。改密还需要旧密码，本不可被 CSRF 利用，
	// 这里属纵深防御，避免将来去掉旧密码校验时留下缺口。
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		if err := checkOriginReferer(r); err != nil {
			writeErr(w, 403, "来源校验失败")
			return
		}
	}
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := decode(r, &req); err != nil {
		writeErr(w, 400, "请求格式错误")
		return
	}
	if req.NewPassword == "" {
		writeErr(w, 400, "新密码不能为空")
		return
	}
	if msg := validatePassword(req.NewPassword); msg != "" {
		writeErr(w, 400, msg)
		return
	}
	hash, ok, err := pg.GetSetting(authPassKey)
	if err != nil {
		writeErr(w, 503, errDataSourceUnavailable)
		return
	}
	if !ok || hash == "" {
		writeErr(w, 403, "密码未初始化，请先设置密码")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.OldPassword)); err != nil {
		writeErr(w, 401, "当前密码错误")
		return
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeErr(w, 500, "密码加密失败")
		return
	}
	s.auditEvent(r, "auth.password_change", "ARTEX", "口令已修改（旧会话全部吊销）", 200)
	if err := pg.SetSetting(authPassKey, string(newHash)); err != nil {
		writeErr(w, 500, "保存失败: "+err.Error())
		return
	}
	// 改密后吊销所有已签发 token（旧会话全部失效），包括本次改密使用的会话本身。
	if _, err := s.bumpTokenVersion(); err != nil {
		log.Printf("[auth] change_password bump_token_ver ip=%q user=%q err=%v", clientIP(r), "ARTEX", err)
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}

// ---------- 登录限速 ----------
//
// 进程内内存限速器：键 = 客户端 IP + 规范化用户名。目的：给爆破 /api/auth/login 的
// 攻击者加一道进程级速率闸门，避免 bcrypt（故意慢）被恶意高频调用拖垮 CPU。
// 内存 map + 互斥锁，重启即清零，对单机部署足够，不引入外部依赖。
//
// 规则：同一键 15 分钟窗口内失败 ≥5 次即锁定；锁定时长按阶梯 60s/120s/300s，封顶
// 15 分钟。成功登录清零。失败与锁定各写一条结构化日志（不含密码、不含 token）。

const (
	loginWindow      = 15 * time.Minute
	loginMaxFails    = 5
	loginCleanupTick = 10 * time.Minute
)

// loginLockLadder 是锁定时长阶梯：第 1 次锁 60s，其后 120s、300s，封顶 900s(15 分钟)。
var loginLockLadder = []time.Duration{
	60 * time.Second,
	120 * time.Second,
	300 * time.Second,
	15 * time.Minute,
}

type loginAttempt struct {
	fails       int
	windowStart time.Time
	lockLevel   int // 0=首次锁定；每次锁定 +1，封顶 len(loginLockLadder)-1
	lockedUntil time.Time
}

type loginLimiter struct {
	mu       sync.Mutex
	attempts map[string]*loginAttempt
}

var (
	loginLim            = &loginLimiter{attempts: make(map[string]*loginAttempt)}
	loginLimCleanupOnce sync.Once
)

// ensureCleanup 启动一次后台清理：删除 15 分钟窗口已过期且不在锁定中的条目，
// 防止长期运行后 map 无界增长。
func (l *loginLimiter) ensureCleanup() {
	loginLimCleanupOnce.Do(func() {
		go func() {
			t := time.NewTicker(loginCleanupTick)
			defer t.Stop()
			for range t.C {
				l.sweep(time.Now())
			}
		}()
	})
}

func (l *loginLimiter) sweep(now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for k, a := range l.attempts {
		if now.Sub(a.windowStart) >= loginWindow && !a.lockedUntil.After(now) {
			delete(l.attempts, k)
		}
	}
}

// loginKey 规范化用户名（去空白、小写）后与客户端 IP 组合，作为限速键。
func loginKey(ip, username string) string {
	return ip + "|" + strings.ToLower(strings.TrimSpace(username))
}

// allow 返回该键当前是否被锁定；锁定则给出剩余秒数（≥1）。
func (l *loginLimiter) allow(key string) (locked bool, retryAfterSecs int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	a := l.attempts[key]
	if a == nil {
		return false, 0
	}
	now := time.Now()
	if a.lockedUntil.After(now) {
		secs := int(a.lockedUntil.Sub(now).Seconds()) + 1
		if secs < 1 {
			secs = 1
		}
		return true, secs
	}
	return false, 0
}

// recordFailure 记录一次失败；达到阈值则加锁并返回锁定时长。
func (l *loginLimiter) recordFailure(key string) (locked bool, retryAfterSecs int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	a := l.attempts[key]
	if a == nil || now.Sub(a.windowStart) >= loginWindow {
		a = &loginAttempt{windowStart: now}
		l.attempts[key] = a
	}
	a.fails++
	if a.fails >= loginMaxFails {
		idx := a.lockLevel
		if idx >= len(loginLockLadder) {
			idx = len(loginLockLadder) - 1
		}
		dur := loginLockLadder[idx]
		a.lockedUntil = now.Add(dur)
		if a.lockLevel < len(loginLockLadder)-1 {
			a.lockLevel++
		}
		// 加锁后重置失败窗口：锁到期后需再累积 5 次失败才进入下一档时长。
		a.fails = 0
		a.windowStart = now
		return true, int(dur.Seconds())
	}
	return false, 0
}

// clear 成功登录后清零该键的全部状态。
func (l *loginLimiter) clear(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
}

// clientIP 取直连层地址。不要默认信任 X-Forwarded-For：该头可被客户端伪造。
// TODO: 只有部署在可信反代后（反代会覆盖/追加该头）才应按白名单信任其最左一跳。
func clientIP(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// loginFail 记录一次失败；若因此触发锁定则回 429，否则回 401。日志不含密码/token。
func (s *Server) loginFail(w http.ResponseWriter, r *http.Request, ip, key, user, reason string) {
	if locked, retry := loginLim.recordFailure(key); locked {
		log.Printf("[auth] login_locked ip=%q user=%q retry_after=%ds", ip, user, retry)
		// 锁定也是安全事件，进审计（便于发现撞库/爆破）。
		s.auditEvent(r, "auth.login.locked", user, "失败次数过多，已临时锁定", 429)
		w.Header().Set("Retry-After", strconv.Itoa(retry))
		writeErr(w, 429, "尝试过于频繁，请稍后再试")
		return
	}
	log.Printf("[auth] login_failed ip=%q user=%q reason=%s", ip, user, reason)
	s.auditEvent(r, "auth.login.failed", user, reason, 401)
	writeErr(w, 401, "用户名或密码错误")
}

// POST /api/auth/login — validates username/password and returns a JWT.
func (s *Server) authLogin(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decode(r, &req); err != nil {
		writeErr(w, 400, "请求格式错误")
		return
	}
	// 限速在解出用户名之后、bcrypt 之前：已锁定的键直接 429，不消耗 bcrypt。
	loginLim.ensureCleanup()
	ip := clientIP(r)
	key := loginKey(ip, req.Username)
	user := strings.ToLower(strings.TrimSpace(req.Username))
	if locked, retry := loginLim.allow(key); locked {
		log.Printf("[auth] login_locked ip=%q user=%q retry_after=%ds", ip, user, retry)
		w.Header().Set("Retry-After", strconv.Itoa(retry))
		writeErr(w, 429, "尝试过于频繁，请稍后再试")
		return
	}
	if req.Username != "ARTEX" {
		s.loginFail(w, r, ip, key, user, "bad_username")
		return
	}
	hash, ok, err := pg.GetSetting(authPassKey)
	if err != nil {
		writeErr(w, 503, errDataSourceUnavailable)
		return
	}
	if !ok || hash == "" {
		writeErr(w, 403, "密码未初始化，请先设置密码")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		s.loginFail(w, r, ip, key, user, "bad_password")
		return
	}
	loginLim.clear(key)
	log.Printf("[auth] login_ok ip=%q user=%q", ip, user)
	tok, err := signJWT(s.jwtKey)
	if err != nil {
		writeErr(w, 500, "token 生成失败")
		return
	}
	setSessionCookie(w, r, tok)
	s.auditEvent(r, "auth.login", "ARTEX", "登录成功", 200)
	writeJSON(w, 200, map[string]any{"token": tok})
}

// POST /api/auth/logout — 吊销所有已签发 token（token 版本 +1）并清除会话 cookie。
// 受 requireAuth 保护（未加入豁免白名单），必须持有有效会话才能触发吊销。
func (s *Server) authLogout(w http.ResponseWriter, r *http.Request) {
	// 该路由在 requireAuth 里被豁免（否则未登录无法清理 cookie），因此 CSRF 来源
	// 校验要在这里补：cookie 认证的写操作若带异源 Origin/Referer 直接拒绝，
	// 防止第三方站点把用户"静默登出"（DoS 类）。Authorization 头调用不受影响。
	if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		if err := checkOriginReferer(r); err != nil {
			writeErr(w, 403, "来源校验失败")
			return
		}
	}
	ip := clientIP(r)
	newVer, err := s.bumpTokenVersion()
	if err != nil {
		log.Printf("[auth] logout ip=%q user=%q token_ver=%d db_err=%v", ip, "ARTEX", newVer, err)
	} else {
		log.Printf("[auth] logout ip=%q user=%q token_ver=%d", ip, "ARTEX", newVer)
	}
	s.auditEvent(r, "auth.logout", "ARTEX", "会话吊销 + 清 cookie", 200)
	clearSessionCookie(w, r)
	writeJSON(w, 200, map[string]any{"ok": true})
}

// POST /api/auth/stream-token — 签发一次性 SSE 流凭证。受 requireAuth 保护。
func (s *Server) authStreamToken(w http.ResponseWriter, r *http.Request) {
	streamNonces.ensureCleanup()
	tok, err := s.issueStreamToken()
	if err != nil {
		writeErr(w, 500, "流凭证生成失败")
		return
	}
	writeJSON(w, 200, map[string]any{"token": tok, "expires_in": 60})
}

// GET /api/settings/security — 安全状态总览（受 requireAuth 保护）。绝不返回密钥/密文/明文。
func (s *Server) settingsSecurity(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{
		"api_key_encryption":    db.APIKeyEncryptionState(),
		"token_version":         curTokenVer(),
		"sse_token_ttl_seconds": 60,
		"cookie": map[string]any{
			"http_only": true,
			"same_site": "Strict",
			"secure":    "auto",
		},
		"cors_configured": strings.TrimSpace(os.Getenv("ARTEX_CORS_ORIGINS")) != "",
	})
}
