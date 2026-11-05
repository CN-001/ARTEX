---
name: awd-defense
description: AWD 攻防赛防守方向：WebShell 检测与清除、Web 攻击识别与自动封禁、二进制补丁对比定位漏洞、服务加固与应急响应。适用于"守机器/守服务"的比赛与演练。
license: MIT
compatibility: 仅限 CTF / 靶场 / 授权环境
---

# AWD 防守（守机器 / 守服务）

## 0. 先明确：AWD 防守的胜负逻辑

AWD 不是"不被攻破"，而是**在被攻破后最快发现、最快止损、最快修回来**。所以优先级是：

```
① 发现（检测）   → WebShell 落地立刻知道
② 止损（封禁）   → 攻击 IP 秒封
③ 修复（补洞）   → 用补丁对比定位改动点
④ 反打（拿分）   → 这是攻击侧的事，专注守方时不要分心
```

## 1. WebShell 检测（被上传木马的第一时间发现）

```bash
# 全量扫描 + 建基线（首轮无告警，之后新增/改动文件都会报）
python webshell_scan.py

# 指定 Web 目录
python webshell_scan.py --roots /var/www/html,/app/public

# 重建基线（确认当前状态干净后）
python webshell_scan.py --rebaseline

# 看报告
cat webshell_report.txt
```

**判读**：
- 评分 ≥15 → **P0**，几乎确定是木马（菜刀/冰蝎/哥斯拉/蚁剑特征）
- 评分 8–14 → **P1**，可疑（混淆、危险函数密集），需人工确认
- "Web 目录新增文件" → 只要不是你自己传的，**先隔离再分析**（别直接删，删了就没证据了）

**手工确认一个可疑文件**：
```bash
head -c 500 suspect.php           # 看开头
grep -oE 'eval|assert|system|base64_decode' suspect.php | sort | uniq -c
# 用 imgread 之类的思路：把二进制/混淆内容转成可读文本再判断
```

## 2. Web 攻击识别 + 自动封禁（fail2ban 等价物）

```bash
# 配置日志路径后运行（web_waf_config.json 的 log_paths）
python web_waf.py --once     # 跑一轮
python web_waf.py            # 常驻
```

识别：SQL注入 / 路径穿越 / WebShell 访问 / 漏洞扫描器(UA) / 敏感文件探测 / Log4Shell / 命令注入 / XSS / XXE / 文件上传攻击
处置：超阈值 → 调防火墙封禁 **`DEFENSE-IPBLOCK-<ip>`**（与溯源封禁共用命名，桌面「11-解除全部自动封禁」可一键解开）

**注意**：内网/本机 IP 在 `whitelist_ips` 里，不会误封自己。

## 3. 二进制补丁对比（定位漏洞 / 验证补丁）

赛方给了新版二进制？对比就知道他们改了哪里 —— **改动量最小的那个函数就是漏洞点**。

```bash
bindiff old_binary new_binary            # 文本报告
bindiff old new --json diff.json         # 导出 JSON
bindiff old new --all                    # 显示全部改动函数
```

输出解读：
```
未变: 116   改动: 1   新增: 0   删除: 0
  ● 函数名
    大小 120 → 132 (+12)   指令差异 3 条
      ~ cmp $0x10,%eax          ← 原来只检查 16
        cmp $0x100,%eax         ← 补丁改成检查 256
```
- **改动小（1–5 条指令）+ 涉及 cmp/test/jmp** → 典型边界检查修复 = **漏洞点**
- **新增函数** → 可能是新功能，也可能是**后门**（`--all` 看全）
- 被 strip 的二进制只能看到地址形式的函数名，此时上 Ghidra

配合 Ghidra 深挖：
```bash
/app/data/tools/ghidra/support/analyzeHeadless /tmp/proj proj1 -import suspect.exe \
  -postScript DumpDecomp.java -scriptPath /app/data/tools/ghidra-scripts
```

## 4. 应急响应顺序（发现被攻破时）

```
1. 保证据（先做！）   : 导日志、留木马样本、记进程/连接
2. 判断入口           : WebShell? 弱口令? 已知漏洞? → 对应日志段
3. 断外联             : 封攻击 IP；必要时断网（桌面「13-断网」）
4. 清木马             : 但要先留样本！清完再扫一遍确认
5. 修漏洞             : 补丁对比定位 → 改代码/配置 → 复测
6. 复盘               : 写进 incident report，改检测规则防再犯
```

**永远不要**：不保证据就删木马、不确认就重启服务（会丢内存证据）。

## 5. 纪律

- 检测工具的结果要**交叉验证**（特征命中 ≠ 一定是木马，可能是正常的管理工具）
- 封禁前确认不是自己人（内网段、跳板机、扫描器）
- 所有动作留日志，比赛结束后要能复盘
