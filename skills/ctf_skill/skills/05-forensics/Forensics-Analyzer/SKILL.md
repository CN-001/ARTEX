---
name: forensics-analyzer
version: 1.0.0
description: CTF Forensics 题专项 playbook —— 在 CTF-Analyzer-V2.0 判定题目为 Forensics 后被调用。提供文件真实类型识别、结构完整性检查、隐藏数据提取、磁盘/内存/流量/隐写分析决策树。本 Skill 是知识层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在判定题目类型为 Forensics 后调度调用
  - 用户提供磁盘镜像、内存镜像、PCAP、图片/音频文件、日志并希望提取隐藏信息
  - 用户要求分析文件结构、元数据、隐写或异常数据
---

# ctf-forensics —— CTF Forensics 题专项 Playbook

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 判断题目类型、维护 Session State、Level 1/2/3、假设追踪与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| 文件类型与结构分析、工具链、数据提取手法、隐写检测 | **本 Skill（playbook）** |
| 决定下一步验证什么、解读测试结果 | **CTF-Analyzer-V2.0** |
| 被询问取证原理时提供纵深技术内容 | **本 Skill** |

**本 Skill 不输出题目类型判断、不维护假设状态、不自行宣布 Level 等级。**

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| 磁盘镜像 / 内存镜像 | 磁盘与内存取证 |
| PCAP / PCAPNG | 流量分析、文件传输还原 |
| 图片 / 音频 / 文档 | 隐写与元数据分析 |
| `file` / `xxd` / `binwalk` 输出 | 真实类型与结构判定 |
| 日志文件 | 行为痕迹分析 |
| 文件属性 / EXIF | 元数据分析 |

## 三、分析流程：Forensics 决策树

**核心原则：先修真实类型，再修结构，最后才做内容层分析。没看到的东西不要声称发现。**

```text
文件 / 数据
    │
    ↓  真实类型识别（file / magic number）—— 扩展名是否匹配？
    │
    ├── 不匹配 → 修正为真实类型，按真实类型重新进入决策树
    └── 匹配 → 继续
    │
    ↓  结构是否完整？（PNG 的 IEND / ZIP 的中央目录 / ELF 的节区 / JPEG 的 FFD9）
    │
    ├── 不完整 → 结构修复方向（修复后重新进入）
    └── 完整 → 继续
    │
    ↓  是否有附加 / 隐藏数据？
    │
    ├── 文件尾部（IEND / 中央目录 / FFD9 之后）有数据 → 提取追加内容
    ├── 元数据异常（EXIF / 注释字段 / 文件属性）→ 提取元数据
    └── 无 → 继续
    │
    ↓  内容层面分析（按真实文件类型选择）
    │
    ├── 图片 → LSB 隐写 / 长宽截断 / 颜色异常 / 分离通道
    ├── 音频 → 频谱图 / 摩斯电码 / DTMF / LSB
    ├── 流量（PCAP）→ 敏感数据提取 / 文件传输还原 / 协议异常
    ├── 日志 → 异常行为 / 命令执行痕迹 / 爆破痕迹 / 时间线
    ├── 压缩包 → 伪加密 / 密码 / 注释字段
    ├── 磁盘镜像 → 文件恢复 / 时间线 / 删除数据
    ├── 内存镜像 → 进程 / 注入 / 凭据提取
    └── 纯文本 → 编码识别 / 隐藏字符（零宽字符 / 不可见字符）
    │
    ↓  证据不足 → Unknown（明确指出需要什么工具输出或哪部分数据）
```

### 常见取证工具与信号

| 场景 | 工具 | 关键信号 |
|---|---|---|
| 磁盘镜像 | `ewfinfo` / `mmls` / `fls` / `icat` / `foremost` / `scalpel` / `mactime` | 时间线异常、未来日期文件、日志缺口 |
| 内存镜像 | `volatility` / `malfind` / `pslist` | 进程注入、异常进程、凭据 |
| 流量 | `tshark` / `wireshark` / `binwalk` | 文件传输、明文凭证、协议异常 |
| 隐写 | `steghide` / `zsteg` / `stegsolve` / `outguess` | LSB、通道差异、文件尾追加 |
| 文件雕刻 | `foremost` / `scalpel` / `bulk_extractor` | 按文件头恢复删除文件 |
| 元数据 | `exiftool` | GPS、作者、注释、时间戳异常 |

## 四、输出格式

返回给 CTF-Analyzer-V2.0 的结果（**不是面向用户的最终格式**）：

```markdown
## Forensics 专项分析结果

### 真实类型判定
<类型> ← <magic number / file 输出>

### 结构与隐藏数据检查
- 结构完整性：<>
- 尾部追加数据：<有/无/未检查>
- 元数据异常：<>

### 决策树命中分支
<节点> ← <Evidence>

### 建议假设（供 Analyzer 纳入）
- <假设名>：<技术依据>

### 可选分析步骤（按成本从低到高）
1. <步骤> —— 成本：<>，预期信号：<>

### 风险与不确定性
- <哪些是推断，哪些缺证据>

### 参考语料
- <相关 .md>
```

## 五、与 CTF-Analyzer-V2.0 的协作方式

```text
用户输入
  → CTF-Analyzer-V2.0（识别类型、提取 Evidence、建立假设）
  → 判定为 Forensics
  → 调用 ctf-forensics playbook
  → 本 Skill 返回类型判定与提取方向
  → Analyzer 汇总（更新假设、计算优先级、生成 Smart Next Step）
  → Level 1/2/3 输出
```

**协作纪律：**

1. 只被 Analyzer 调用，不反向调用。
2. 不修改 Session State。
3. **没有实际看到的文件内容，不声称发现了隐藏数据。**
4. 隐写类假设在无证据时保持低置信（图片打不开时不应优先谈 LSB）。
5. 需要原始文件或工具输出而未提供时返回 Unknown + 清单。

## 六、安全边界

- 仅用于 CTF、靶场、课程实验与授权取证分析。
- 不对他人所有的设备、账户、真实数据进行取证分析或凭据提取。
- 越界时转为取证原理讲解、靶场复现、企业取证规范与数据保护建议。
- 保留只读分析纪律（不修改原始证据的完整性）。

## 七、自检规则

- [ ] 真实类型是否用 magic number / file 输出确认，而非仅凭扩展名？
- [ ] 是否在没看到文件内容的情况下声称发现了隐写或隐藏数据？
- [ ] 结构检查的顺序是否正确（类型 → 结构 → 隐藏 → 内容）？
- [ ] 是否编造了工具输出或提取结果？
- [ ] 是否返回假设供 Analyzer 纳入，而非自己宣布结论？

## 八、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `disk-and-memory.md`、`disk-advanced.md`、`disk-recovery.md`、`network.md`、`network-advanced.md`、`steganography.md`、`stego-advanced*.md`、`stego-image.md`、`windows.md`、`linux-forensics.md`、`signals-and-hardware.md`、`peripheral-capture.md`、`3d-printing.md`

**hack-skills 语料（VillanCh，MIT）：** `memory-forensics-volatility/`、`steganography-techniques/`、`traffic-analysis-pcap/`

**cybersecurity-skills 语料（Bri Russell，MIT）：** `disk-forensics-governance.md`（企业级取证工作流与证据保管链规范，强调只读挂载与 CTF 场景契合）、`secrets-and-credential-hunting.md`（git 历史与捆绑包中的凭据猎取）

## 九、本地工具资产（tools/）

`tools/` 内置取证/隐写/压缩包类工具（用户本机 CTF 工具集）。**按决策树命中节点后，先查 `tools/README.md` 取对应工具。** **工具权威副本固定位于 `E:\ctf\工具，脚本\misc`（本机稳定路径）；若 `tools/` 目录缺失或被应用周期性清空，直接从该路径调用；并可运行 `python E:\ctf\restore_ctf_tools.py` 一键恢复全部工具到实时技能。**

| 场景 | 工具 | 用法 |
|---|---|---|
| 十六进制/文件结构查看 | `tools/010EditorWin64Portable14.0-二次修整版/010EditorPortable.exe` | 双击运行（内置 PNG/ZIP/PDF/EXE 模板） |
| LSB 隐写 | `tools/LSB-Steganography-master/LSBSteg.py` | `encode -i 载体 -o 输出 -f 文件`；`decode -i 载体 -o 输出` |
| 双图盲水印 | `tools/双图盲水印/bwmforpy3.py` | `python bwmforpy3.py decode file.jpg file0.jpg out.jpg` |
| JPEG 隐写检测 | `tools/检测jpeg文件隐写/` | `.\Stegdetect.exe -tjopi -s 10.0 test.jpg` |
| MP3 隐写（MP3Stego） | `tools/MP3加密隐写/MP3Stego/Decode.exe` | `Decode.exe -X -P 密码 x.mp3`（输出 x.mp3.txt） |
| 音频摩斯自动分析 | `tools/全自动分析音频的摩斯/main.py` | `python main.py 文件名.wav` |
| PDF/文本隐写 | `tools/wbs43open(pdf隐写)/wbStego4.3open.exe` | 双击按向导操作 |
| PCAPNG 流量分析 | `tools/pcapngscan-main/main.py` | `python main.py x.pcapng`（统计/拓扑/HTTP 文件提取/异常检测） |
| 路由备份提取账号 | `tools/.bin结尾，宽带流量泄露(routepassview)/` | RouterPassView.exe 打开 .bin/.cfg 备份 |
| PNG 宽高爆破 | `tools/png宽高爆破全自动/p.py` | `python p.py -i x.png` |
| 压缩包 CRC 已知明文爆破 | `tools/压缩包CRC碰撞/` | 1-4字节 `python CRC32-Tools.py -h x.zip`；5字节 `5字节/CRC.py`；6字节 `6字节/crc32.py` |
| ZIP/RAR 密码爆破 | `tools/压缩包zip或rar爆破/ARCHPR.exe` | 双击运行；另 `tools/Ziperello（压缩包解密工具）/Ziperello/Ziperello.exe` |
| 压缩包修复/伪加密 | `tools/压缩包修复（损坏的压缩包并破解伪加密）/ZipCenOp.jar` | `java -jar ZipCenOp.jar r x.zip` |
| 套娃压缩包 | `tools/压缩包套娃/` | 改脚本内压缩包名后 `python 1.py` / `2.py` |
| .kdbx 保险库 | `tools/KeePass-2.57(结尾是kdbx)/KeePass.exe` | 双击打开 kdbx 文件 |
| .dat 文件分析 | `tools/分析.dat文件/二选一/超级快速解/WpfApp1.exe` | GUI 快速分析 |
