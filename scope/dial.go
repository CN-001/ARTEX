package scope

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// ValidateURL 校验一个绝对 URL 的主机是否命中网络范围门禁。
// 仅允许 http/https；任何解析失败都按拒绝处理（fail-closed），因为到这一步
// 参数完全由调用方（模型/工具参数）控制，宽容处理等于留下绕过面。
func ValidateURL(raw string) error {
	if raw == "" {
		return errors.New("url 为空")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("url 解析失败: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("仅允许 http/https，收到 %q", u.Scheme)
	}
	host := u.Host
	if host == "" {
		return errors.New("url 缺少主机")
	}
	if blocked, reason := Blocked(host); blocked {
		return fmt.Errorf("目标命中平台网络范围门禁: %s", reason)
	}
	return nil
}

// CheckRedirect 用于 http.Client.CheckRedirect：每一跳重定向都重新过门禁，
// 防止「公网域名 302 到 127.0.0.1」这类经典绕过。同时限制最多 10 跳。
func CheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("重定向次数过多")
	}
	return ValidateURL(req.URL.String())
}

// Dial 是 http.Transport.DialContext 的实现：先门禁、再解析、再用「已校验的 IP」
// 实际建连，避免 DNS rebinding（解析时公网、连接时私网）。
// 允许列表命中的主机直接按原名拨号（授权场景，例如本机/私网靶场）。
func Dial(network, addr string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return DialContext(ctx, network, addr)
}

// DialContext 与 Dial 相同，但接受调用方 context。
func DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		// 没有端口（非 http 语义）时保守拒绝。
		return nil, fmt.Errorf("地址缺少端口: %w", err)
	}
	d := &net.Dialer{Timeout: 15 * time.Second}

	// 授权例外：允许列表里的主机不做私网判定，按原名直连。
	if Allowed(host) {
		return d.DialContext(ctx, network, net.JoinHostPort(host, port))
	}

	var ips []net.IP
	if ip := net.ParseIP(trimBrackets(host)); ip != nil {
		ips = []net.IP{ip}
	} else {
		ips, err = net.DefaultResolver.LookupIP(ctx, "ip", trimBrackets(host))
		if err != nil {
			// 与 Blocked 的取舍一致：DNS 层错误 fail-open（连接自然会失败），
			// 但一旦解析出任何私网地址就拒绝。
			return d.DialContext(ctx, network, net.JoinHostPort(host, port))
		}
	}
	for _, ip := range ips {
		if blocked, reason := blockedIP(ip); blocked {
			return nil, fmt.Errorf("目标命中平台网络范围门禁: %s → %s", trimBrackets(host), reason)
		}
	}
	if len(ips) == 0 {
		return nil, errors.New("解析结果为空")
	}
	// 用已校验的 IP 建连，避免二次解析（DNS rebinding）。
	return d.DialContext(ctx, network, net.JoinHostPort(ips[0].String(), port))
}

func trimBrackets(h string) string {
	if len(h) >= 2 && h[0] == '[' && h[len(h)-1] == ']' {
		return h[1 : len(h)-1]
	}
	return h
}
