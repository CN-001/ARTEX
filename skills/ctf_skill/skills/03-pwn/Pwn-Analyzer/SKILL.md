---
name: pwn-analyzer
version: 1.0.0
description: CTF Pwn 题专项 playbook —— 在 CTF-Analyzer-V2.0 判定题目为 Pwn 后被调用。提供输入点→边界检查→漏洞类型→保护机制→利用方向的完整决策树，checksec 逐项解读，以及可运行的 pwntools 模板脚本。本 Skill 是知识层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在判定题目类型为 Pwn 后调度调用
  - 用户提供二进制程序、源码、checksec 输出、反编译结果并希望分析内存安全漏洞
  - 用户要求解释保护机制（Canary / NX / PIE / RELRO）对利用方式的影响
---

# ctf-pwn —— CTF Pwn 题专项 Playbook

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 判断题目类型、维护 Session State、Level 1/2/3、假设追踪与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| 漏洞类型判定、保护机制解读、利用链构造、pwntools 脚本模板 | **本 Skill（playbook）** |
| 决定下一步验证什么、解读测试结果 | **CTF-Analyzer-V2.0** |
| 被询问利用原理时提供纵深技术内容 | **本 Skill** |

**本 Skill 不输出题目类型判断、不维护假设状态、不自行宣布 Level 等级。**

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| C 源码 | 定位输入点与缓冲区大小 |
| ELF / PE 二进制 | 反汇编 / 反编译分析 |
| `checksec` 输出 | 保护机制判定（**必须逐项解释**） |
| 反编译 / 反汇编结果 | 定位比较逻辑、后门函数 |
| 程序运行输出 / 报错 | 行为观察 |
| gdb / pwndbg 调试记录 | 偏移实测、地址泄漏 |

## 三、分析流程：Pwn 决策树

```text
程序（源码 / 二进制 / 反编译）
    │
    ↓  输入点在哪？（gets / scanf / read / fgets / 参数传递）
    │
    ├── 无证据 → Unknown（要求补充程序输入方式）
    │
    └── 找到输入点
          │
          ↓  输入是否有长度 / 边界检查？
          │
          ├── 否（无检查）→ 溢出类
          │     │
          │     ↓  栈还是堆？
          │     ├── 栈缓冲区 → 栈溢出
          │     └── 堆缓冲区 → 堆溢出
          │
          └── 是（有检查）→ 检查其他缺陷
                │
                ├── printf 等格式化函数参数可控 → 格式化字符串漏洞
                ├── 释放后仍被使用 / 重复释放 → UAF / Double Free
                └── 大小计算涉及整数运算 → 整数溢出
    │
    ↓  checksec 保护机制 + ASLR
    │
    ↓  确定利用方向
    │
    ↓  利用链所需信息是否齐全？（缺地址 / 缺偏移 / 缺 libc 版本 → Unknown）
```

### 保护机制 → 利用方向对照表

用户提供 `checksec` 时**必须逐项解释**：

| 保护 | 含义 | 利用影响 |
|---|---|---|
| **Canary** | 栈溢出保护，函数返回前检查栈上随机值 | 需先泄漏或绕过 canary 才能栈溢出 |
| **NX** | 栈不可执行 | 不能写 shellcode 到栈上执行，需 ROP |
| **PIE** | 程序基址随机化 | 需要先泄漏程序基址 |
| **RELRO** | GOT 表保护（Partial / Full） | Full RELRO 下不能改 GOT 表 |
| **ASLR** | 地址空间随机化（系统级） | 需要泄漏 libc / 程序地址 |

**保护组合 → 首选策略：**

| 组合 | 首选方向 |
|---|---|
| Partial RELRO + No PIE | GOT 表改写 |
| No canary + NX | 栈溢出 → ROP |
| No canary + No NX | 栈溢出 → shellcode |
| PIE + NX | 需先泄漏基址 → ret2libc / ROP |
| 存在后门函数 + No PIE | ret2win（最短路径） |

### ⚠️ 偏移量纪律

**覆盖偏移量必须实测确认**（如 `cyclic` / `pattern_offset`），**不要套用「buf 大小 + 8」这类公式**。对齐、padding、编译器优化都会改变实际偏移。公式计算值只作参考。

## 四、输出格式

返回给 CTF-Analyzer-V2.0 的结果（**不是面向用户的最终格式**）：

```markdown
## Pwn 专项分析结果

### 输入点与缓冲区
- 输入函数：<>
- 缓冲区大小：<>
- 边界检查：<有/无>

### 漏洞类型判定
<类型> ← <Evidence>

### 保护机制解读（若用户提供了 checksec）
| 保护 | 状态 | 影响 |
|---|---|---|

### 建议利用方向（供 Analyzer 纳入假设）
- <方向>：<技术依据>

### 利用链缺失信息
- <缺什么地址 / 偏移 / libc 版本>

### 可选验证步骤（按成本从低到高）
1. <步骤> —— 成本：<>，预期信号：<>

### 风险与不确定性
- <哪些是推断，哪些缺证据>

### 参考语料
- <相关 .md / scripts/>
```

## 五、与 CTF-Analyzer-V2.0 的协作方式

```text
用户输入
  → CTF-Analyzer-V2.0（识别类型、提取 Evidence、建立假设）
  → 判定为 Pwn
  → 调用 ctf-pwn playbook
  → 本 Skill 返回漏洞判定、保护解读、利用方向
  → Analyzer 汇总（更新假设、计算优先级、生成 Smart Next Step）
  → Level 1/2/3 输出
```

**协作纪律：**

1. 只被 Analyzer 调用，不反向调用。
2. 不修改 Session State。
3. 缺关键程序信息（位数、保护机制、源码）时返回 Unknown + 需补充清单。
4. **偏移量不给公式答案**，引导实测。
5. 未确认利用链可行性前，不宣称「可以拿到 shell」。

## 六、安全边界

- 仅用于 CTF、靶场、本地实验与授权安全研究。
- 不针对真实生产系统的二进制提供利用链。
- 越界时转为漏洞原理讲解、安全编译选项建议（ Canary / NX / RELRO 加固）与防御方案。

## 七、自检规则

- [ ] checksec 是否逐项解释并说明对利用的影响？
- [ ] 偏移量是否要求实测，而非套公式？
- [ ] 是否编造了地址、偏移量、libc 版本或执行结果？
- [ ] 保护组合与建议方向是否一致？
- [ ] 是否返回假设供 Analyzer 纳入，而非自己宣布结论？

## 八、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `overflow-basics.md`、`rop-and-shellcode.md`、`rop-advanced.md`、`format-string.md`、`heap-techniques*.md`、`heap-fsop.md`、`advanced-exploits*.md`、`kernel*.md`、`sandbox-escape.md`、`field-notes.md`、`scripts/`（5 个可运行 pwntools 模板）

**hack-skills 语料（VillanCh，MIT）：** `stack-overflow-and-rop/`、`format-string-exploitation/`（含 pwntools 自动化、GOT/hook 表、FORTIFY 绕过、盲打方法）、`heap-exploitation/`、`binary-protection-bypass/`、`arbitrary-write-to-rce/`、`kernel-exploitation/`、`sandbox-escape-techniques/`
