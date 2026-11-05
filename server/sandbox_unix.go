//go:build !windows

package server

import (
	"os/exec"
	"strconv"
	"syscall"
)

// sandboxUID/GID = nobody/nogroup（Debian slim 默认 65534）。
const (
	sandboxUID = 65534
	sandboxGID = 65534
)

// sandboxWrapCommand 把一条 shell 命令包装成「降权 + 资源上限 + 独立 TMP/HOME」形式：
//
//	setpriv --reuid=65534 --regid=65534 --clear-groups \
//	  env HOME=<dir> TMPDIR=<dir> PATH=... sh -c 'ulimit ...; cd <dir>; exec <script>'
//
// 命令体先落盘为脚本再执行，避免二次转义带来的注入风险。
// 返回 "" 表示未启用沙箱。
func sandboxWrapCommand(dir, scriptPath string) string {
	if dir == "" || scriptPath == "" {
		return ""
	}
	// 资源上限：虚拟内存 2GiB、文件大小 256MiB、打开文件数 4096；
	// ulimit 在降权后对 setuid/硬上限之外的项仍可收紧。
	return "setpriv --reuid=" + strconv.Itoa(sandboxUID) + " --regid=" + strconv.Itoa(sandboxGID) +
		" --clear-groups env HOME=" + dir + " TMPDIR=" + dir +
		" PATH=/usr/local/bin:/usr/bin:/bin" +
		" sh -c 'ulimit -v 2097152 2>/dev/null; ulimit -f 262144 2>/dev/null; ulimit -n 4096 2>/dev/null; cd " +
		dir + " && exec sh " + scriptPath + "'"
}

// applySandboxToCmd 让 *exec.Cmd 以 nobody 身份、在沙箱目录内执行（脚本类工具）。
func applySandboxToCmd(c *exec.Cmd, dir string) {
	if c == nil || dir == "" {
		return
	}
	c.Dir = dir
	c.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: sandboxUID, Gid: sandboxGID},
	}
	c.Env = append(c.Env,
		"HOME="+dir,
		"TMPDIR="+dir,
		"PATH=/usr/local/bin:/usr/bin:/bin",
	)
	// 不再放宽目录权限：沙箱目录由 sandboxDir 以 chown 65534 + 0700 创建，
	// 如果 chown 失败则调用方应 fail-closed（拒绝以 root 执行），而不是退回 0777。
}
