package traffic

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/Aurox/artex/scope"
	mproxy "github.com/lqqyt2423/go-mitmproxy/proxy"
)

// scopeDialCheck 是代理拨号前的门禁 + IP 钉住回调（挂在 mproxy.Proxy.SetDialCheck）。
//
// 为什么需要它：scopeGate 在请求头阶段判定的是「主机名」，真正拨号时 DNS 会再次解析，
// 存在 DNS rebinding（校验时公网、拨号时私网）的 TOCTOU 窗口。这里在校验通过后把
// 解析出的 IP 返回给代理，代理直接对该 IP 拨号，校验对象与实际连接对象被绑定。
//
// 返回 ("", nil) 表示按原地址拨号（如命中 ARTEX_SCOPE_ALLOW 白名单、或 DNS 解析失败时
// 维持与 scope.Blocked 一致的 fail-open 取舍）。
func scopeDialCheck(host string) (string, error) {
	if host == "" {
		return "", nil
	}
	if scope.Allowed(host) {
		return "", nil
	}
	h, port := host, "443"
	if hh, pp, err := net.SplitHostPort(host); err == nil {
		h, port = hh, pp
	}
	h = strings.Trim(h, "[]")
	// 本身就是 IP：直接判定，无需钉住。
	if ip := net.ParseIP(h); ip != nil {
		if blocked, reason := scope.Blocked(ip.String()); blocked {
			return "", fmt.Errorf("拨号门禁拒绝 %s: %s", host, reason)
		}
		return "", nil
	}
	ips, err := net.LookupIP(h)
	if err != nil || len(ips) == 0 {
		// 拨号层默认 fail-closed：解析失败就无法证明目标是公网，
		// 宁可拒绝（HTTP 工具的 scope.Blocked 仍保持 fail-open 以避免误伤，
		// 需要旧行为时设 ARTEX_SCOPE_DIAL_FAILOPEN=1）。
		if os.Getenv("ARTEX_SCOPE_DIAL_FAILOPEN") == "1" {
			return "", nil
		}
		return "", fmt.Errorf("拨号门禁拒绝 %s: DNS 解析失败（fail-closed）", host)
	}
	for _, ip := range ips {
		if blocked, reason := scope.Blocked(ip.String()); blocked {
			return "", fmt.Errorf("拨号门禁拒绝 %s(%s): %s", host, ip, reason)
		}
	}
	return net.JoinHostPort(ips[0].String(), port), nil
}

// scopeGate 在录制代理层执行网络范围门禁。
//
// 为什么必须在代理层做：容器内没有直连公网，工具的出网路径是「工具 → 录制代理 →
// 上游」。沙箱工具被 iptables 限制只能连 127.0.0.1:8788（本代理），因此这里就是
// 沙箱外网的唯一咽喉——若不在此判定，沙箱里的 curl -x 就能打到 127.0.0.1/云元数据/
// 内网，绕过 scope 包对 HTTP 工具的限制。
//
// 覆盖的入口（只挂 Request 钩子是不够的：CONNECT 隧道、WebSocket 升级、大包体
// 的 stream 模式不会走完整 Request 流程）：
//   - Requestheaders：所有请求（含 CONNECT/WS 升级）最先触发，必须在这里拦；
//   - Request：兜底再判一次。
//
// 判定复用 scope.Blocked（回环/链路本地/私网/保留主机名/写法归一化 +
// ARTEX_SCOPE_ALLOW 例外），保证「HTTP 工具」与「代理出网」策略一致。
type scopeGate struct {
	mproxy.BaseAddon
}

func (g *scopeGate) Requestheaders(f *mproxy.Flow) { g.gate(f) }
func (g *scopeGate) Request(f *mproxy.Flow)       { g.gate(f) }

func (g *scopeGate) gate(f *mproxy.Flow) {
	if f == nil || f.Request == nil {
		return
	}
	host := requestHost(f)
	if host == "" {
		return
	}
	// 已被其它钩子判过就别重复拦截（避免日志刷屏）。
	if f.Response != nil && f.Response.StatusCode == http.StatusForbidden {
		return
	}
	blocked, reason := scope.Blocked(host)
	if !blocked {
		return
	}
	log.Printf("[traffic] 代理门禁拒绝 %s %s（%s）", f.Request.Method, host, reason)
	f.Response = &mproxy.Response{
		StatusCode: http.StatusForbidden,
		Header:     http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
		Body:       []byte("blocked by ARTEX network scope gate: " + reason),
	}
}

// requestHost 从 flow 里取目标主机：普通请求用 URL.Host；CONNECT 隧道用
// URL.Host（go-mitmproxy 把 authority 放在这里），兜底 URL.Opaque。
func requestHost(f *mproxy.Flow) string {
	if f.Request == nil {
		return ""
	}
	if f.Request.URL != nil {
		if h := strings.TrimSpace(f.Request.URL.Host); h != "" {
			return h
		}
		if o := strings.TrimSpace(f.Request.URL.Opaque); o != "" {
			return o
		}
	}
	return ""
}
