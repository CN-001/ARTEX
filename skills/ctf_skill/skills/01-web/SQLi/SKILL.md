---
name: sqli
version: 1.0.0
description: SQL 注入专项 playbook —— 由 CTF-Analyzer-V2.0 在 Web 题判定存在注入方向后调度。覆盖字符型/数字型判定、数据库指纹识别、过滤绕过、盲注与 NoSQL/EL 注入。本 Skill 是知识层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在判定 SQL 注入方向后调度调用
  - 用户提供含数据库交互的源码、SQL 报错信息、或注入测试响应
  - 用户要求分析注入类型、绕过过滤或构造 payload
---

# sqli —— SQL 注入专项 Playbook

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 题型识别、Session State、L1/L2/L3、假设与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| 注入类型判定、指纹识别、绕过技术、payload 构造、盲注方法 | **本 Skill** |
| 决定下一步验证什么、解读测试结果 | **CTF-Analyzer-V2.0** |

本 Skill 不输出题型判断、不维护假设状态、不宣布 Level 等级。

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| 源码中的 SQL 拼接 / ORM 调用 | 判断输入是否进入查询语法层 |
| 数据库报错信息 | 指纹识别与报错注入方向 |
| 正常/异常响应对比 | 判定闭合方式与过滤边界 |
| 已执行 payload 与结果 | 由 Analyzer 归类，本 Skill 解释技术含义 |

## 三、分析流程：注入决策树

```text
输入点进入数据库交互
    │
    ↓  输入是否进入 SQL 语法层？
    │   （证据：源码拼接 / 报错特征 / 响应随输入语法变化）
    │
    ├── 否（参数化查询）→ 无注入 → 通知 Analyzer 该方向淘汰
    │
    └── 是
        │
        ↓  注入类型判定
        │
        ├── 数字型（无引号闭合）→ `id=1 AND 1=1` / `AND 1=2` 响应差异判定
        ├── 字符型（引号闭合）→ 判定闭合符 `'` / `"` / `')` 等
        └── 模板/NoSQL（非 SQL）→ 转NosQL/EL 注入方向
        │
        ↓  是否有回显？
        │
        ├── 有回显 → 联合查询（UNION）方向：判定列数与显示位
        ├── 有报错 → 报错注入方向（extractvalue / updatexml / 类型转换）
        └── 无回显无报错 → 盲注方向
              ├── 布尔盲注：构造真/假条件，观察响应差异
              └── 时间盲注：`SLEEP()` / `WAITFOR DELAY` 观察响应延迟
        │
        ↓  是否存在过滤 / WAF？
        │
        ├── 大小写混合、注释符（`/**/`、`--`、`#`）
        ├── 编码绕过（URL 双重编码、Hex、Unicode）
        ├── 等价函数（`MID`↔`SUBSTR`、`ASCII`↔`ORD`）
        └── 关键字拆分 / 注释内嵌
```

### 数据库指纹速查

| 特征 | 数据库 |
|---|---|
| 报错含 `You have an error in your SQL syntax` | MySQL |
| `ORA-` 前缀 | Oracle |
| `SQLSTATE` / `PG_` 前缀 | PostgreSQL |
| `Unclosed quotation mark` | MSSQL |
| `SQLite3::query` / `SQLITE_ERROR` | SQLite |

## 四、输出格式

```markdown
## SQLi 专项分析结果

### 注入可行性判定
<可注入 / 不可注入> ← <Evidence>

### 注入类型判定
<类型> ← <闭合方式 / 回显情况>

### 建议假设（供 Analyzer 纳入）
- <假设名>：<技术依据>

### 可选验证步骤（按成本从低到高）
1. <步骤> —— 成本：<>，预期信号：<>

### 风险与不确定性
- <哪些是推断，哪些缺证据>

### 参考语料
- <本目录下相关文件>
```

## 五、与 CTF-Analyzer-V2.0 的协作方式

```text
CTF-Analyzer-V2.0 判定 Web 题且存在注入方向
  → 调用 sqli playbook（传入 Evidence 与当前假设）
  → 本 Skill 返回注入判定与验证方法
  → Analyzer 汇总（更新假设、计算优先级、生成 Smart Next Step）
```

**协作纪律：** 只被 Analyzer 调用，不反向调用；不修改 Session State；源码显示参数化查询时**明确返回「该方向淘汰」**，不用「二次注入」等话术救援；未验证前不宣称「存在注入」。

## 六、安全边界

所有操作仅限 CTF / 靶场 / 本地实验 / 明确授权环境。不针对未授权真实数据库提供 payload。越界时转为注入原理讲解、参数化查询加固建议与授权测试流程。

## 七、自检规则

- [ ] 是否区分了「源码可见拼接」（风险）与「已验证可注入」（结论）？
- [ ] 注入类型判定是否有 Evidence 支撑？
- [ ] 是否编造了数据库响应或 payload 执行结果？
- [ ] 参数化查询证据出现时，是否明确返回了淘汰信号？

## 八、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `sql-injection.md`
**hack-skills 语料（VillanCh，MIT）：** `sqli-sql-injection/`、`nosql-injection/`（NoSQL 注入）、`expression-language-injection/`（EL 表达式注入）
**相关跨目录引用：** `../Web-Analyzer/`（通用 Web 流程、fuzzer 工具链）
