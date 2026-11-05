//go:build windows

package server

import "os/exec"

// Windows 下不提供 setuid 降权沙箱（没有 nobody/setpriv 等价物）。
// 保留同名接口，保证跨平台编译；ARTEX_TOOL_SANDBOX=1 在 Windows 上为空操作，
// 工具仍以当前用户身份运行（与改造前一致）。
func sandboxWrapCommand(dir, scriptPath string) string { return "" }

func applySandboxToCmd(c *exec.Cmd, dir string) {}
