package server

import (
	"testing"

	"github.com/Aurox/artex/scope"

	actool "github.com/Autumn-27/norma/tool"
)

// TestScopeBlockedConsistency 防止「两份门禁实现」漂移：
// artex/scope.Blocked 用于 HTTP 自定义工具/代理/拨号层，norma/tool.ScopeBlocked
// 是 WebFetch 用的副本（norma 不能 import artex）。两者对同一目标必须给出相同结论，
// 否则会出现「一个通道封了、另一个通道漏了」。
func TestScopeBlockedConsistency(t *testing.T) {
	cases := []string{
		"127.0.0.1",
		"127.0.0.1:8787",
		"169.254.169.254",
		"[::1]:80",
		"::1",
		"10.0.0.5",
		"172.16.3.4",
		"192.168.1.2",
		"8.8.8.8",
		"example.com",
		"localhost",
		"metadata.google.internal",
		"2130706433",
		"0x7f000001",
		"0177.0.0.1",
		"127.0.0.1.",
		"0.0.0.0",
		"224.0.0.1",
	}
	for _, h := range cases {
		want, _ := scope.Blocked(h)
		got, _ := actool.ScopeBlocked(h)
		if want != got {
			t.Errorf("门禁策略漂移 host=%q: artex/scope.Blocked=%v norma/tool.ScopeBlocked=%v", h, want, got)
		}
	}
}
