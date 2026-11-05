---
name: linux-security
version: 1.0.0
description: CTF Linux 环境专项 playbook —— 在 CTF-Analyzer-V2.0 判定题目涉及 Linux 提权、后渗透、受限环境逃逸时被调用。提供从信息收集到具体提权路径的决策树与检查清单。本 Skill 是知识层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在题目涉及 Linux 提权 / 后渗透 / 受限 shell 逃逸时调度调用
  - 用户提供 Linux shell 访问、`id` / `sudo -l` / `uname -a` 输出并希望提权
  - 用户要求分析 SUID、Cron、PATH 劫持、能力（capabilities）等提权向量
---

# ctf-linux —— CTF Linux 环境专项 Playbook

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 判断题目类型、维护 Session State、Level 1/2/3、假设追踪与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| 提权向量枚举、具体提权技术、检查清单、工具链 | **本 Skill（playbook）** |
| 决定下一步验证什么、解读测试结果 | **CTF-Analyzer-V2.0** |
| 被询问提权原理时提供纵深技术内容 | **本 Skill** |

**本 Skill 不输出题目类型判断、不维护假设状态、不自行宣布 Level 等级。**

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| shell 会话与当前用户信息（`id`、`whoami`） | 权限基线 |
| `sudo -l` 输出 | sudo 滥用向量 |
| `uname -a` / 发行版信息 | 内核漏洞匹配 |
| SUID / SGID 文件清单 | SUID 滥用向量 |
| Cron 任务清单 | 计划任务利用 |
| 文件系统 / 服务清单 | 配置缺陷发现 |
| 受限 shell 报错 | jail 逃逸方向 |

## 三、分析流程：Linux 提权决策树

**核心原则：先枚举，后定向。没有枚举结果就不指定提权路径。**

```text
获得低权限 shell
    │
    ↓  环境基线收集
    │   id / whoami / sudo -l / uname -a / 发行版 / 网络接口 / 进程
    │
    ├── 信息不足 → Unknown（要求提供上述输出）
    │
    └── 有基线信息
          │
          ↓  逐层检查提权向量（按命中证据的强度排序）
          │
          ├── sudo 滥用 → sudo -l 命中可执行程序（GTFOBins 查询）
          ├── SUID / SGID 滥用 → 可疑 SUID 二进制
          ├── 内核漏洞 → 发行版与内核版本匹配已知 CVE
          ├── Cron 任务 → 可写脚本 / PATH 劫持 / 通配符注入
          ├── 服务配置 → 弱配置 / 硬编码凭据 / 可写服务文件
          ├── 文件系统 → 可写敏感文件 / SSH 密钥 / 备份文件
          ├── capabilities → 非特权但高危的 cap（如 cap_setuid）
          ├── PATH 劫持 → 可写目录在 PATH 前部 / 脚本调用无绝对路径的命令
          └── 受限 shell → rbash / 受限字符集 / 禁用命令（转 linux-security-bypass/）
```

### 常见向量速查

| 向量 | 检查命令 | 利用提示 |
|---|---|---|
| sudo | `sudo -l` | 命中程序查 GTFOBins（如 `find`、`vim`、`python3` 的 sudo 提权方式） |
| SUID | `find / -perm -4000 2>/dev/null` | 命中程序查 GTFOBins |
| 内核 | `uname -a`、`cat /etc/os-release` | 匹配 Dirty Pipe / Dirty COW 等 |
| Cron | `ls -la /etc/cron*`、`cat /etc/crontab` | 脚本可写？PATH 可控？ |
| Capabilities | `getcap -r / 2>/dev/null` | `cap_setuid` 类直接提权 |
| PATH | `echo $PATH` | 脚本中无绝对路径的命令可劫持 |
| 凭据 | 配置文件、历史记录、备份 | 硬编码密码 / 密钥复用 |

## 四、输出格式

返回给 CTF-Analyzer-V2.0 的结果（**不是面向用户的最终格式**）：

```markdown
## Linux 专项分析结果

### 环境基线摘要
- 用户 / 组：<>
- 内核 / 发行版：<>
- sudo 权限：<>

### 命中的提权向量
- <向量> ← <Evidence>

### 建议提权方向（供 Analyzer 纳入假设）
- <方向>：<技术依据>

### 可选验证步骤（按成本从低到高）
1. <步骤> —— 成本：<>，预期信号：<>

### 缺失信息
- <还需哪些枚举输出>

### 风险与不确定性
- <哪些是推断，哪些缺证据>

### 参考语料
- <相关目录 / .md>
```

## 五、与 CTF-Analyzer-V2.0 的协作方式

```text
用户输入
  → CTF-Analyzer-V2.0（识别场景、提取 Evidence、建立假设）
  → 判定涉及 Linux 提权 / 后渗透
  → 调用 ctf-linux playbook
  → 本 Skill 返回向量枚举结果与提权方向
  → Analyzer 汇总（更新假设、计算优先级、生成 Smart Next Step）
  → Level 1/2/3 输出
```

**协作纪律：**

1. 只被 Analyzer 调用，不反向调用。
2. 不修改 Session State。
3. **没有枚举证据时不指定具体提权路径**——先要求提供基线信息。
4. 提权建议给出验证方式与预期信号，不宣称「一定能提权」。
5. 涉及命令必须限定在 CTF / 本地靶场 / 授权环境。

## 六、安全边界

- 仅用于 CTF、靶场（HTB / THM / VulnHub）、本地实验与**明确授权**的测试。
- **不针对未授权的真实服务器、企业内网、他人账户提供提权或后渗透步骤。**
- 越界时转为提权原理解释、企业加固建议（最小权限、SUID 审计、sudo 策略）、授权测试流程说明。

## 七、自检规则

- [ ] 是否在缺少枚举信息时就断言了提权路径？
- [ ] 每个建议向量是否有用户提供的 Evidence 支撑？
- [ ] 是否编造了 sudo / SUID / cron 输出或执行结果？
- [ ] 建议命令是否限定在授权环境内？
- [ ] 是否返回假设供 Analyzer 纳入，而非自己宣布结论？

## 八、参考语料索引

**hack-skills 语料（VillanCh，MIT）：** `linux-privilege-escalation/`（系统性提权枚举与向量）、`linux-lateral-movement/`（后渗透与横向移动）

**相关跨目录引用：** `../../06-misc/Misc-Analyzer/linux-security-bypass/`（受限 shell / jail 逃逸）、`../../06-misc/Misc-Analyzer/pyjails.md`、`../../06-misc/Misc-Analyzer/bashjails.md`、`../../05-forensics/Forensics-Analyzer/linux-forensics.md`（取证视角的 Linux 痕迹分析）
