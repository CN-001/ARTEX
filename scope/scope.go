// Package scope implements the network-range hard gate (P0-2): before any
// platform-side request reaches the network, the target host is normalized and
// checked against reserved/private ranges so SSRF against the loopback, cloud
// metadata (169.254.169.254) or the internal network is denied by default.
//
// 门禁默认开启、无需额外配置:命中回环/链路本地/私网/保留段的目标一律拦截;唯一例外
// 是环境变量 ARTEX_SCOPE_ALLOW(逗号分隔的 IP/CIDR/域名),给「本机靶场在私网」这类
// 已授权场景留口。DNS 解析失败按 fail-open 放行(取舍见 Blocked 注释)。
package scope

import (
	"net"
	"os"
	"strings"
)

// 预留/私网 CIDR 判定表。顺序无关,命中任一即 blocked。
var (
	cidrLoopback4  = mustCIDR("127.0.0.0/8")
	cidrLoopback6  = mustCIDR("::1/128")
	cidrLinkLocal  = mustCIDR("169.254.0.0/16") // 含云元数据 169.254.169.254
	cidrLinkLocal6 = mustCIDR("fe80::/10")
	cidrPrivate10  = mustCIDR("10.0.0.0/8")
	cidrPrivate172 = mustCIDR("172.16.0.0/12")
	cidrPrivate192 = mustCIDR("192.168.0.0/16")
	cidrPrivate6   = mustCIDR("fc00::/7")
	cidrUnspec4    = mustCIDR("0.0.0.0/32")
	cidrUnspec6    = mustCIDR("::/128")
	cidrMulticast  = mustCIDR("224.0.0.0/4")
)

func mustCIDR(s string) *net.IPNet {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		panic("scope: invalid built-in CIDR " + s)
	}
	return n
}

// normalizeHost strips a trailing :port (including [ipv6]:port), lowercases and
// trims a trailing dot (FQDN 写法 "localhost." 与 "localhost" 等价，否则可绕过后缀/精确判定),
// returning the bare host for matching. Bare IPs / bare domains are unchanged.
func normalizeHost(host string) string {
	host = strings.TrimSpace(host)
	if strings.HasPrefix(host, "[") {
		host = trimBrackets(host)
	} else if h, _, err := net.SplitHostPort(host); err == nil {
		host = trimBrackets(h)
	}
	host = strings.ToLower(host)
	// 仅去一个尾部点：多个点属于非法域名，交给 DNS 失败 → fail-open 处理。
	host = strings.TrimSuffix(host, ".")
	return host
}

// parseNumericHost 解析「非标准写法的 IPv4」（十进制整数 / 0x 十六进制 / 前导零八进制），
// 例如 2130706433、0x7f000001、0177.0.0.1 都等价于 127.0.0.1。net.ParseIP 不认这些写法，
// 若不处理，它们会被当作域名交给 DNS，从而绕过私网判定。
func parseNumericHost(h string) net.IP {
	// 十进制整数（无点）。
	if isAllDigits(h) {
		var n uint64
		for i := 0; i < len(h); i++ {
			n = n*10 + uint64(h[i]-'0')
			if n > 0xFFFFFFFF {
				return nil
			}
		}
		return uint32ToIP(uint32(n))
	}
	// 0x 开头十六进制整数。
	if strings.HasPrefix(h, "0x") {
		var n uint64
		for i := 2; i < len(h); i++ {
			c := h[i]
			var v uint64
			switch {
			case c >= '0' && c <= '9':
				v = uint64(c - '0')
			case c >= 'a' && c <= 'f':
				v = uint64(c-'a') + 10
			default:
				return nil
			}
			n = n*16 + v
			if n > 0xFFFFFFFF {
				return nil
			}
		}
		return uint32ToIP(uint32(n))
	}
	// 点分但含前导零（八进制段），如 0177.0.0.1。
	if strings.Contains(h, ".") {
		parts := strings.Split(h, ".")
		if len(parts) > 4 {
			return nil
		}
		octal := false
		for _, p := range parts {
			if len(p) > 1 && p[0] == '0' {
				octal = true
			}
			if !isAllDigits(p) {
				return nil
			}
		}
		if !octal {
			return nil // 普通点分 IPv4 交给 net.ParseIP
		}
		ip := make(net.IP, 4)
		for i := 0; i < 4; i++ {
			v := 0
			if i < len(parts) {
				for j := 0; j < len(parts[i]); j++ {
					v = v*8 + int(parts[i][j]-'0')
				}
			}
			ip[i] = byte(v)
		}
		return ip
	}
	return nil
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func uint32ToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

// Blocked reports whether host is forbidden by the network-range gate. The host
// may carry a port ("10.0.0.5:8080") or IPv6 brackets ("[::1]:80") or be a bare
// domain. 判定顺序:Allowed(放行) → 私网/保留段(拦截) → 放行。
//
// 取舍说明:除 DNS 解析错误外,门禁采用「默认拦截」(fail-closed)——只有能被明确判定为
// 公网/允许列表的目标才放行。唯一例外是 net.LookupIP 返回错误(NXDOMAIN 等):此时
// fail-open 放行,避免误伤合法的公网目标。DNS rebinding(先解析公网、连接时再解析回
// 私网)属于后续 TODO,不在本次范围内。
func Blocked(host string) (bool, string) {
	// 1. 允许列表优先:命中则即使私网也放行(授权靶场留口)。
	if Allowed(host) {
		return false, ""
	}
	h := normalizeHost(host)

	// 2. 主机名精确/后缀判定(便宜且确定,先做)。
	if b, reason := blockedHostname(h); b {
		return true, reason
	}

	// 3. 非标准写法的 IPv4（十进制/十六进制/八进制）先归一化成 IP 再判定，
	//    否则它们会被当作域名交给 DNS（解析结果不确定或失败 → fail-open 绕过）。
	if ip := parseNumericHost(h); ip != nil {
		if b, reason := blockedIP(ip); b {
			return true, reason + "(非标准写法 " + h + ")"
		}
		return false, ""
	}

	// 4. 若本身就是 IP，直接判定。
	if ip := net.ParseIP(h); ip != nil {
		if b, reason := blockedIP(ip); b {
			return true, reason
		}
		return false, "" // 公网 IP → 放行
	}

	// 5. 域名:解析后逐个判定,任一命中私网即拦截(原因保留原始 host)。
	if ips, err := net.LookupIP(h); err == nil {
		for _, ip := range ips {
			if b, reason := blockedIP(ip); b {
				return true, "域名 " + host + " 解析到 " + ip.String() + "(" + reason + ")"
			}
		}
	}
	// DNS 失败 → fail-open(见上方注释)。
	return false, ""
}

// blockedHostname 判定保留主机名:精确 localhost/metadata 及 .internal/.local 后缀。
func blockedHostname(h string) (bool, string) {
	switch h {
	case "localhost", "metadata.google.internal", "metadata":
		return true, "保留主机名(" + h + ")"
	}
	if strings.HasSuffix(h, ".internal") || strings.HasSuffix(h, ".local") {
		return true, "保留主机名后缀(" + h + ")"
	}
	return false, ""
}

// blockedIP 判定单个 IP 命中的保留/私网段,返回中文原因。
func blockedIP(ip net.IP) (bool, string) {
	switch {
	case cidrLoopback4.Contains(ip) || cidrLoopback6.Contains(ip):
		return true, "回环地址"
	case cidrLinkLocal.Contains(ip) || cidrLinkLocal6.Contains(ip):
		return true, "链路本地地址(含云元数据 169.254.169.254)"
	case cidrPrivate10.Contains(ip) || cidrPrivate172.Contains(ip) ||
		cidrPrivate192.Contains(ip) || cidrPrivate6.Contains(ip):
		return true, "私网地址"
	case cidrUnspec4.Contains(ip) || cidrUnspec6.Contains(ip):
		return true, "未指定地址(0.0.0.0/::)"
	case cidrMulticast.Contains(ip):
		return true, "组播地址"
	}
	return false, ""
}

// Allowed reports whether host is explicitly allow-listed via ARTEX_SCOPE_ALLOW
// (comma-separated IP / CIDR / domain). A domain entry matches the exact host, any
// subdomain (".domain" suffix), or — when host is a literal IP — any IP that domain
// resolves to. 允许列表命中则 Blocked 直接放行,即便目标是私网。
func Allowed(host string) bool {
	raw := strings.TrimSpace(os.Getenv("ARTEX_SCOPE_ALLOW"))
	if raw == "" {
		return false
	}
	h := normalizeHost(host)
	hostIP := net.ParseIP(h)
	// host:port 的原始小写写法，用于「端口级白名单」精确匹配（见下）。
	rawHost := strings.ToLower(strings.TrimSpace(host))
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		entry = strings.TrimSuffix(entry, ".")
		// CIDR(含 "/"):对主机 IP 做 CIDR 包含匹配。
		if strings.Contains(entry, "/") {
			if _, n, err := net.ParseCIDR(entry); err == nil && hostIP != nil && n.Contains(hostIP) {
				return true
			}
			continue
		}
		// 带端口的白名单项（如 host.docker.internal:7877）：只放行该 host:port 组合，
		// 避免「放行一个域名 = 放行该主机所有端口」的过度授权。
		// 排除 IPv6 字面量（其中也含冒号）——IPv6 项按下面的裸 IP/CIDR 处理。
		if strings.Count(entry, ":") == 1 && !strings.HasPrefix(entry, "[") && !strings.Contains(entry, "]") {
			if strings.EqualFold(entry, rawHost) {
				return true
			}
			continue
		}
		// 裸 IP:精确相等匹配。
		if eip := net.ParseIP(entry); eip != nil {
			if hostIP != nil && hostIP.Equal(eip) {
				return true
			}
			continue
		}
		// 域名:精确 / 子域后缀 / host 是 IP 且等于该域名解析出的任一 IP。
		e := strings.ToLower(entry)
		if h == e || strings.HasSuffix(h, "."+e) {
			return true
		}
		if hostIP != nil {
			if ips, err := net.LookupIP(e); err == nil {
				for _, ip := range ips {
					if ip.Equal(hostIP) {
						return true
					}
				}
			}
		}
	}
	return false
}
