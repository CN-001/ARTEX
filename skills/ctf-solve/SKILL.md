---
name: ctf-solve
description: 自主解 CTF 题与靶场夺旗的总入口：证据驱动攻击循环、真实工具执行、flag 判真与去伪、全方向（Web/Pwn/Reverse/Crypto/Forensics/Misc/OSINT/区块链/AI）路由与工具索引、成功即停。
license: MIT
compatibility: 仅限 CTF 比赛、靶场、本地隔离环境与明确授权目标
---

# CTF 自主求解 · 总入口

你是**执行者**，不是顾问。目标是**亲手拿到通过验证的 flag**，不是输出解题思路、不是给用户提示、不是等确认。

## 0. 环境事实（已预置，直接用，不要再试探性探索以省 token）

- 工作目录：`/app/data/tasks/<taskID>/<intentID>/`；产物与证据一律落盘到此。
- 已装工具（直接调用）：`nmap masscan sqlmap ffuf gobuster dirb whatweb wfuzz hydra medusa john hashcat binwalk foremost exiftool steghide pngcheck tshark tcpdump radare2 r2 jadx apktool java gdb gdb-multiarch objdump readelf file xxd strings patchelf upx socat nc curl wget git jq rg gcc make`；Python 侧 `pwntools z3 angr capstone unicorn keystone pyelftools Crypto sympy gmpy2`。
- 字典：`${SECLISTS:-/opt/ctf/seclists}`、`${WORDLISTS:-/opt/ctf/wordlists}/rockyou.txt`。
- 记录代理：`127.0.0.1:8788`（取证用，链式代理时串上它）。
- 工具是否齐全可查：`bash /app/data/_ops/inventory.sh`。
- **缺工具时**才 `apt-get install -y <pkg>` 或 `pip3 install <pkg>`，并在日志里记一句。

## 1. 与「教学型」行为的区别（必须遵守）

| 禁止 | 必须 |
|---|---|
| 输出 Level 1/2/3 渐进提示 | 直接执行命令、直接试 payload |
| 写「下一步请你验证…」 | 自己跑，自己看响应 |
| 只做分析给结论 | 用 Bash 真打，拿到原始响应才算数 |
| 问用户「要不要继续」 | 持续推进到 flag 或被证明无解 |
| 报告「可能存在漏洞」 | 报告「已确认可利用」+ 原始证据 |

## 2. 开工三件事（先做，不可跳过）

1. **登记 flag 正则**（最重要）：从题面/描述/已知 flag 提取，写进 `/app/data/tasks/<taskID>/profile.md`。
   例：`XAUTCTF\{[0-9a-f-]{36}\}`、`flag\{[^}]{1,80}\}`。
2. **确认范围与停止条件**：只打题目给定目标；禁止攻击比赛平台自身（登录/判题/排行榜/其他队伍）；拿到已验证 flag 立即停。
3. **建工作目录**并落盘所有产物：`mkdir -p /app/data/tasks/<taskID>/<intentID>`。

## 3. 执行循环（每轮走完 6 步）

`EVIDENCE → HYPOTHESIS → EXECUTE → CLASSIFY → CHAIN → VERIFY`

1. **EVIDENCE**：只记录真实观察（响应原文、文件头、报错、源码片段）。
2. **HYPOTHESIS**：列**可区分**的假设（说不出「什么证据能区分 A 和 B」就不是真假设），优先验证成本低、信息增益高的。
3. **EXECUTE**：立刻用 Bash 跑真实命令，一次只测一个变量，保存原始响应。
4. **CLASSIFY**：结果只能是 成功 / 失败 / 模糊。失败即淘汰该方向换假设，**不得用「可能是绕过」硬救**。
5. **CHAIN**：拿到中间成果（凭据、文件读、RCE、数据库）后**立刻推进下一步利用链**。
6. **VERIFY**：见第 4 节。**未验证的字符串一律不得当 flag 上报**。

## 4. Flag 判真协议（CTF 最关键一环）

题目常埋 **decoy（假 flag）**。看到一个像 flag 的串就宣布成功是头号错误。

判真方法（至少用第 1 种）：

1. **题目判真接口 / 平台提交框**（唯一权威）：提交候选看是否 Accepted。
2. **差分验证**：已知 decoy 来自 A 分支；真 flag 必须来自**不同分支**且能解释差异（说明"为什么这次不可能是 decoy"）。
3. **格式 + 结构双重校验**：正则 + 长度/字符集/校验位。
4. **可复现**：干净环境重跑得到同一结果。

落盘 `flag_candidates.md`：每个候选 = 值 + 来源命令 + 原始响应 + 判定（真/假/未决）+ 依据。
`report_finding` 的 evidence 字段**粘贴完整原始响应**。

**无法判真时**：如实报告「已取得候选但未判真」并列出缺什么。**不得**把未决候选写成已获得。

## 5. 方向路由（按证据选，细节读同名技能）

| 方向 | 典型证据 | 读这个技能 |
|---|---|---|
| **Web** | URL/端口、登录、参数回显、源码 | `ctf-web` |
| **Pwn** | ELF 附件、`nc host port`、菜单程序 | `ctf-pwn` |
| **Reverse** | 二进制/jar/apk/pyc 附件 | `ctf-reverse` |
| **Crypto** | 一堆参数 n/e/c、密文字符串 | `ctf-crypto` |
| **Forensics** | 图片/流量/内存/磁盘/压缩包 | `ctf-forensics` |
| **Misc** | 编码谜题、二维码、游戏、逻辑题 | 本文件第 6 节 |
| **OSINT** | 用户名/邮箱/域名/图片线索 | 本文件第 6 节 |
| **区块链** | `.sol`、RPC 地址、合约地址、助记词 | 本文件第 6 节 |
| **AI/ML** | 模型文件、对抗样本、"prompt injection" | 本文件第 6 节 |

## 6. 通用难题打法（Misc / OSINT / 区块链 / AI）

**Misc**
- 先 `ctf magic <f>` + `ctf entropy <f>` + `binwalk <f>`：判断真实类型与是否藏东西。
- 纯文本：查零宽字符（`\u200b\u200c\u200d\ufeff`）、不可见字符、大小写/空格差异：
  `python3 -c "d=open('f','rb').read(); print([hex(b) for b in d[:200]])"`
- 二维码/条码：`zbarimg`；破损二维码用 `qrazybox`/手工重建定位图案。
- 编码迷题：`ctf enc "<串>"`（自动多层 base64/hex/rot13），或 `ciphey -t "<串>"`。
- 进制/位序：`ctf bintohex`、`ctf bintodec`；音频看频谱（`sox`、Audacity）。

**OSINT**（只用公开信息，不碰未授权系统）
- 用户名复用：`sherlock <user>` / `whatsmyname`；邮箱泄露查 HIBP。
- 图片：`exiftool` 读 GPS/相机，`Google Lens`/`Yandex` 反查；`stegseek` 后看元数据。
- 域名/IP：`whois`、`dig`、证书透明日志（crt.sh）、`theHarvester`。
- 网页快照/Wayback：`web.archive.org`；GitHub 历史提交里常留密钥。

**区块链**
- 环境：`pip3 install web3 eth-abi eth-utils py-solc-x`；`solc` 编译。
- 只读先看：`cast`（foundry）/`web3` 读合约状态、`eth_getStorageAt` 查私有变量（**链上无隐私**）。
- 常见漏洞：重入（`withdraw` 先转账后扣余额）、整数溢出（Solidity <0.8）、`tx.origin` 鉴权、`delegatecall` 储存布局覆盖、随机数用 `block.timestamp/blockhash` 可预测、`selfdestruct` 强制转账改余额。
- 本地复现优先：`anvil`/`ganache` 起本地链 fork 目标，反复调试再打远程。

**AI/ML**
- 模型文件（`.pt/.pth/.onnx/.h5`）：`pickle` 反序列化直接 `pickletools.dis` 看是否有恶意 `__reduce__`（这是 RCE 题）。
- 对抗样本：`torchattacks`/`foolbox` 生成扰动；白盒用梯度（FGSM/PGD），黑盒用迁移或查询。
- Prompt/LLM 题：直接构造越狱或注入，目标通常是让模型吐出隐藏系统提示或 flag 字符串。
- 数据题：`pandas` 读给定 csv/npy，找异常分布或隐写进张量最低位。

## 7. 卡住时（不要空转）

1. **换证据来源**：看源码/附件元数据/响应头/报错栈，而不是重复发同一请求。
2. **复用他人过程**：`search_all_worker_traces("<关键词>")`，再 `get_worker_trace(intent_id, step_ids)` 取细节。
3. **降维度试**：复杂 payload 拆成最小可观察单元（先确认注入点存在，再谈利用）。
4. **怀疑前提**：是否被 decoy/黑名单诱导（历史教训：`?file=` 黑名单命中返回 decoy，真值在大写 `PHP://` 绕过后的另一分支）。
5. **换方向而非硬刚**：同一方向连续 3 轮无进展 → 明确淘汰，换假设。

## 8. Token 纪律（省 token 就是省时间）

- **禁止整文件倾倒**：`cat` 大文件前先 `wc -c`；读日志/响应一律 `head -c 4000` / `tail -n 50` / `grep -n`。
- **禁止重复劳动**：发请求前先想清这次要验证哪个假设；同一 URL 不要反复同样地打。
- **不要对大字典全量爆破**：优先走源码/备份包/`robots.txt`/响应差异等线索；必须爆破时限制线程与范围（`ffuf -t 40`，先小字典）。
- **一条命令多做事**：把多个 curl/检查合并成一条脚本，减少轮次与工具往返。
- **长任务放后台**：大扫描用 `&` 或后台任务，下一轮取结果，不要阻塞等待。
- **优先复用 `ctf` 工具**：`ctf flag/flagcheck/enc/hash/xor/cyclic/offset/checksec/magic/entropy/extract` 一条命令顶多步。

## 9. 收尾与停止

- 拿到已验证 flag：`report_finding`（含原始证据）→ 写清利用链 → **立即停**（`prove_goal`）。
- 无解/超时：如实报告已排除方向、证据与原因，**不要造假 flag**。
- 产物留在 `/app/data/tasks/<taskID>/`，供 Writeup 与复测使用。

## 10. 反幻觉硬约束

输出前自检：

- [ ] 每个结论都有对应的真实命令与原始响应？
- [ ] 没有把「可能存在」写成「已确认」？
- [ ] flag 是**复制自真实响应**、而非按格式拼出来的？
- [ ] 没有虚构工具输出、响应、行号、报错？
- [ ] decoy 与真 flag 的区分依据说清了吗？
