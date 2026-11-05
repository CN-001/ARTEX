---
name: misc-analyzer
version: 1.0.0
description: CTF Misc 题专项 playbook —— 在 CTF-Analyzer-V2.0 判定题目为 Misc 后被调用。覆盖编码识别、Pyjail/Bashjail 逃逸、沙箱逃逸、游戏与 VM、RF/SDR、DNS、OSINT、AI/ML 攻击与赛后平台操作。本 Skill 是知识层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在判定题目类型为 Misc 后调度调用
  - 用户提供混杂编码、pyjail / bashjail 等沙箱逃逸题、异常数据、OSINT 线索或 CTFd 平台相关问题
  - 用户要求分析不属于 Web/Crypto/Pwn/Reverse/Forensics 的杂项题目
---

# ctf-misc —— CTF Misc 题专项 Playbook

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 判断题目类型、维护 Session State、Level 1/2/3、假设追踪与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| 编码识别、jail 逃逸、沙箱突破、各类杂项技术、工具链 | **本 Skill（playbook）** |
| 决定下一步验证什么、解读测试结果 | **CTF-Analyzer-V2.0** |
| 被询问杂项技术原理时提供纵深技术内容 | **本 Skill** |

**本 Skill 不输出题目类型判断、不维护假设状态、不自行宣布 Level 等级。**

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| 未知编码字符串 | 编码识别与解码链 |
| 受限 shell / 沙箱环境提示 | jail 逃逸方向 |
| 游戏规则 / VM 指令集 | 游戏与 VM 分析 |
| 图片 / 音频 / 无线电信号 | RF/SDR 与信号解码 |
| DNS 记录 / 查询 | DNS 类题目 |
| 社媒截图 / 坐标 / 照片元数据 | OSINT |
| 模型文件 / 权重 / LoRA | AI/ML 攻击 |
| 赛后平台（CTFd）问题 | 平台操作 |

## 三、分析流程：Misc 决策树

```text
输入数据 / 环境描述
    │
    ↓  是否为可识别的编码或数据格式？
    │
    ├── 是 → 编码识别（Base64/32/85、Hex、URL、ROT13、Morse、核心价值观、敲击码）
    │        ↓ 解码后是否得到可读结果？
    │        ├── 是 → 检查是否为 flag 或第二层编码
    │        └── 否 → 回到格式分析
    │
    ├── 受限环境（沙箱 / jail）→ 逃逸方向
    │     ├── Python sandbox → pyjails.md（字符集受限 / 内建禁用 / AST 过滤）
    │     ├── bash 受限 shell → bashjails.md（受限 PATH / 禁用命令 / rbash）
    │     └── 容器 / 沙箱 → linux-security-bypass/、sandbox-escape-techniques/
    │
    ├── 网络与协议 → DNS 类 / 隧道与流量伪装（tunneling-and-pivoting/）
    │
    ├── 信号与硬件 → RF/SDR 解码（rf-sdr.md）、外设捕获
    │
    ├── 游戏与虚拟机 → 规则分析 / VM 指令集逆向（games-and-vms*.md）
    │
    ├── OSINT 线索 → 时间戳解码 / 地理定位 / 社媒痕迹（osint/）
    │
    ├── AI/ML 相关 → 权重扰动 / 适配器合并 / 模型反演（ai-ml/）
    │
    └── 证据不足 → Unknown
```

### 编码自动识别优先级（经验序）

ROT13 → Base64 → Base32 → Base85 → Hex → URL 编码 → Morse → 敲击码 → 核心价值观密码 → 二进制/八进制

### Misc 高频考点

| 考点 | 关键信号 | 参考 |
|---|---|---|
| 编码套娃 | 多层可解码结构 | `encodings.md`、`encodings-advanced.md` |
| Pyjail | Python 环境提示 `>>>`、黑名单报错 | `pyjails.md` |
| Bashjail | `rbash`、受限 PATH 报错 | `bashjails.md` |
| 反弹 shell / 隧道 | 需要 shell 或穿透隔离网络 | `reverse-shell-techniques/`、`tunneling-and-pivoting/` |
| 常见服务未授权访问 | 端口与服务指纹 | `unauthorized-access-common-services/` |
| 受限 shell 绕过 | 命令被拦截 | `linux-security-bypass/` |
| CTFd 平台操作 | 比赛平台 API 调用 | `ctfd-navigation.md`（含 CTFd REST API 与 Python 客户端） |
| OSINT 定位 | 照片 / 时间戳 / 坐标 | `osint/geolocation-and-media.md`（雪花时间戳、Plus Codes、What3Words、MGRS） |

## 四、输出格式

返回给 CTF-Analyzer-V2.0 的结果（**不是面向用户的最终格式**）：

```markdown
## Misc 专项分析结果

### 输入性质判定
<编码 / 受限环境 / 信号 / OSINT / 平台操作> ← <Evidence>

### 决策树命中分支
<节点> ← <Evidence>

### 建议假设（供 Analyzer 纳入）
- <假设名>：<技术依据>

### 可选处理方法（按成本从低到高）
1. <方法> —— 成本：<>，预期信号：<>

### 风险与不确定性
- <哪些是推断，哪些缺证据>

### 参考语料
- <相关 .md>
```

## 五、与 CTF-Analyzer-V2.0 的协作方式

```text
用户输入
  → CTF-Analyzer-V2.0（识别类型、提取 Evidence、建立假设）
  → 判定为 Misc
  → 调用 ctf-misc playbook
  → 本 Skill 返回输入性质判定与处理方向
  → Analyzer 汇总（更新假设、计算优先级、生成 Smart Next Step）
  → Level 1/2/3 输出
```

**协作纪律：**

1. 只被 Analyzer 调用，不反向调用。
2. 不修改 Session State。
3. 编码识别给出判定依据（字符集 / 长度 / 特征字符），不靠猜。
4. jail 逃逸类建议限定在 CTF / 本靶场环境。
5. OSINT 类分析不针对真实个人隐私挖掘。

## 六、安全边界

- 仅用于 CTF、靶场、课程实验与授权研究。
- **OSINT 类题目不得用于对真实个人的人肉搜索或隐私挖掘**——转为方法学讲解与公开数据练习。
- 反弹 shell / 隧道技术只在本地隔离实验环境演示。
- 越界时转为原理讲解、靶场搭建与防御建议。

## 七、自检规则

- [ ] 编码判定是否给出字符集 / 长度 / 特征依据？
- [ ] 是否在无证据时声称了编码类型或隐藏内容？
- [ ] 是否编造了解码结果、环境响应或 flag？
- [ ] OSINT 内容是否避开了针对真实个人的隐私挖掘？
- [ ] 是否返回假设供 Analyzer 纳入，而非自己宣布结论？

## 八、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `encodings*.md`、`pyjails.md`、`bashjails.md`、`games-and-vms*.md`、`rf-sdr.md`、`dns.md`、`ctfd-navigation.md`、`osint/`（geolocation-and-media、social-media、web-and-dns）、`ai-ml/`（adversarial-ml、llm-attacks、model-attacks）

**hack-skills 语料（VillanCh，MIT）：** `linux-security-bypass/`（受限 shell / jail 逃逸）、`tunneling-and-pivoting/`、`reverse-shell-techniques/`、`unauthorized-access-common-services/`

## 九、本地工具资产（tools/）

`tools/` 内置编码/杂项/爆破类工具（用户本机 CTF 工具集）。**命中编码套娃、Base 系列、哈希爆破等节点后，先查 `tools/README.md` 取工具。** **工具权威副本固定位于 `E:\ctf\工具，脚本\misc`（本机稳定路径）；若 `tools/` 目录缺失或被应用周期性清空，直接从该路径调用；并可运行 `python E:\ctf\restore_ctf_tools.py` 一键恢复全部工具到实时技能。**

| 场景 | 工具 | 用法 |
|---|---|---|
| Base 系列自动解码 | `tools/basecrack-v4/basecrack.py` | 单条 `python basecrack.py`（交互）；指定串 `python basecrack.py -b "…"`；多层自动 `python basecrack.py --magic`；文件批量 `-f file.txt` |
| 万能编解码/魔术操作 | `tools/CyberChef_v10.19.4/` | 浏览器打开 `CyberChef_v10.19.4.html`（离线可用） |
| 编码/加密杂项 GUI | `tools/d_safe_2.1.8.1/D_Safe_Manage.exe` | 双击运行（部分功能需管理员） |
| 哈希爆破 | `tools/hashcat-6.2.6/hashcat-6.2.6/hashcat.exe` | `hashcat.exe -m <模式> hash.txt wordlist.txt`；模式：0=MD5、100=SHA1、1400=SHA256、22000=WPA/WPA2 |
| 乱码恢复 | `tools/乱码恢复/新建 文本文档.txt` | 内含在线乱码恢复工具链接（mytju.com） |
| flag 特征扫描 | `tools/破空_flag查找3.5/破空_flag查找工具3.5.exe` | GUI 批量扫描 flag 字符串 |
| Base64 隐写解码 | `tools/ctf密码/base64隐写/b64steg.py` | `python b64steg.py` |
| GIF 帧分离 / 二维码 | `tools/ctf密码/gif分离.py`、`01转二维码.py`、`批量扫二维码.py` | Python 运行 |
