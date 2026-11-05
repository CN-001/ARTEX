package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Aurox/artex/db"
	"github.com/Aurox/artex/scope"
	"github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
)

// 本文件实现自定义工具执行器(docs/自定义工具设计.md)。system=false 的 tools 行按
// kind 分派:command(渲染命令→复用 Bash 底层 run)、script(仅 Python;写临时文件、
// stdin=参数 JSON + env TOOL_*、用配置的解释器)、http(原生请求+可设代理)。这些工具
// 像流量/编排工具一样 seed 不需要(它们本就在 tools 表),经 hostTools 注入、按绑定过滤。

// ---------- 自定义工具 CRUD ----------

type customToolReq struct {
	Key         string          `json:"key"`
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema"`
	Agents      []string        `json:"agents"`
	Enabled     bool            `json:"enabled"`
	Kind        string          `json:"kind"` // command | script | http
	Exec        json.RawMessage `json:"exec"`
	Deferred    bool            `json:"deferred"`
}

var reToolKey = reAgentKey // 同 agent key 规则:小写字母开头 + 小写字母/数字/下划线

func (s *Server) pgCreateCustomTool(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		return
	}
	var req customToolReq
	if err := decode(r, &req); err != nil {
		writeErr(w, 400, err.Error())
		return
	}
	req.Key = strings.TrimSpace(req.Key)
	if !reToolKey.MatchString(req.Key) {
		writeErr(w, 400, "key 需小写字母开头，仅含小写字母/数字/下划线")
		return
	}
	if req.Kind != "command" && req.Kind != "script" && req.Kind != "http" && req.Kind != "shell" {
		writeErr(w, 400, "kind 需为 command / script / http / shell")
		return
	}
	if req.Kind == "http" && !hasSchemaProps(req.Schema) {
		writeErr(w, 400, "http 工具必须提供参数 JSON Schema(不能留空)")
		return
	}
	if exist, _ := pg.GetTool(req.Key); exist != nil {
		writeErr(w, 409, "该 key 已存在(内置或自定义工具)")
		return
	}
	if err := pg.CreateCustomTool(&db.Tool{
		Key: req.Key, Description: req.Description, Schema: req.Schema, Agents: req.Agents,
		Enabled: req.Enabled, Kind: req.Kind, Exec: req.Exec, Deferred: req.Deferred,
	}); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"key": req.Key})
}

func (s *Server) pgUpdateCustomTool(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		return
	}
	key := r.PathValue("key")
	existing, err := pg.GetTool(key)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	if existing == nil || existing.System {
		writeErr(w, 400, "只能编辑自定义工具")
		return
	}
	var req customToolReq
	if err := decode(r, &req); err != nil {
		writeErr(w, 400, err.Error())
		return
	}
	if req.Kind != "command" && req.Kind != "script" && req.Kind != "http" && req.Kind != "shell" {
		writeErr(w, 400, "kind 需为 command / script / http / shell")
		return
	}
	if req.Kind == "http" && !hasSchemaProps(req.Schema) {
		writeErr(w, 400, "http 工具必须提供参数 JSON Schema(不能留空)")
		return
	}
	if err := pg.UpdateCustomTool(&db.Tool{
		Key: key, Description: req.Description, Schema: req.Schema, Agents: req.Agents,
		Enabled: req.Enabled, Kind: req.Kind, Exec: req.Exec, Deferred: req.Deferred,
	}); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) pgDeleteCustomTool(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		return
	}
	key := r.PathValue("key")
	if err := pg.DeleteCustomTool(key); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"deleted": key})
}

// testToolReq is a dry-run request from the editor: run the given (possibly unsaved)
// exec spec with sample params, without persisting the tool. Same executor path as a
// real tool call — it runs arbitrary command/script/http on the server, which the
// custom-tool feature already allows, so no new capability is granted.
type testToolReq struct {
	Kind   string          `json:"kind"` // command | script | http
	Exec   json.RawMessage `json:"exec"`
	Params map[string]any  `json:"params"`
}

// hostsInParams 尽力从命令/脚本工具的参数里提取「明显的主机目标」，用于试运行端点的
// 网络范围预检。只提取 URL、IPv4[:port] 两类高置信度写法，不做域名猜测，以免把普通
// 参数误判成目标（真正的硬门禁在 runHTTPTool / scope.Dial）。
func hostsInParams(params map[string]any) []string {
	var out []string
	seen := map[string]bool{}
	add := func(h string) {
		h = strings.TrimSpace(h)
		if h == "" || seen[h] {
			return
		}
		seen[h] = true
		out = append(out, h)
	}
	for _, v := range params {
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		for _, m := range urlHostRe.FindAllString(s, -1) {
			if u, err := url.Parse(m); err == nil && u.Host != "" {
				add(u.Host)
			}
		}
		for _, m := range ipv4Re.FindAllString(s, -1) {
			if strings.Contains(m, ":") || net.ParseIP(m) != nil {
				add(m)
			}
		}
	}
	return out
}

var (
	// urlHostRe 匹配 http(s) URL（持续到空白/引号为止）。
	urlHostRe = regexp.MustCompile(`(?i)https?://[^\s"'<>\\]+`)
	// ipv4Re 匹配 IPv4，可带 :port。
	ipv4Re = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}(?::\d{1,5})?\b`)
)

// pgTestCustomTool executes an exec spec once and returns its raw output + error
// flag, so the editor can debug a tool before saving it. Per-kind timeouts still
// apply from the exec spec (with defaults); the outer ceiling is a hard backstop.
func (s *Server) pgTestCustomTool(w http.ResponseWriter, r *http.Request) {
	pg := s.pg(w)
	if pg == nil {
		return
	}
	var req testToolReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "无效的请求体")
		return
	}
	params := req.Params
	if params == nil {
		params = map[string]any{}
	}
	// P0-1/P0-2 管控:先审计(dry-run,只记 kind 与 http 的 host,绝不记完整 exec/params,
	// 避免把 header/query/命令里的明文密钥写进日志),再做网络范围预检。
	auditTarget := ""
	if req.Kind == "http" {
		var spec httpExec
		_ = json.Unmarshal(req.Exec, &spec)
		if u, err := url.Parse(renderTemplate(spec.URL, params, identity)); err == nil {
			auditTarget = u.Host
		}
	}
	log.Printf("[custom-tool] test_exec kind=%s target=%s", req.Kind, auditTarget)

	// 网络范围硬门禁(P0-2):http 直接检查 URL 的 host;command/script 尽力从参数里
	// 提取明显主机,提取不到就跳过(真正的硬门禁在 runHTTPTool)。
	if auditTarget != "" {
		if blocked, _ := scope.Blocked(auditTarget); blocked {
			writeErr(w, 403, "目标命中平台网络范围门禁")
			return
		}
	}
	if req.Kind == "command" || req.Kind == "script" {
		for _, host := range hostsInParams(params) {
			if blocked, _ := scope.Blocked(host); blocked {
				writeErr(w, 403, "目标命中平台网络范围门禁")
				return
			}
		}
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()
	tc := &actool.ToolContext{WorkingDir: s.m.dir} // run in the project dir, like a real call
	var res actool.Result
	switch req.Kind {
	case "command":
		res, _ = s.runCommandTool(ctx, req.Exec, params, tc)
	case "script":
		res, _ = s.runScriptTool(ctx, "test", req.Exec, params, tc)
	case "http":
		res, _ = s.runHTTPTool(ctx, req.Exec, params, tc)
	case "shell":
		writeErr(w, 400, "shell 类型工具是 bash 环境声明，无可执行内容")
		return
	default:
		writeErr(w, 400, "未知工具类型: "+req.Kind)
		return
	}
	writeJSON(w, 200, map[string]any{"output": res.Flatten(), "is_error": res.IsError})
}

// ---------- Python 解释器(检测 + 入库 + 覆盖) ----------

const settingPythonInterp = "python_interpreter"

// detectPython finds a python interpreter absolute path (python3 preferred).
func detectPython() string {
	for _, c := range []string{"python3", "python"} {
		if p, err := exec.LookPath(c); err == nil {
			return p
		}
	}
	return ""
}

// pythonInterpreter resolves the interpreter: user-set > stored auto-detect > live
// detect. "" only when truly none found.
func (s *Server) pythonInterpreter() string {
	if v, ok, _ := s.m.pg.GetSetting(settingPythonInterp); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return detectPython()
}

// seedPythonInterpreter stores the auto-detected interpreter on startup if unset
// (never clobbers a user-set value).
func (s *Server) seedPythonInterpreter() {
	if v, ok, _ := s.m.pg.GetSetting(settingPythonInterp); ok && strings.TrimSpace(v) != "" {
		return
	}
	if p := detectPython(); p != "" {
		_ = s.m.pg.SetSetting(settingPythonInterp, p)
		log.Printf("[custom-tool] 自动检测到 python 解释器: %s", p)
	}
}

// ---------- exec 规格 ----------

type commandExec struct {
	Command   string `json:"command"`
	TimeoutMs int    `json:"timeout_ms"`
}
type scriptExec struct {
	Code      string `json:"code"`
	TimeoutMs int    `json:"timeout_ms"`
}
type httpExec struct {
	Method            string            `json:"method"`
	URL               string            `json:"url"`
	Headers           map[string]string `json:"headers"`
	Body              string            `json:"body"`
	TimeoutMs         int               `json:"timeout_ms"`
	Proxy             string            `json:"proxy"`
	UseRecordingProxy bool              `json:"use_recording_proxy"`
}

func timeoutOr(ms, def int) time.Duration {
	if ms <= 0 {
		return time.Duration(def) * time.Millisecond
	}
	return time.Duration(ms) * time.Millisecond
}

// ---------- 通用工具构造 ----------

// customTools builds CoreTools for every user-defined (system=false) tool row.
// shell-kind tools are environment hints only — they surface in the Bash tool
// description via ToolResolve and do NOT create callable tool entries here.
func (s *Server) customTools() ([]actool.CoreTool, error) {
	rows, err := s.m.pg.ListCustomTools()
	if err != nil {
		return nil, err
	}
	out := make([]actool.CoreTool, 0, len(rows))
	for _, t := range rows {
		if t.Kind == "shell" {
			continue // shell hints are handled by ToolResolve → Bash description
		}
		out = append(out, s.buildCustomTool(t))
	}
	return out, nil
}

// buildCustomTool turns one custom-tool row into a CoreTool. Empty schema → a thin
// {args:string} (薄壳工具), so command/http templates can use {args}.
func (s *Server) buildCustomTool(t *db.Tool) actool.CoreTool {
	schema := ensureSchema(t.Schema)
	key, kind, execRaw := t.Key, t.Kind, t.Exec
	run := func(ctx context.Context, in json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
		var params map[string]any
		if len(in) > 0 {
			_ = json.Unmarshal(in, &params)
		}
		if params == nil {
			params = map[string]any{}
		}
		switch kind {
		case "command":
			return s.runCommandTool(ctx, execRaw, params, tc)
		case "script":
			return s.runScriptTool(ctx, key, execRaw, params, tc)
		case "http":
			return s.runHTTPTool(ctx, execRaw, params, tc)
		default:
			return actool.Errorf("未知自定义工具类型: " + kind), nil
		}
	}
	return actool.Build(actool.Spec{
		Name: key, Description: t.Description, Schema: schema,
		Permissions: func(context.Context, json.RawMessage, permission.Context) permission.Decision {
			return permission.Allowed()
		},
		Run: run,
	})
}

// hasSchemaProps reports whether raw is a JSON-Schema object with ≥1 property.
// http tools require an explicit schema (the auto {args} shell can't name the
// {param} placeholders in URL/headers/body), so an empty schema is rejected.
func hasSchemaProps(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var m map[string]any
	if json.Unmarshal(raw, &m) != nil {
		return false
	}
	props, _ := m["properties"].(map[string]any)
	return len(props) > 0
}

// ensureSchema returns the tool's schema, or a thin {args:string} when none given.
func ensureSchema(raw json.RawMessage) map[string]any {
	var m map[string]any
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &m)
	}
	props, _ := m["properties"].(map[string]any)
	if len(props) > 0 {
		return m
	}
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"args": map[string]any{"type": "string", "description": "命令/参数(自由文本)"},
		},
	}
}

// ---------- command:渲染命令 → 复用 Bash 底层 run ----------

func (s *Server) runCommandTool(ctx context.Context, execRaw json.RawMessage, params map[string]any, tc *actool.ToolContext) (actool.Result, error) {
	var spec commandExec
	_ = json.Unmarshal(execRaw, &spec)
	if strings.TrimSpace(spec.Command) == "" {
		return actool.Errorf("command 为空"), nil
	}
	cmd := renderTemplate(spec.Command, params, shellQuote)
	// 注意：command 工具走 norma 的 Bash 底层，而 norma 已在 Bash 执行层统一注入沙箱
	// （third_party/norma/tool/sandbox_unix.go，ARTEX_TOOL_SANDBOX=1 时 setpriv 降权）。
	// 这里不能再包一层，否则 setpriv 套 setpriv 会因 setgroups 权限失败。
	// 复用 Bash 也在用的底层 run(经 Bash CoreTool.Call):自动继承安全 floor/超时/
	// 代理 env/输出溢出。工具与 Bash 平级、共用底层,不经过 Bash 这个工具让模型调。
	bashIn, _ := json.Marshal(map[string]any{"command": cmd})
	if spec.TimeoutMs > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeoutOr(spec.TimeoutMs, 120000))
		defer cancel()
	}
	return actool.NewBash().Call(ctx, bashIn, tc)
}

// ---------- script(仅 Python):临时文件 + stdin JSON + env ----------

func (s *Server) runScriptTool(ctx context.Context, key string, execRaw json.RawMessage, params map[string]any, tc *actool.ToolContext) (actool.Result, error) {
	var spec scriptExec
	_ = json.Unmarshal(execRaw, &spec)
	if strings.TrimSpace(spec.Code) == "" {
		return actool.Errorf("script code 为空"), nil
	}
	interp := s.pythonInterpreter()
	if interp == "" {
		return actool.Errorf("未配置且未检测到 python 解释器(在系统配置里设置)"), nil
	}
	workDir := s.m.dir
	var sessionEnv []string
	if tc != nil {
		if tc.WorkingDir != "" {
			workDir = tc.WorkingDir
		}
		sessionEnv = tc.Env
	}
	body, err := execPython(ctx, interp, key, spec.Code, params, workDir, sessionEnv, timeoutOr(spec.TimeoutMs, 120000))
	if err != nil {
		return actool.Errorf(err.Error()), nil
	}
	return actool.Text(actool.Capture(tc, body)), nil
}

// execPython writes the code to a temp .py under workDir/.tools, runs it via interp
// with the params JSON on stdin + scalar params mirrored to TOOL_<NAME> env, and
// returns combined stdout+stderr (with a timeout/exit note). Standalone + testable.
func execPython(ctx context.Context, interp, key, code string, params map[string]any, workDir string, sessionEnv []string, timeout time.Duration) (string, error) {
	// 沙箱（P1）：启用时把脚本写进沙箱目录（nobody 可读），进程降权 + 资源上限。
	// fail-closed：沙箱开启但初始化失败时拒绝执行，绝不退回以 root 裸跑。
	sbDir, sbErr := sandboxDir("")
	if sbErr != nil {
		return "", sbErr
	}
	var scriptPath string
	runDir := workDir
	if sbDir != "" {
		p, err := sandboxScriptPath(sbDir, "tool.py", code)
		if err != nil {
			return "", err
		}
		scriptPath = p
		runDir = sbDir
	} else {
		toolsDir := filepath.Join(workDir, ".tools")
		if err := os.MkdirAll(toolsDir, 0o755); err != nil {
			return "", err
		}
		f, err := os.CreateTemp(toolsDir, key+"-*.py")
		if err != nil {
			return "", err
		}
		scriptPath = f.Name()
		defer os.Remove(scriptPath)
		if _, err := f.WriteString(code); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	c := exec.CommandContext(runCtx, interp, scriptPath)
	c.Dir = runDir
	if sbDir != "" {
		// 沙箱模式：绝不透传 os.Environ()（其中含 ARTEX_PG_DSN/postgres 口令、
		// ARTEX_KEY_SECRET_FILE、LLM key 等）。只给最小环境 + 会话代理 env + TOOL_*。
		c.Env = minimalSandboxEnv(sbDir)
		for _, e := range sessionEnv {
			// 代理类变量只放行与出网相关的白名单，避免把宿主机凭据带进沙箱。
			if isSandboxSafeEnv(e) {
				c.Env = append(c.Env, e)
			}
		}
		applySandboxToCmd(c, sbDir) // 降权到 nobody + HOME/TMPDIR 指向沙箱
	} else {
		c.Env = append(os.Environ(), sessionEnv...) // 非沙箱：保持原行为
	}
	for k, v := range params { // 标量参数镜像成 TOOL_<NAME>
		if sv, ok := scalarStr(v); ok {
			c.Env = append(c.Env, "TOOL_"+strings.ToUpper(k)+"="+sv)
		}
	}
	pj, _ := json.Marshal(params)
	c.Stdin = bytes.NewReader(pj) // 参数 JSON 走 stdin
	out, err := c.CombinedOutput()
	body := string(out)
	if runCtx.Err() == context.DeadlineExceeded {
		body += "\n... [超时终止] ..."
	} else if err != nil {
		body += "\n[exit: " + err.Error() + "]"
	}
	return body, nil
}

// ---------- http:原生请求 + 代理 ----------

func (s *Server) runHTTPTool(ctx context.Context, execRaw json.RawMessage, params map[string]any, tc *actool.ToolContext) (actool.Result, error) {
	var spec httpExec
	_ = json.Unmarshal(execRaw, &spec)
	method := strings.ToUpper(strings.TrimSpace(spec.Method))
	if method == "" {
		method = "GET"
	}
	rawURL := renderTemplate(spec.URL, params, identity)
	if strings.TrimSpace(rawURL) == "" {
		return actool.Errorf("http url 为空"), nil
	}
	// 网络范围硬门禁(P0-2)：http 自定义工具是最容易被诱导打内网/云元数据的出口，
	// 这里在发请求前先做一次主机判定；真正防 DNS rebinding 的判定在 dial 层。
	if err := scope.ValidateURL(rawURL); err != nil {
		// 门禁拒绝是安全事件：落审计（便于发现 Agent 被诱导打内网/元数据）。
		if s.m != nil && s.m.pg != nil {
			s.insertAudit(db.AuditEntry{
				Actor: "tool", Action: "scope.denied", Target: truncateStr(rawURL, auditTargetMax),
				Status: 403, Detail: truncateStr(err.Error(), auditDetailMax),
			})
			s.securityAlert(nil, "scope.denied", truncateStr(rawURL, auditTargetMax), truncateStr(err.Error(), auditDetailMax))
		}
		return actool.Errorf("网络范围门禁拒绝: " + err.Error()), nil
	}
	var bodyReader io.Reader
	if spec.Body != "" {
		bodyReader = strings.NewReader(renderTemplate(spec.Body, params, identity))
	}
	runCtx, cancel := context.WithTimeout(ctx, timeoutOr(spec.TimeoutMs, 30000))
	defer cancel()
	req, err := http.NewRequestWithContext(runCtx, method, rawURL, bodyReader)
	if err != nil {
		return actool.Errorf(err.Error()), nil
	}
	for k, v := range spec.Headers {
		req.Header.Set(k, renderTemplate(v, params, identity))
	}
	// 每跳重定向都重新过门禁，防止「公网 302 到 127.0.0.1」绕过。
	client := &http.Client{
		Timeout:       timeoutOr(spec.TimeoutMs, 30000),
		CheckRedirect: scope.CheckRedirect,
	}
	if tr := s.httpProxyTransport(spec); tr != nil {
		// 显式代理场景（录制代理等）：目标由代理解析，dial 层看不到真实目标，
		// 因此只有上面的 URL 预检生效；这是已知取舍，已在报告中记录。
		client.Transport = tr
	} else {
		// 直连场景：关闭环境变量代理，改用带门禁的 DialContext，连接前校验并
		// 用已校验 IP 建连，避免 DNS rebinding。
		client.Transport = &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return scope.DialContext(ctx, network, addr)
		}}
	}
	resp, err := client.Do(req)
	if err != nil {
		return actool.Errorf("请求失败: " + err.Error()), nil
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	out := map[string]any{"status": resp.StatusCode, "body": string(respBody)}
	b, _ := json.Marshal(out)
	return actool.Text(actool.Capture(tc, string(b))), nil
}

// httpProxyTransport builds a Transport for the http tool's proxy config, or nil
// (direct). use_recording_proxy routes through the recording proxy + trusts its CA.
func (s *Server) httpProxyTransport(spec httpExec) *http.Transport {
	proxyStr := strings.TrimSpace(spec.Proxy)
	var caFile string
	if spec.UseRecordingProxy {
		if addr := s.m.ProxyAddr(); addr != "" {
			proxyStr = "http://" + addr
			caFile = s.m.ProxyCACert()
		}
	}
	if proxyStr == "" {
		return nil
	}
	pu, err := url.Parse(proxyStr)
	if err != nil {
		return nil
	}
	tr := &http.Transport{Proxy: http.ProxyURL(pu)}
	if caFile != "" {
		if pem, err := os.ReadFile(caFile); err == nil {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM(pem) {
				tr.TLSClientConfig = &tls.Config{RootCAs: pool}
			}
		}
	}
	return tr
}

// ---------- helpers ----------

func identity(s string) string { return s }

// shellQuote single-quotes a value for safe shell interpolation.
func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// renderTemplate replaces {name} placeholders with each param's rendered value.
func renderTemplate(tmpl string, params map[string]any, quote func(string) string) string {
	out := tmpl
	for k, v := range params {
		out = strings.ReplaceAll(out, "{"+k+"}", quote(valToStr(v)))
	}
	return out
}

// valToStr renders a param value: scalars as-is, arrays/objects as compact JSON.
func valToStr(v any) string {
	if sv, ok := scalarStr(v); ok {
		return sv
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// scalarStr returns (string, true) for scalar values, ("", false) for arrays/objects.
func scalarStr(v any) (string, bool) {
	switch x := v.(type) {
	case string:
		return x, true
	case bool:
		return fmt.Sprintf("%t", x), true
	case float64:
		if x == float64(int64(x)) {
			return fmt.Sprintf("%d", int64(x)), true
		}
		return fmt.Sprintf("%g", x), true
	case nil:
		return "", true
	default:
		return "", false
	}
}
