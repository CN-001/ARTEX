package server

import (
	"net/http"
	"strings"
)

// securityHeaders 给所有响应加一组浏览器侧安全头。
//
// 说明：
//   - 这套头挡不住服务端漏洞，但能显著缩小 XSS/点击劫持/MIME 嗅探的利用面，
//     而 ARTEX 的 UI 是静态导出（Next），没有服务端模板，靠头来兜底最划算。
//   - CSP 里的 'unsafe-inline' 是 Next 静态导出 hydration 内联脚本所必需；
//     没有外部源，因此 script-src 不允许任何第三方域。
//   - HSTS 仅在请求经 TLS（含可信反代 X-Forwarded-Proto: https）时下发，
//     避免纯 HTTP 内网部署被浏览器永久锁死。
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		if requestIsTLS(r) {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		// API 响应不需要 CSP（不是文档），给 UI 文档加。
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			h.Set("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self' 'unsafe-inline'; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' data: blob:; "+
					"font-src 'self' data:; "+
					"connect-src 'self' ws: wss:; "+
					"worker-src 'self' blob:; "+
					"frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		}
		next.ServeHTTP(w, r)
	})
}
