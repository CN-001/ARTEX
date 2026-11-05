---
name: web-analyzer
version: 1.0.0
description: CTF Web 题专项 playbook —— 在 CTF-Analyzer-V2.0 判定题目为 Web 后被调用，提供可执行的漏洞决策树、payload 索引与验证方法。本 Skill 是知识层与工具层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在判定题目类型为 Web 后调度调用
  - 用户明确要求分析 Web 类漏洞方向（认证 / 授权 / 反序列化 / 文件上传 / 命令注入 / 模板注入等）
  - 用户提供 HTTP 请求/响应、Web 源码、Burp 抓包并希望深入某个具体漏洞方向
---

# web-analyzer —— CTF Web 题专项 Playbook

## 一、核心功能

本 Skill 是 **CTF-Analyzer-V2.0 的 Web 知识与工具层**，职责边界严格如下：

| 职责 | 归属 |
|---|---|
| 判断题目类型、维护 Session State、Level 1/2/3 渐进提示、假设追踪与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| Web 漏洞的具体技术细节、payload 构造、验证手法、工具命令、绕过技巧 | **本 Skill（playbook）** |
| 决定下一步验证什么、如何解读用户返回的测试结果 | **CTF-Analyzer-V2.0** |
| 被询问「这个漏洞的原理和利用方式」时提供纵深技术内容 | **本 Skill** |

**本 Skill 不输出题目类型判断、不维护假设状态、不自行宣布 Level 等级。** 所有面向用户的最终回答由 CTF-Analyzer-V2.0 汇总生成。

## 二、输入识别

本 Skill 处理以下输入（由 CTF-Analyzer-V2.0 传入或用户直接补充）：

| 输入类型 | 用途 |
|---|---|
| HTTP 请求（方法/路径/参数/Header/Cookie/Body） | 定位输入点 |
| HTTP 响应（状态码/响应头/响应体） | 判断回显与上下文 |
| Web 源码（PHP / Python / Java / JS / Go） | 判断输入如何进入后端 |
| Burp / 抓包记录 | 还原完整交互 |
| 用户已执行的测试与结果 | 由 Analyzer 归类，本 Skill 负责解释技术含义 |
| 前端代码 / 模板文件 | 判断客户端与模板上下文 |

## 三、分析流程：Web 漏洞决策树

**走哪条分支必须由 Evidence 决定。没有证据时进入 Unknown，不要自行假设后端实现。**

```text
输入点（参数 / Header / Body / Cookie / URL 路径 / 文件上传）
    │
    ↓  是否有证据表明输入进入后端处理？
    │
    ├── 无证据 → Unknown（要求补充后端信息或响应）
    │
    └── 是
        │
        ↓  输入是否回显到响应？
        │
        ├── 是
        │     │
        │     ↓  回显到什么上下文？（证据：响应体中输入出现的位置）
        │     │
        │     ├── HTML 上下文 → XSS（HTML 注入）
        │     ├── JavaScript 上下文 → XSS（JS / DOM 注入）
        │     ├── 标签属性 / URL 上下文 → XSS（属性注入）
        │     ├── 模板语法上下文（{{ }}、${ }、<%= %>）→ SSTI
        │     └── SQL 报错信息 → SQL 注入（报错注入方向）
        │
        └── 否（或无回显信息）
              │
              ↓  有证据表明后端把输入用在了哪里？
              │
              ├── 无证据 → Unknown
              │
              ├── 数据库相关 → SQL Injection / NoSQL Injection
              │     └── 有过滤 → 绕过方向（大小写/编码/注释/等价函数）
              ├── 文件相关 → 文件上传 / 路径穿越 / 文件包含
              ├── 发起网络请求 → SSRF
              ├── 拼接系统命令 → 命令注入 / 代码执行
              ├── 反序列化入口 → 反序列化利用
              └── 用于身份/权限判断 → 认证绕过 / IDOR / 访问控制 / 逻辑漏洞
```

### 常见 Web 漏洞快速判定表

| 漏洞 | 关键证据信号 | 首选验证方式 | 参考位置 |
|---|---|---|---|
| SQL 注入 | 输入直接拼接进查询；报错含数据库特征 | 字符型/数字型判定 + 闭合测试 | 专项 Skill `../SQLi/` |
| XSS | 输入回显到 HTML/JS 上下文 | 上下文判定 + 标签探针 | 专项 Skill `../XSS/` |
| SSRF | 参数为 URL 且服务端发起请求 | 协议与内网地址探测 | 专项 Skill `../SSRF/` |
| SSTI | 回显位置出现模板语法 | 模板探针（如 `{{7*7}}`） | `server-side.md`、`server-side-advanced*.md` |
| IDOR | 参数变化导致返回对象变化 | 水平/垂直越权对比 | `auth-and-access*.md`、`auth-jwt.md`、`idor-broken-object-authorization/` |
| 命令注入 | 输入拼接到系统命令 | 无回显用延时/OOB 探针 | `server-side-exec*.md`、`cmdi-command-injection/` |
| 反序列化 | 传输序列化对象（`O:`、`¬í`、pickle、fastjson） | 魔术头识别 + gadget 链 | `server-side-deser.md`、`deserialization-insecure/` |
| 文件上传 | 存在上传接口 | 后缀/内容/MIME 绕过矩阵 | `upload-insecure-files/`、`path-traversal-lfi/` |
| 逻辑漏洞 | 业务流程可被重放/乱序/负数 | 竞态与边界值测试 | `business-logic-vulnerabilities/`、`race-condition/`、`type-juggling/` |
| WAF 绕过 | 特殊字符被拦截 | 编码/等价/分块/协议层绕过 | `waf-bypass-techniques/`、`request-smuggling/` |
| XXE | 接收 XML 输入 | 实体探测 | `xxe-xml-external-entity/` |
| 表达式注入 | EL / OGNL / SpEL 特征 | 表达式探针 | `../SQLi/expression-language-injection/` |
| JNDI 注入 | Java + 可控 lookup 串 | 协议端点探测 | `jndi-injection/` |

## 四、输出格式

本 Skill 返回给 CTF-Analyzer-V2.0 的结果采用以下结构（**不是面向用户的最终格式**）：

```markdown
## Web 专项分析结果

### 涉及输入点
- <参数/位置>

### 决策树命中分支
<走到哪个节点，依据哪条 Evidence>

### 建议假设（供 Analyzer 纳入假设系统）
- <假设名>：<技术依据>

### 可选验证方法（按成本从低到高）
1. <方法> —— 成本：<低/中/高>，预期信号：<成功标志>
2. ...

### 风险与不确定性
- <哪些结论是推断，哪些缺乏证据>

### 参考语料
- <本目录下相关 .md 文件名>
```

## 五、与 CTF-Analyzer-V2.0 的协作方式

**调用关系（严格单向）：**

```text
用户输入
  → CTF-Analyzer-V2.0（识别类型、提取 Evidence、建立假设）
  → 判定为 Web
  → 调用 web-analyzer playbook（传入 Evidence 与当前假设）
  → web-analyzer 返回技术分析结果（上述格式）
  → CTF-Analyzer-V2.0 汇总（更新假设、计算优先级、生成 Smart Next Step）
  → 面向用户的 Level 1/2/3 输出
```

**协作纪律：**

1. 本 Skill **只被 Analyzer 调用，不反向调用 Analyzer**。
2. 本 Skill **不修改 Session State**——状态只由 Analyzer 维护。
3. 本 Skill 返回的每个技术判断都标注「这是推断还是证据」。
4. 当 Evidence 不足以走决策树时，本 Skill **返回 Unknown 并列出所需信息**，不编造后端实现。
5. 用户直接对本 Skill 说「继续/给我完整解析」时，本 Skill 应提示由 Analyzer 统一控制提示等级。

## 六、安全边界

**本 Skill 的所有操作指令仅适用于：**

- CTF 比赛与练习平台
- Hack The Box / TryHackMe / DVWA / OWASP Juice Shop / pikachu 等靶场
- 本地实验环境
- 用户**明确授权**的测试目标

**禁止：**

- 针对未授权真实网站、服务器、账号的攻击步骤
- 提供针对特定真实目标的定制化攻击 payload
- 大规模扫描、爆破真实服务的具体操作指令

**越界时的行为：** 转为漏洞原理解释、CTF 模拟题构造、本地靶场搭建方案、防御与安全编码建议。

## 七、自检规则

输出前检查：

- [ ] 决策树每个分支的走法是否有 Evidence 支撑？无证据的分支是否进入了 Unknown？
- [ ] 是否把「存在风险」写成了「已确认可利用」？
- [ ] 是否编造了服务器响应、源码或 payload 执行结果？
- [ ] 验证方法是否按成本排序？是否给出了明确的成功/失败信号？
- [ ] 是否返回了假设供 Analyzer 纳入，而不是自己宣布结论？
- [ ] 涉及的操作是否限定在授权环境内？

## 八、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `auth-and-access*.md`、`auth-infra.md`、`auth-jwt.md`、`cves.md`、`field-notes.md`、`pat-reference.md`（PayloadsAllTheThings 按需索引）、`python-requests.md`（Burp-Intruder 级 fuzzer 工具链）、`server-side*.md`、`web3.md`、`scripts/`

**hack-skills 语料（VillanCh，MIT）：** `cmdi-command-injection/`、`deserialization-insecure/`、`path-traversal-lfi/`、`upload-insecure-files/`、`prototype-pollution-advanced/`、`waf-bypass-techniques/`、`business-logic-vulnerabilities/`、`race-condition/`、`type-juggling/`、`idor-broken-object-authorization/`、`xxe-xml-external-entity/`、`jndi-injection/`、`request-smuggling/`、`ssti-server-side-template-injection/`

**拆分出去的专项 Skill：** `../SQLi/`（SQL 注入）、`../XSS/`（跨站脚本）、`../SSRF/`（服务端请求伪造）

**首次快速流程：** 先映射信任边界再 fuzz；先验证输入是否进入后端，再谈漏洞类型。
