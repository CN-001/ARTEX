package scope

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 覆盖回环/元数据/私网/特殊/端口/IPv6/NXDOMAIN fail-open 以及允许列表的例外通道。
func TestBlocked(t *testing.T) {
	cases := []struct {
		name       string
		host       string
		allow      string
		wantBlock  bool
		wantReason string // 命中时要求包含的中文关键字(可空)
	}{
		// 回环
		{"loopback-ip", "127.0.0.1", "", true, "回环"},
		{"loopback-ip-port", "127.0.0.1:8080", "", true, "回环"},
		{"loopback-ipv6", "::1", "", true, "回环"},
		{"loopback-ipv6-bracket", "[::1]:80", "", true, "回环"},
		{"loopback-host", "localhost", "", true, "保留主机名"},
		{"loopback-host-port", "localhost:8080", "", true, "保留主机名"},
		// 链路本地(云元数据)
		{"metadata", "169.254.169.254", "", true, "链路本地"},
		{"metadata-port", "169.254.169.254:80", "", true, "链路本地"},
		{"link-local", "169.254.1.2", "", true, "链路本地"},
		{"link-local-v6", "fe80::1", "", true, "链路本地"},
		// 私网
		{"private-10", "10.0.0.5", "", true, "私网"},
		{"private-172", "172.16.0.1", "", true, "私网"},
		{"private-192", "192.168.1.1", "", true, "私网"},
		{"private-v6", "fc00::1", "", true, "私网"},
		// 特殊
		{"unspecified4", "0.0.0.0", "", true, "未指定"},
		{"unspecified6", "::", "", true, "未指定"},
		{"multicast", "224.0.0.1", "", true, "组播"},
		{"metadata-host", "metadata.google.internal", "", true, "保留主机名"},
		{"metadata-bare", "metadata", "", true, "保留主机名"},
		{"internal-suffix", "foo.internal", "", true, "保留主机名后缀"},
		{"local-suffix", "bar.foo.local", "", true, "保留主机名后缀"},
		// 公网放行
		{"public-ip", "8.8.8.8", "", false, ""},
		{"public-domain", "example.com", "", false, ""},
		// NXDOMAIN → fail-open
		{"nxdomain-fail-open", "nonexistent-host-xyz.invalid", "", false, ""},
		// 允许列表:命中私网也放行
		{"allow-cidr", "10.0.0.5", "10.0.0.0/8", false, ""},
		{"allow-exact-ip", "192.168.1.10", "192.168.1.10", false, ""},
		{"allow-exact-ip-miss", "192.168.1.11", "192.168.1.10", true, "私网"},
		{"allow-domain-exact", "internal.example.com", "internal.example.com", false, ""},
		{"allow-domain-sub", "api.internal.example.com", "internal.example.com", false, ""},
		{"allow-domain-miss", "internal.example.org", "internal.example.com", false, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("ARTEX_SCOPE_ALLOW", c.allow)
			blocked, reason := Blocked(c.host)
			if blocked != c.wantBlock {
				t.Fatalf("Blocked(%q) = %v (reason=%q), want blocked=%v", c.host, blocked, reason, c.wantBlock)
			}
			if c.wantBlock && c.wantReason != "" && !strings.Contains(reason, c.wantReason) {
				t.Fatalf("Blocked(%q) reason %q missing keyword %q", c.host, reason, c.wantReason)
			}
		})
	}
}

// TestBlockedPortAndIPv6 ensures the normalize path keeps the port out of the IP
// / hostname comparisons while IPv6 brackets are stripped.
func TestBlockedPortAndIPv6(t *testing.T) {
	t.Setenv("ARTEX_SCOPE_ALLOW", "")
	if b, _ := Blocked("10.0.0.5:8080"); !b {
		t.Fatal("10.0.0.5:8080 should be blocked (port stripped)")
	}
	if b, _ := Blocked("[::1]:80"); !b {
		t.Fatal("[::1]:80 should be blocked (brackets/port stripped)")
	}
}

// TestBlockedBypassForms 覆盖常见「写法绕过」：十进制/十六进制/八进制 IPv4、
// 尾部点、IPv4-mapped IPv6。这些写法 net.ParseIP 不认或易被当成域名，必须归一化后拦截。
func TestBlockedBypassForms(t *testing.T) {
	t.Setenv("ARTEX_SCOPE_ALLOW", "")
	cases := []string{
		"2130706433",      // 127.0.0.1 十进制
		"0x7f000001",      // 127.0.0.1 十六进制
		"0177.0.0.1",      // 127.0.0.1 八进制段
		"localhost.",      // 尾部点
		"127.0.0.1.",      // IP 尾部点
		"::ffff:127.0.0.1", // IPv4-mapped IPv6
		"::ffff:10.0.0.1",
		"[::ffff:169.254.169.254]:80",
		"169.254.169.254.",
	}
	for _, h := range cases {
		if b, reason := Blocked(h); !b {
			t.Errorf("Blocked(%q) = false, want true (bypass form); reason=%q", h, reason)
		}
	}
	// 公网十进制写法（8.8.8.8 = 134744072）不应被误拦。
	if b, reason := Blocked("134744072"); b {
		t.Errorf("Blocked(134744072) = true (%s), want false", reason)
	}
}

// TestValidateURL 覆盖 scheme 限制、主机门禁与非法输入（fail-closed）。
func TestValidateURL(t *testing.T) {
	t.Setenv("ARTEX_SCOPE_ALLOW", "")
	bad := []string{
		"http://127.0.0.1/admin",
		"https://169.254.169.254/latest/meta-data/",
		"http://[::1]:8080/",
		"http://localhost:18888/healthz",
		"file:///etc/passwd",
		"ftp://example.com/x",
		"http://10.1.2.3:80/",
		"",
		"http://",
	}
	for _, u := range bad {
		if err := ValidateURL(u); err == nil {
			t.Errorf("ValidateURL(%q) = nil, want error", u)
		}
	}
	good := []string{
		"https://example.com/x",
		"http://8.8.8.8/",
	}
	for _, u := range good {
		if err := ValidateURL(u); err != nil {
			t.Errorf("ValidateURL(%q) = %v, want nil", u, err)
		}
	}
}

// TestCheckRedirect 确保重定向到私网会被拦下。
func TestCheckRedirect(t *testing.T) {
	t.Setenv("ARTEX_SCOPE_ALLOW", "")
	req := httptest.NewRequest("GET", "http://127.0.0.1:8000/x", nil)
	if err := CheckRedirect(req, []*http.Request{req}); err == nil {
		t.Fatal("redirect to 127.0.0.1 should be rejected")
	}
	req2 := httptest.NewRequest("GET", "https://example.com/x", nil)
	if err := CheckRedirect(req2, []*http.Request{req2}); err != nil {
		t.Fatalf("redirect to public target should pass, got %v", err)
	}
}

// TestDialBlocksPrivate 验证 dial 层门禁（不依赖 DNS）：直接连 127.0.0.1 必须失败。
func TestDialBlocksPrivate(t *testing.T) {
	t.Setenv("ARTEX_SCOPE_ALLOW", "")
	if _, err := Dial("tcp", "127.0.0.1:9"); err == nil {
		t.Fatal("Dial(127.0.0.1) should be blocked by scope gate")
	}
	if _, err := Dial("tcp", "169.254.169.254:80"); err == nil {
		t.Fatal("Dial(169.254.169.254) should be blocked by scope gate")
	}
}
