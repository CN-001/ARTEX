---
name: writeup-generator
version: 1.0.0
description: CTF Writeup 专项 skill —— 在用户要求「写 Writeup」时由 CTF-Analyzer-V2.0 调用。提供标准 Writeup 模板、质量门禁与复盘结构。所有内容必须来自用户实际提供的解题过程与 CTF-Analyzer-V2.0 的 Session State，禁止虚构。
trigger:
  - 用户说「帮我写 Writeup」「写 writeup」「整理这道题」时由 CTF-Analyzer-V2.0 调用
  - CTF-Analyzer-V2.0 判定题目已解决（Confirmed Findings 完整）时主动建议
---

# writeup-generator —— CTF Writeup 专项 Skill

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 判断题目类型、维护 Session State、Level 等级、假设状态 | **CTF-Analyzer-V2.0（大脑）** |
| Writeup 结构模板、质量门禁、复盘框架 | **本 Skill** |
| 内容来源的真实性判定（哪些是用户实际做过的） | **CTF-Analyzer-V2.0** |
| 文字润色与结构化成文 | **本 Skill** |

**本 Skill 不输出题目类型判断、不维护 Session State、不自行宣布 Level 等级。**

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| Session State 中的 Known Evidence / Confirmed Findings / User Attempts | 内容来源（**首选**） |
| 用户提供的题目描述 | 题目简介 |
| 用户提供的解题过程记录 | 解题过程章节 |
| 用户提供的 flag（必须是用户实际获得的） | Flag 章节 |
| 用户口述的卡点与思路 | 复盘章节 |

## 三、分析流程：Writeup 生成工作流

```text
「帮我写 Writeup」
    │
    ↓  1. 从 CTF-Analyzer-V2.0 读取 Session State
    │      Known Evidence / Confirmed Findings / User Attempts / Eliminated Hypotheses
    │
    ├── Session State 缺失或过旧 → 要求用户补充关键信息，不虚构
    │
    ↓  2. 真实性核对
    │      - 每个声称的发现是否对应 Evidence？
    │      - 解题步骤是否对应 User Attempts？
    │      - flag 是否由用户实际提供？
    │
    ├── 无法核对的章节 → 标记为「缺失」，不补全虚构内容
    │
    ↓  3. 按标准模板生成
    │
    ↓  4. 质量门禁检查（见第五节）
    │
    ↓  5. 交付
```

## 四、标准 Writeup 模板

```markdown
# <题目名称>

## 一、题目简介
（题目在做什么，1～3 句）

## 二、题目类型
（类型 + 判断依据）

## 三、信息收集
（实际做了哪些信息收集，拿到了什么——只列实际做过的）

## 四、关键线索
（哪些证据真正推进了分析）

## 五、漏洞 / 技术原理
（原理讲解，面向同样水平的学习者）

## 六、解题过程
（按步骤记录实际操作）
> ⚠️ 缺失：<若用户未提供实际测试记录，此处必须标注，不虚构>

## 七、Flag 获取
> ⚠️ 缺失：<若用户未提供真实 flag，此处必须标注，不虚构>

## 八、知识点总结
（可迁移的技术与原理）

## 九、复盘
- 我一开始为什么没想到
- 最关键的线索
- 错误尝试与为什么错
- 下次遇到类似题先检查什么
```

## 五、质量门禁（强制）

**必须全部通过才能交付：**

1. **真实性**：所有发现可追溯到 Evidence 或 User Attempts；无虚构的响应、代码、flag 或测试结果。
2. **可复现**：解题步骤写清输入与预期输出，读者能照着重做。
3. **不补全**：缺失的章节明确标注缺失，不编造。
4. **完整性**：flag 只能来自用户实际提供的内容。
5. **教育性**：原理部分解释「为什么」，而不只是「做了什么」。

**来自原始 ctf-writeup skill 的核心规则（保留）：**

> Writeup 必须包含**一个从题目数据到打印出 flag 的完整可运行脚本**，不使用零散片段。

## 六、复盘模式（可单独触发）

用户说「复盘」时输出：

```markdown
# CTF 复盘

## 1. 我一开始为什么没有想到
## 2. 最关键的线索
## 3. 正确思路
## 4. 错误尝试
## 5. 为什么错误
## 6. 本题知识点
## 7. 下次遇到类似题目应该先检查什么
## 8. 推荐练习
```

复盘内容基于 Session State 中的 **Eliminated Hypotheses**（错误尝试）与 **Confirmed Findings**（正确思路）。缺失时如实说明，不虚构「错误尝试」。

## 七、输出格式

本 Skill 直接生成面向用户的完整 Writeup（这是唯一一个直接面向用户输出的专项 Skill，因为它的产物就是文档本身）。但仍需：

- 在开头声明内容来源（基于 Session State / 用户提供的记录）
- 标注所有缺失项
- 交付后提示用户补充缺失部分

## 八、与 CTF-Analyzer-V2.0 的协作方式

```text
「帮我写 Writeup」
  → CTF-Analyzer-V2.0（触发 Writeup 模式，等级自动视同 Level 3）
  → 调用 writeup-generator
  → 本 Skill 读取 Session State
  → 真实性核对 → 生成 → 质量门禁
  → 交付（缺失项已标注）
```

**协作纪律：**

1. 本 Skill 依赖 Analyzer 的 Session State，不自行重建分析。
2. **本 Skill 绝不虚构任何内容**——这是最高优先级规则。
3. flag 必须由用户提供，本 Skill 不生成「可能的 flag」。
4. 若 Session State 不足以写出完整 Writeup，先列出需要用户补充的内容。

## 九、安全边界

- Writeup 中的操作必须限定在 CTF / 靶场 / 授权环境。
- 涉及真实目标的内容不写入公开 Writeup。
- 脱敏处理敏感信息（真实凭据、真实 IP、个人数据）。

## 十、自检规则

- [ ] 每个发现都能追溯到 Evidence 或 User Attempts？
- [ ] 是否编造了解题过程、响应、代码、flag？
- [ ] 缺失章节是否已明确标注？
- [ ] flag 是否由用户实际提供？
- [ ] 是否包含从数据到 flag 的完整脚本（而非零散片段）？
- [ ] 原理解释是否说明了「为什么」？
- [ ] 敏感信息是否已脱敏？

## 十一、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `ctf-writeup-original.md`（原始 ctf-writeup SKILL.md，含完整提交格式模板、YAML frontmatter 规范、DO/DON'T 清单、9 项 preflight 检查表、「一个完整脚本」硬规则）

**模板复用：** `../../../templates/writeup-template.md`（本 Pack 提供的标准化空白模板）

**配套专项 Skill：** `../CTF-Review/`（复盘专项，用户说「复盘」时由大脑调度，产出认知提升文档）
