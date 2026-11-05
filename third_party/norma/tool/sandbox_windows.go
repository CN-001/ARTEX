//go:build windows

package tool

// Windows 没有 setpriv/nobody 等价机制：保持原有行为（不改变命令语义）。
// 跨平台编译需要同名符号。
func sandboxCommand(command string) string { return command }
