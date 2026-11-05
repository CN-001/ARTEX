package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// toolSandboxEnabled 报告是否启用「工具执行沙箱」（P1）。
// 通过 ARTEX_TOOL_SANDBOX=1 显式开启：命令/脚本类自定义工具将以 nobody 身份、
// 在独立临时目录、带资源上限运行，避免 Agent 的工具调用直接以容器 root 执行。
//
// 默认关闭的原因：ARTEX 本身就是渗透测试平台，很多工具链需要写 /app/data、
// 读工具目录、绑定低端口；盲目降权会把正常功能打坏。生产/共享环境建议开启。
func toolSandboxEnabled() bool {
	return strings.TrimSpace(os.Getenv("ARTEX_TOOL_SANDBOX")) == "1"
}

// sandboxDir 创建一个 nobody 可读写的独立工作目录。属主直接 chown 给 65534 并设为
// 0700：避免历史上用 0777 带来的同机其他进程/符号链接竞态问题。
//
// 返回 ("", err)：沙箱启用但初始化失败时调用方必须 fail-closed（拒绝执行），
// 绝不能退回「以 root 裸跑」——那等于静默取消隔离。
func sandboxDir(base string) (string, error) {
	if !toolSandboxEnabled() {
		return "", nil
	}
	if base == "" {
		base = os.TempDir()
	}
	dir, err := os.MkdirTemp(base, "artex-sbx-")
	if err != nil {
		return "", fmt.Errorf("沙箱目录创建失败: %w", err)
	}
	if err := os.Chown(dir, 65534, 65534); err != nil {
		_ = os.RemoveAll(dir)
		return "", fmt.Errorf("沙箱目录 chown 失败（不允许以 root 回落）: %w", err)
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		_ = os.RemoveAll(dir)
		return "", fmt.Errorf("沙箱目录 chmod 失败: %w", err)
	}
	return dir, nil
}

// minimalSandboxEnv 构造沙箱进程的最小环境变量集合。
// 关键：**不继承** os.Environ()——父进程环境里有数据库连接串、主密钥路径、
// LLM API key 等，降权后的工具进程若能读到，沙箱的隔离意义就被抵消了。
//
// 沙箱工具唯一出网通道是录制代理（容器内 127.0.0.1:8788）：该代理内置 scope 门禁
// （traffic/gate.go），因此 HTTP(S)_PROXY 指向它是安全的——既是出网路径，也保证
// 私网/云元数据在代理层被拒。
func minimalSandboxEnv(dir string) []string {
	proxy := sandboxProxy()
	env := []string{
		"HOME=" + dir,
		"TMPDIR=" + dir,
		"PATH=/usr/local/bin:/usr/bin:/bin",
		"LANG=C.UTF-8",
	}
	if proxy != "" {
		env = append(env, "HTTP_PROXY="+proxy, "HTTPS_PROXY="+proxy, "http_proxy="+proxy, "https_proxy="+proxy)
	}
	return env
}

// sandboxProxy 返回沙箱工具应使用的出网代理（默认容器内录制代理 8788）。
func sandboxProxy() string {
	if v := strings.TrimSpace(os.Getenv("ARTEX_SANDBOX_PROXY")); v != "" {
		if v == "off" {
			return ""
		}
		return v
	}
	return "http://127.0.0.1:8788"
}

// isSandboxSafeEnv 只放行与运行相关、不含凭据的环境变量。
//
// 刻意 **不放行** HTTP(S)_PROXY 等代理变量：容器里这些指向宿主录制代理
// (host.docker.internal:7877)，放行等于给沙箱一条绕过 scope 门禁/出网限制的通道，
// 且代理本身可能被用来访问内网。沙箱工具应走直连（并由内核层按 uid 限制私网）。
func isSandboxSafeEnv(kv string) bool {
	k := kv
	if i := strings.IndexByte(kv, '='); i >= 0 {
		k = kv[:i]
	}
	switch k {
	case "LANG", "LC_ALL", "TZ":
		return true
	}
	return false
}

// sandboxScriptPath 把命令/脚本落盘到沙箱目录，并赋予可执行权限。
func sandboxScriptPath(dir, name, content string) (string, error) {
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o755); err != nil {
		return "", err
	}
	_ = os.Chmod(p, 0o755)
	return p, nil
}
