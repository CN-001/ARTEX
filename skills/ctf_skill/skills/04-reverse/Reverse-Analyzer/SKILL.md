---
name: reverse-analyzer
version: 1.0.0
description: CTF Reverse 题专项 playbook —— 在 CTF-Analyzer-V2.0 判定题目为 Reverse 后被调用。提供「先尝试低成本动态手段」的解题工作流、反混淆与反调试对策、符号执行/模拟路径。本 Skill 是知识层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在判定题目类型为 Reverse 后调度调用
  - 用户提供 CrackMe、反编译/反汇编代码、加壳程序并希望分析输入验证逻辑
  - 用户要求解释反调试、代码混淆、VM 保护或符号执行技术
---

# ctf-reverse —— CTF Reverse 题专项 Playbook

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 判断题目类型、维护 Session State、Level 1/2/3、假设追踪与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| 输入验证逻辑定位、算法识别、反混淆/反调试对策、工具链 | **本 Skill（playbook）** |
| 决定下一步验证什么、解读测试结果 | **CTF-Analyzer-V2.0** |
| 被询问逆向原理时提供纵深技术内容 | **本 Skill** |

**本 Skill 不输出题目类型判断、不维护假设状态、不自行宣布 Level 等级。**

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| 反编译 / 反汇编代码 | 定位比较逻辑 |
| 二进制文件 | 工具链分析（strings / ltrace / strace / Frida / angr） |
| 程序运行行为（输入输出对） | 黑盒推断验证规则 |
| 加壳 / 混淆特征 | 选择脱壳与反混淆策略 |
| 题目描述 | 判断是「求输入」型还是「分析算法」型 |

## 三、分析流程：Reverse 决策树

**核心原则：先用低成本的动态手段拿信息，不要一上来就全量静态反编译。**

```text
程序
    │
    ↓  是否有混淆 / 加壳 / 反调试特征？
    │
    ├── 是 → 先识别保护类型 → 脱壳 / 反调试绕过 / 反混淆
    │        （参考：anti-debugging-techniques/、code-obfuscation-deobfuscation/）
    │
    └── 否 → 直接进入分析
          │
    ↓  输入验证逻辑在哪？（比较发生在哪一步）
    │
    ├── 无证据 → Unknown（要求提供关键函数的反编译）
    │
    └── 找到比较逻辑
          │
          ↓  比较方式是什么？
          │
          ├── 与硬编码常量 / 数组逐字节比较 → 提取常量，检查输入到比较之间是否有变换
          ├── 与运行时计算结果比较 → 分析计算算法
          └── 与外部数据（文件 / 网络 / license）比较 → 分析外部依赖
          │
          ↓  输入到比较之间是否有变换？
          │
          ├── 无变换 → 直接反解输入
          ├── 简单变换（异或 / 加减 / 查表）→ 逆运算还原
          ├── 标准算法或其变形 → 识别算法（AES / 自实现 / VM 保护）
          └── 复杂控制流 / VM → 符号执行 / 模拟执行（angr / unicorn / qiling）
          │
          ↓  还原输入 → 实际输入程序验证（返回成功即确认）
```

### 低成本手段优先级（先试这些）

1. `strings` —— 找提示字符串、格式串、错误信息
2. 动态运行 + `ltrace` / `strace` —— 看系统调用与库调用序列
3. Frida hook —— 在比较点直接拦截输入输出
4. angr / 符号执行 —— 约束求解自动找路径
5. Qiling / unicorn 模拟 —— 无需真实环境运行

### 算法识别常用常量

| 常量 | 指向算法 |
|---|---|
| `0x67452301`、`0xEFCDAB89` | MD5/SHA1 初始向量 |
| `0x9E3779B9` | TEA / XXTEA |
| `0xC6EF3720`、`0x09E667F3` 等 | MD5 常量表 |
| AES S-box `0x63 0x7C 0x77 0x7B...` | AES |
| `expand 32-byte k` | ChaCha20 |
| CRC32 多项式 `0xEDB88320` | CRC |

## 四、输出格式

返回给 CTF-Analyzer-V2.0 的结果（**不是面向用户的最终格式**）：

```markdown
## Reverse 专项分析结果

### 程序功能
<输入什么、验证什么>

### 保护特征
<是否有混淆 / 加壳 / 反调试>

### 比较逻辑定位
<函数> <比较方式>

### 建议还原方向（供 Analyzer 纳入假设）
- <方向>：<技术依据>

### 可选分析手段（按成本从低到高）
1. <手段> —— 成本：<>，预期信号：<>

### 风险与不确定性
- <哪些是推断，哪些缺证据>

### 参考语料
- <相关 .md>
```

## 五、与 CTF-Analyzer-V2.0 的协作方式

```text
用户输入
  → CTF-Analyzer-V2.0（识别类型、提取 Evidence、建立假设）
  → 判定为 Reverse
  → 调用 ctf-reverse playbook
  → 本 Skill 返回验证逻辑定位与还原方向
  → Analyzer 汇总（更新假设、计算优先级、生成 Smart Next Step）
  → Level 1/2/3 输出
```

**协作纪律：**

1. 只被 Analyzer 调用，不反向调用。
2. 不修改 Session State。
3. **解释关键函数的思路，不逐行机械翻译反编译代码。**
4. 关键比较逻辑未见时返回 Unknown + 需补充清单。
5. 未实际验证前，不宣称「正确输入是 X」。

## 六、安全边界

- 仅用于 CTF、靶场、课程实验与授权研究。
- 不对真实商业软件的许可验证、 DRM 保护提供绕过帮助。
- 越界时转为逆向工程原理讲解、软件保护设计与防御建议。

## 七、自检规则

- [ ] 是否优先建议了低成本动态手段，而非一上来就全量反编译？
- [ ] 是否在解释关键函数思路，而非逐行翻译代码？
- [ ] 是否编造了反编译内容、常量或正确输入？
- [ ] 比较逻辑的证据是否明确？
- [ ] 是否返回假设供 Analyzer 纳入，而非自己宣布结论？

## 八、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `tools*.md`、`tools-dynamic.md`、`tools-emulation.md`、`tools-advanced*.md`、`patterns*.md`、`patterns-ctf*.md`、`patterns-runtime.md`、`languages*.md`、`platforms*.md`、`unicorn-emulation.md`、`anti-analysis*.md`、`field-notes.md`、`malware-analysis.md`（源自 ctf-malware：C2 / PE / .NET / 脚本混淆 / 恶意软件配置提取）

**hack-skills 语料（VillanCh，MIT）：** `anti-debugging-techniques/`、`code-obfuscation-deobfuscation/`、`symbolic-execution-tools/`、`vm-and-bytecode-reverse/`
