//go:build !windows

package tool

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// sandboxUID/GID = nobody/nogroup（Debian slim 默认 65534）。
const (
	sandboxUID = 65534
	sandboxGID = 65534
)

// sandboxCommand 在 ARTEX_TOOL_SANDBOX=1 时把整条 shell 命令包装成「降权到 nobody +
// 资源上限 + 独立 HOME/TMPDIR」的形式执行，避免 Agent 的 Bash 工具以容器 root 运行。
//
// 设计要点：
//   - 命令先落盘成脚本再执行，避免二次转义引入注入面；
//   - 保留调用方的工作目录（exec.Cmd.Dir），只替换 HOME/TMPDIR，避免破坏相对路径；
//   - 资源上限可用环境变量调整：ARTEX_SANDBOX_MEM_MB（默认 4096，0/负数=不限）、
//     ARTEX_SANDBOX_FSIZE_MB（默认 512）。
//   - 需要 root 能力的扫描（如 nmap -sS、绑定低端口）会被限制；这类场景需要
//     把 ARTEX_TOOL_SANDBOX 设为 0 或用平台外的方式执行。
// sandboxProxy 返回沙箱出网代理（容器内录制代理，含 scope 门禁）。
func sandboxProxy() string {
	if v := strings.TrimSpace(os.Getenv("ARTEX_SANDBOX_PROXY")); v != "" {
		if v == "off" {
			return ""
		}
		return v
	}
	return "http://127.0.0.1:8788"
}

func sandboxCommand(command string) string {
	if strings.TrimSpace(os.Getenv("ARTEX_TOOL_SANDBOX")) != "1" {
		return command
	}
	dir, err := os.MkdirTemp("", "artex-sbx-")
	if err != nil {
		return "echo 'artex sandbox init failed: refusing to run unsandboxed' >&2; exit 1"
	}
	// 属主 chown 给 nobody 并 0700：目录里是待执行脚本，避免 0777 带来的
	// 同机其他进程改写脚本（符号链接/竞态的提权风险）。chown 失败时退回 0777。
	// fail-closed：chown 失败说明没法真正降权落盘，直接拒绝执行而不是退回 0777。
	if err := os.Chown(dir, sandboxUID, sandboxGID); err != nil {
		_ = os.RemoveAll(dir)
		return "echo 'artex sandbox init failed (chown): refusing to run as root' >&2; exit 1"
	}
	_ = os.Chmod(dir, 0o700)
	script := filepath.Join(dir, "cmd.sh")
	if err := os.WriteFile(script, []byte(command), 0o755); err != nil {
		return "echo 'artex sandbox init failed (write): refusing to run unsandboxed' >&2; exit 1"
	}
	_ = os.Chmod(script, 0o755)

	memMB := envInt("ARTEX_SANDBOX_MEM_MB", 4096)
	fsMB := envInt("ARTEX_SANDBOX_FSIZE_MB", 512)
	var limits []string
	if memMB > 0 {
		limits = append(limits, fmt.Sprintf("ulimit -v %d 2>/dev/null", memMB*1024))
	}
	if fsMB > 0 {
		limits = append(limits, fmt.Sprintf("ulimit -f %d 2>/dev/null", fsMB*2048))
	}
	limits = append(limits, "ulimit -n 8192 2>/dev/null")

	// trap 在命令结束/超时被杀时清理沙箱目录，避免 /tmp 里堆积分发脚本（含命令内容）。
	//
	// 关键：必须用 `env -i` 清空继承环境。GNU env 只做「继承 + 覆盖」，不加 -i 时
	// 父进程的 ARTEX_PG_DSN / ANTHROPIC_API_KEY / 主密钥路径 等会原样传给降权进程，
	// 沙箱就形同虚设。这里只保留运行必需变量 + 代理（代理本身带 scope 门禁）。
	proxy := sandboxProxy()
	return fmt.Sprintf(
		"setpriv --reuid=%d --regid=%d --clear-groups env -i HOME=%s TMPDIR=%s PATH=/usr/local/bin:/usr/bin:/bin HTTP_PROXY=%s HTTPS_PROXY=%s http_proxy=%s https_proxy=%s bash -c 'trap \"rm -rf %s\" EXIT; %s; exec bash %s'",
		sandboxUID, sandboxGID, dir, dir, proxy, proxy, proxy, proxy, dir, strings.Join(limits, "; "), script,
	)
}

func envInt(name string, def int) int {
	s := strings.TrimSpace(os.Getenv(name))
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
