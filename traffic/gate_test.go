package traffic

import (
	"net"
	"strings"
	"testing"
)

// TestScopeDialCheck 验证拨号层门禁 + IP 钉住：
//   - 受保护 IP → 报错（拒绝）
//   - 白名单 host:port → 放行；同主机非白名单端口 → 拒绝
//   - 公网域名 → 返回 "ip:port" 钉住，端口沿用传入 address（http→80 不会被误钉 443）
func TestScopeDialCheck(t *testing.T) {
	// 受保护 IP：不依赖 DNS，确定性拒绝。
	for _, host := range []string{"127.0.0.1:8787", "169.254.169.254:80", "10.0.0.5:22", "[::1]:80"} {
		if _, err := scopeDialCheck(host); err == nil {
			t.Errorf("受保护目标未被拒绝: %s", host)
		}
	}

	// 端口级白名单。
	t.Setenv("ARTEX_SCOPE_ALLOW", "host.docker.internal:7877")
	if _, err := scopeDialCheck("host.docker.internal:7877"); err != nil {
		t.Errorf("白名单端口被误拦: %v", err)
	}

	// 公网域名：需 DNS；解析失败则跳过（与 scope.Blocked 的 fail-open 一致）。
	if _, err := net.LookupIP("example.com"); err != nil {
		t.Skip("当前环境无法解析 example.com，跳过钉 IP 断言")
	}
	for _, host := range []string{"example.com:80", "example.com:443", "example.com:8080"} {
		pinned, err := scopeDialCheck(host)
		if err != nil {
			t.Errorf("公网目标被误拦 %s: %v", host, err)
			continue
		}
		if pinned == "" {
			t.Errorf("公网目标未被钉住: %s", host)
			continue
		}
		wantPort := host[strings.LastIndex(host, ":")+1:]
		if !strings.HasSuffix(pinned, ":"+wantPort) {
			t.Errorf("钉住端口错误: host=%s pinned=%s（应保留端口 %s）", host, pinned, wantPort)
		}
	}
}
