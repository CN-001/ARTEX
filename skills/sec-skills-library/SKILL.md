---
name: sec-skills-library
description: 网络安全技能图书馆（982 个攻击/取证/防御专项技能）。做任何安全任务前，先 grep 本库索引定位技能，再 Read 对应 SKILL.md 现学现用。覆盖 Web/内网AD/云/二进制/取证/恶意分析/无线 等全方向。
license: 各技能自带 LICENSE
compatibility: 仅限 CTF / 靶场 / 授权环境
---

# 网络安全技能图书馆（982 Skills）

## 0. 这是什么

本目录下有一个 **982 个专项技能**的图书馆（Claude Skills 格式，每个技能 = 一个子目录）。
它们是现成的"作战手册"：攻击路径、命令矩阵、工具用法、绕过技巧、取证流程。

**你不需要预加载任何一个** —— 用到哪个就现查现读，零上下文负担。

## 1. 标准工作流（顺序不可颠倒）

```bash
# ① 查索引（永远先做这步）：用 1-3 个英文关键词搜
grep -i "<关键词>" /app/skills/sec-skills-library/sec_skills_index.md

#   例: grep -i "kerberos" .../sec_skills_index.md
#   例: grep -i "forensic.*memory" .../sec_skills_index.md
#   例: grep -i "bypass" .../sec_skills_index.md

# ② 按索引给出的技能名读它
Read /app/skills/sec-skills-library/<技能名>/SKILL.md

# ③ 技能里可能有:
#    scripts/     → 直接 Bash 跑（先读脚本头部确认用法，别盲跑）
#    references/  → 补充资料，按需再读
#    assets/      → 模板/样例文件
```

索引格式: `- 技能名 | 一句话描述`，一行一个技能，共 982 行。
另有一个机器可读版: `/app/skills/sec-skills-library/sec_skills_index.json`（name+desc 数组）。

## 2. 覆盖方向（速查关键词）

| 方向 | 搜这些关键词 |
|---|---|
| Web 攻防 | `bypass` `waf` `sqli` `ssrf` `rce` `deserial` `jwt` `oauth` `graphql` `cors` |
| 内网/AD | `active-directory` `kerberos` `lateral` `dcsync` `bloodhound` `smb` |
| 云 | `aws` `azure` `gcp` `kubernetes` `serverless` `container` |
| 二进制/逆向 | `reverse` `exploit` `shellcode` `heap` `rop` `malware` `unpack` |
| 取证 | `forensic` `memory` `disk` `timeline` `artifact` `browser` |
| 无线/硬件 | `wifi` `bluetooth` `rfid` `firmware` `iot` |
| 社工/钓鱼 | `phishing` `social` `osint` |
| 防御/检测 | `detect` `hunting` `blue` `siem` `hardening` `incident` |

## 3. 纪律

- **查了索引再动手**：瞎猜技能名会浪费时间；grep 索引 1 秒出结果。
- **先读 SKILL.md 再执行**：技能里的命令按它的上下文用，别照搬乱打。
- **授权边界**：图书馆技能仅用于 CTF/靶场/授权测试。目标任务不在授权范围内就拒绝执行。
- 技能与技能之间常有互引（`RELATED ROUTING` 段），顺着引用能找到更对口的技能。
- 索引是启动时生成的快照；若发现新技能目录没在索引里，重跑
  `python E:\ARTEX-main\data\_ops\build_skills_index.py`（容器内路径 /app/data/_ops/build_skills_index.py 的宿主映射）。
