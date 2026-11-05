# 标准分析输出模板

> 由 `skills/00-core/ctf-analyzer` 第十六章定义。这是 CTF Analyzer 面向用户的默认输出格式。

```markdown
# CTF Analysis

## 1. Challenge Summary
<1～3 句话，提炼不复述>

## 2. Challenge Type
类型：<类型>
置信度：<高/中/低>
判断依据：<具体依据>

## 3. Evidence
用户实际提供的信息：
- ……
- ……

## 4. Unknown Information
目前缺少：
- ……

## 5. Hypotheses

### H1：<名称>
支持证据：……
缺失证据：……
当前置信度：……
Evidence Strength：<High/Medium/Low>
Verification Cost：<High/Medium/Low>
Information Gain：<High/Medium/Low>
Priority：<High/Medium/Low>（含排序依据）

### H2：<名称>
（同上结构）

## 6. Current Best Direction
当前最值得关注：<一个方向>
原因：<为什么是这个，包括为什么其他方向暂时不优先>

## 7. Smart Next Step
下一步目标：……
你需要观察：
1. ……
2. ……
成功信号：……
失败信号：……
模糊信号：……

## 8. Current Hint Level
Level <1/2/3>
```

## 各等级输出差异

| 等级 | 输出内容 |
|---|---|
| **Level 1** | 摘要 / 类型 / Evidence / Unknown / 假设（精简）/ 当前方向 / Next Step / 等级 |
| **Level 2** | 增加：证据链表格 + 假设三维评分与 Priority 明细 + 具体验证方法 |
| **Level 3** | 增加：完整证据链 / 完整解题思路 / 技术原理（含术语解释）/ 验证结果解释 / Flag 获取（条件允许时）/ Writeup |

## 证据链表格模板

```markdown
### 证据链

| 观察到的信息 | 推断 | 支持程度 |
|---|---|---|
| <用户提供的真实信息> | <当前可能性> | <高/中/低> |
| <尚未验证的缺口> | <未验证> | —— |
```

## 假设淘汰记录模板

```markdown
### ❌ H<n>：<假设名称> —— 已淘汰

原假设：<假设>
新证据：<触发淘汰的证据>
结论：<H<n> 淘汰>
原因：<为什么该方向被切断>
```

## Simple Challenge Mode 压缩规则

触发条件：证据高度明确 + 无有意义的竞争假设。
可压缩：单假设（附理由）、证据链 1～2 行、Unknown 只列关键项、输出长度。
不可压缩：E/I/U 区分、Level 纪律、Session State 更新、Smart Next Step 的目标与成功信号、反幻觉规则。
