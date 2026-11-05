# CTF Session State 标准模板

> 本模板由 `skills/00-core/ctf-analyzer` 第九章定义。每次收到有效新信息后，**先更新状态，再生成回答**。

```text
=== CTF SESSION STATE ===

Challenge Type:        <Web/Crypto/Pwn/Reverse/Misc/Forensics/Linux/Unknown>
Current Level:         <1/2/3>

Known Evidence:
  E1: <证据1 —— 用户实际提供>
  E2: <证据2>
  E3: <证据3>

Unknown Information:
  U1: <未知1 —— 分析推进所必需的缺口>
  U2: <未知2>

Active Hypotheses:

  H1: <假设名称>
    Status:               活跃
    Confidence:           <高/中/低>
    Supporting Evidence:  <E1, E2>
    Missing Evidence:     <U1>
    Evidence Strength:    <High/Medium/Low>
    Verification Cost:    <High/Medium/Low>
    Information Gain:     <High/Medium/Low>
    Priority:             <High/Medium/Low>  (得分 x/9)
    Next Verification:    <具体验证方式>

  H2: <假设名称>
    Status:
    Confidence:
    Supporting Evidence:
    Missing Evidence:
    Evidence Strength:
    Verification Cost:
    Information Gain:
    Priority:
    Next Verification:

Eliminated Hypotheses:

  H3: <假设名称>
    Eliminated Because: <淘汰原因 + 触发证据>
    （编号不回收；新增假设使用 H4、H5…）

Confirmed Findings:
  F1: <已确认的发现>

User Attempts:
  A1: <用户已执行的操作及结果>

Current Best Direction: <当前优先验证方向 —— Priority 最高者>

Next Step:              <下一步目标>
```

## Priority 计算速查

| 维度 | High | Medium | Low |
|---|---|---|---|
| Evidence Strength | 3 | 2 | 1 |
| Verification Cost | **Low=3** | Medium=2 | High=1 |
| Information Gain | 3 | 2 | 1 |

总分 3～9：**8～9 → High；5～7 → Medium；3～4 → Low**

**护栏：**
1. Evidence Strength = Low → Priority 上限 Medium
2. Evidence Strength = Low 且 Information Gain = Low → Priority = Low
3. ES=High + VC=Low + IG=High → Priority 必为 High
4. 并列时 VC 低者优先；VC 也并列时 IG 高者优先

## 状态更新规则

| 场景 | 行为 |
|---|---|
| 用户说「继续」 | 不重分析；从 Current Level 和 Next Step 继续 |
| 用户提供新测试结果 | 更新 Evidence → 更新假设 → 更新置信度 → 重算 Priority → 必要时淘汰 → 新 Next Step |
| 用户提供反向证据 | 触发假设淘汰，重排其余假设优先级 |
| 用户说「重头分析」 | 清空状态，重建，回到 Level 1 |
| 上下文丢失 | 明确说明，请求最小信息（题目类型 / 进度 / 已排除方向） |
