---
name: bounty-ops
description: 漏洞赏金/SRC 作业流程：范围核对、被动侦察（限速）、常见漏洞初筛、人工验证清单、按平台格式提交报告的完整闭环。用于补天/漏洞盒子/教育SRC/HackerOne 等授权范围内的众测作业。
license: MIT
compatibility: 仅限 SRC/众测平台明确授权的项目范围；禁止越界、禁止未授权测试
---

# 漏洞赏金作业流程（SRC / Bug Bounty Ops）

**核心原则：AI 只做"发现"，人做"验证与提交"。**
SRC 审核看证据质量：误报会扣分/封号，所以每个 finding 必须人工复现后才提交。

## 0. 开工前：范围核对（漏了这步全白干）

```bash
# 把平台给的范围整理成文件（放进任务工作目录）
cat > scope.txt <<'EOF'
# 项目名: xxx  SRC
# 平台: 补天
# 授权域名: *.example.com (不含 admin.example.com — 明确排除)
# 授权时间: 2026-10-01 ~ 2026-12-31
# 禁止: 任何爆破/拒绝服务/社工/物理接触
EOF
```

检查清单见 `references/scope-checklist.md`：
- 是否包含 IP 段？还是只有域名？
- 子域是否包含在 `*.example.com` 内？
- 有无"禁止测试"的路径/功能（支付、删除操作）？
- 是否要求"仅被动/低危"？

## 1. 被动侦察（先跑这步，零风险）

```bash
cd /app/data/tasks/TID/recon
# ① 子域枚举（被动源为主）
subfinder -d example.com -silent -all > subs.txt
# ② 存活与指纹
httpx -l subs.txt -silent -title -tech-detect -status-code -o httpx.txt
# ③ 历史 URL（找被遗忘的接口）
gau example.com | head -2000 > gau.txt
# ④ 爬虫（限制深度与速率，别打崩目标）
katana -u https://example.com -d 2 -rl 30 -silent > crawl.txt
```

**速率纪律（SRC 大忌）**：新目标先用 `-rl 10~30`（每秒请求数）；先看平台规则，很多平台明文规定"禁止高并发"。**不要一上来就跑 nuclei 全量模板**——先 `httpx` 看存活，再针对性跑。

## 2. 初筛（低危优先，找可交的）

```bash
# 常见低危可交项：信息泄露/敏感文件/默认页面/CORS/头缺失（看平台是否收）
nuclei -l alive.txt -tags misconfig,exposure,default-login -rl 30 -silent -o nuclei.txt
# 敏感文件探测（无破坏性）
ffuf -u https://target/FUZZ -w /opt/ctf/seclists/Discovery/Web-Content/raft-small-words.txt \
     -mc 200,301,403 -rate 30 -t 5 -o ffuf.txt
```

**初筛产出**：每条候选要写清「什么资产 + 什么漏洞 + 为什么怀疑 + 怎么验证」，交给人工。

## 3. 人工验证清单（AI 不能跳过的）

- [ ] 请求/响应完整抓包（Burp/curl 保存原始报文）
- [ ] 复现步骤 3-5 步，别人照着能复现
- [ ] 影响评估（能读到什么数据/能执行什么操作，控制在只读范围）
- [ ] 确定平台是否已收录（重复漏洞会被判"已存在"）
- [ ] 严重程度按平台/CVSS 打分

## 4. 提交报告（按平台格式）

见 `references/submission-templates.md`（补天 / 漏洞盒子 / HackerOne 三种格式）。

报告要素：
```
标题: [严重程度] 资产 - 漏洞类型（一句话）
资产: 域名/IP
复现步骤: 编号列表 + 原始请求
证据: 截图/原始 HTTP 报文（脱敏）
影响: 能做什么（不夸大不缩小）
修复建议: 具体到配置/代码位置
```

## 5. 效率提升（你的差异化）

- **ARTEX 任务模板**：把「侦察→初筛→生成待验证清单」固化成固定指令，每次换目标只改域名
- **批量跑**：用 ARTEX 的并发能力同时跑 3-5 个 SRC 项目（每个项目独立任务、独立 rate limit）
- **历史战绩**：把提交过的漏洞按平台归档，复测提醒用 ARTEX 的 retester

## 6. 合规红线

- 只测授权范围；越界一个端口都可能被平台通报
- 不做破坏性测试（删除/写入/DoSS）
- 发现的高危数据立即停止并报告，不下载不传播
- 报告内容脱敏后归档，原始证据加密保存

## 7. 参考技能（图书馆）

`grep -iE 'bug|bounty|recon|subdomain|nuclei' /app/skills/sec-skills-library/sec_skills_index.md`

## 8. 入门与工具

- **入门教程**：`E:\ARTEX-main\data\_ops\docs\SRC入门教程.md`（补天/教育SRC 注册→首次提交全流程）
- **任务模板**：`E:\ARTEX-main\data\_ops\src_task_template.md`
- **一键建任务**：`powershell -File E:\ARTEX-main\data\_ops\new_src_task.ps1 -Domain xxx.com -Platform 补天 -Rate 10`
- 敏感项目用 `-Model local`（本地模型，数据不出本机）