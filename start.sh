#!/bin/sh
# ARTEX 守护启动脚本（Linux / macOS / Docker ENTRYPOINT）
#
# 用法：
#   ./start.sh                       前台运行（Ctrl-C 停止）
#   nohup ./start.sh >artex.log 2>&1 &   后台常驻
#   ./start.sh -addr :9000           额外参数原样透传给 artex
#
# 它只做一件事：把 artex 跑起来，进程退出后按退出码决定要不要再拉起。
#
#   0      用户正常停止        → 退出循环
#   75     程序请求重启        → 立刻重跑（页面点了"一键更新"或"回滚"）
#   其他   崩溃                → 退避后重跑（1→2→4…最多 60 秒）
#
# 刻意不在这里做下载、SHA256 校验或换装：那些逻辑在 sh 和 bat 上要写两套，
# 而它们恰恰是最不能出错的一环——一旦换上跑不起来的二进制，本脚本会忠实地
# 反复拉起它，用户只能上机器手工救。所以校验/换装全部留在 Go 里（selfupdate 包），
# 由 artex 自己在启动时完成，脚本保持傻瓜化。
set -u

cd "$(dirname "$0")" || exit 1

BIN=./artex
[ -x "$BIN" ] || { echo "[artex] 找不到可执行文件 $BIN" >&2; exit 1; }

# ---- 启动期安全加固（失败不阻断启动，仅打印告警）----
# 1) 密钥文件权限收紧：容器内 named volume 上 chmod 生效；Windows bind mount 上
#    无效果（文件永远 0777），所以 compose 里改用 named volume 承载密钥。
harden_secrets() {
	if [ -d /app/secrets ]; then
		chmod 700 /app/secrets 2>/dev/null || true
		find /app/secrets -type f -exec chmod 600 {} \; 2>/dev/null || true
	fi
	if [ -f /run/secrets/artex-master-key ]; then
		chmod 600 /run/secrets/artex-master-key 2>/dev/null || true
	fi
	# Docker secret 在非 swarm 模式下是 Windows bind mount（容器内显示 0777），
	# 降权后的工具进程能直接读到。这里由 root 复制进 named volume 并 0600。
	if [ -f /run/secrets/pg_password ]; then
		cp /run/secrets/pg_password /app/secrets/pg_password 2>/dev/null || true
		chmod 600 /app/secrets/pg_password 2>/dev/null || true
	fi
}

# 2) 沙箱出网限制：降权后的工具进程（uid 65534）不许访问本机/私网/云元数据，
#    但保留 Docker 内置 DNS(127.0.0.11) 以便正常解析公网域名。需要 NET_ADMIN。
#    这是「Bash/command/script 无法用 curl 打内网」的实际落点：scope 包只管
#    HTTP 工具，任意 shell 的网络限制只能靠内核层。
harden_egress() {
	[ "${ARTEX_TOOL_SANDBOX:-0}" = "1" ] || return 0
	command -v iptables >/dev/null 2>&1 || { echo "[artex] 未找到 iptables，跳过出网限制" >&2; return 0; }
	iptables -C OUTPUT -m owner --uid-owner 65534 -d 10.0.0.0/8 -j REJECT 2>/dev/null && return 0
	# 先放行 DNS，再拒绝内网段（顺序即优先级）。
	# 沙箱工具的唯一出网通道：容器内录制代理（该代理内置 scope 门禁）。
	iptables -A OUTPUT -m owner --uid-owner 65534 -d 127.0.0.1 -p tcp --dport 8788 -j ACCEPT 2>/dev/null || true
	iptables -A OUTPUT -m owner --uid-owner 65534 -d 127.0.0.11 -p udp --dport 53 -j ACCEPT 2>/dev/null || true
	iptables -A OUTPUT -m owner --uid-owner 65534 -d 127.0.0.11 -p tcp --dport 53 -j ACCEPT 2>/dev/null || true
	for net in 127.0.0.0/8 169.254.0.0/16 10.0.0.0/8 172.16.0.0/12 192.168.0.0/16 100.64.0.0/10 fc00::/7 fe80::/10; do
		iptables -A OUTPUT -m owner --uid-owner 65534 -d "$net" -j REJECT 2>/dev/null || \
			ip6tables -A OUTPUT -m owner --uid-owner 65534 -d "$net" -j REJECT 2>/dev/null || true
	done
	# IPv6 回环单独处理：不能整段 REJECT，否则 [::1]:8788 的代理通道也被堵死。
	# 先放行代理端口，再拒绝 ::1/128 与 IPv4-mapped ::ffff:0:0/96。
	if command -v ip6tables >/dev/null 2>&1; then
		ip6tables -A OUTPUT -m owner --uid-owner 65534 -d ::1/128 -p tcp --dport 8788 -j ACCEPT 2>/dev/null || true
		ip6tables -A OUTPUT -m owner --uid-owner 65534 -d ::1/128 -j REJECT 2>/dev/null || true
		ip6tables -A OUTPUT -m owner --uid-owner 65534 -d ::ffff:0:0/96 -j REJECT 2>/dev/null || true
		# 沙箱 IPv6 出网一律拒绝（v4 公网直连本环境也不通，出网统一走 v4 代理）：
		# 不这么做的话，IPv6 目标可以完全绕过上面的 v4 规则。
		ip6tables -A OUTPUT -m owner --uid-owner 65534 -d ::/0 -j REJECT 2>/dev/null || true
	else
		echo "[artex] 未找到 ip6tables，IPv6 出网限制未生效" >&2
	fi
}

harden_secrets
harden_egress

RESTART_CODE=75
MAX_DELAY=60

child=0
stopping=0

# 转发停止信号给 artex 本体。
#
# Docker 下这是必需的：docker stop 只把 SIGTERM 发给 PID 1（也就是本脚本），
# 不会发给子进程。不转发的话 artex 收不到信号、做不了优雅关闭，10 秒后被 SIGKILL
# 硬杀，正在跑的任务直接断在半路。
forward() {
	stopping=1
	if [ "$child" -ne 0 ]; then
		kill -TERM "$child" 2>/dev/null || true
	fi
}
trap forward INT TERM

delay=1
while :; do
	"$BIN" "$@" &
	child=$!

	# 信号会打断 wait 并让它返回 >128。此时子进程其实还在做优雅关闭，
	# 必须再 wait 一次才能拿到它真正的退出码。
	wait "$child"
	code=$?
	if [ "$code" -gt 128 ]; then
		wait "$child"
		code=$?
	fi
	child=0

	if [ "$stopping" -eq 1 ]; then
		echo "[artex] 已停止"
		exit 0
	fi

	case "$code" in
		0)
			echo "[artex] 正常退出"
			exit 0
			;;
		"$RESTART_CODE")
			# 更新/回滚已就绪：重跑后 artex 会在启动时完成换装（见 selfupdate.Bootstrap）。
			echo "[artex] 请求重启（应用新版本）…"
			delay=1
			;;
		*)
			echo "[artex] 异常退出 (code=$code)，${delay}s 后重启" >&2
			sleep "$delay"
			delay=$((delay * 2))
			[ "$delay" -gt "$MAX_DELAY" ] && delay=$MAX_DELAY
			;;
	esac
done
