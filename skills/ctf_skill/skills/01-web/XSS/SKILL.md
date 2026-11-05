---
name: xss
version: 1.0.0
description: XSS 专项 playbook —— 由 CTF-Analyzer-V2.0 在 Web 题判定存在跨站脚本方向后调度。覆盖 HTML / JavaScript / 属性 / DOM 四类上下文的判定与绕过、CSP 绕过与原型链利用。本 Skill 是知识层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在判定 XSS 方向后调度调用
  - 用户提供含输入回显的源码或响应，希望分析跨站脚本
  - 用户要求分析回显上下文、过滤绕过、CSP 绕过或原型链污染
---

# xss —— XSS 专项 Playbook

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 题型识别、Session State、L1/L2/L3、假设与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| 回显上下文判定、payload 构造、过滤与 CSP 绕过 | **本 Skill** |
| 决定下一步验证什么、解读测试结果 | **CTF-Analyzer-V2.0** |

本 Skill 不输出题型判断、不维护假设状态、不宣布 Level 等级。

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| 源码中的输出点（echo / render / template） | 判定输出上下文与是否有编码 |
| 响应体中输入出现的位置 | 上下文判定 |
| 已执行 payload 与浏览器渲染结果 | 由 Analyzer 归类，本 Skill 解释 |
| CSP 响应头 | 绕过方向判定 |

## 三、分析流程：XSS 决策树

```text
输入回显到响应
    │
    ↓  回显到什么上下文？（证据：响应体中输入出现的位置）
    │
    ├── HTML 上下文（标签之间）
    │     → 标签注入：`<script>`、`<img onerror>`、`<svg onload>`
    │     → 若标签被过滤：变形（大小写、嵌套、不常见标签）
    │
    ├── JavaScript 上下文（字符串内 / 代码内）
    │     → 引号闭合 + 语句注入：`';alert(1)//`
    │     → 若引号被转义：模板字符串、反斜杠、Unicode 转义
    │
    ├── 标签属性 / URL 上下文
    │     → 属性闭合：`" onfocus=alert(1) autofocus`
    │     → 事件处理器注入；`javascript:` 伪协议
    │
    └── DOM 上下文（不经过服务器，前端 sink）
          → `location.hash`、`document.write`、`innerHTML` 等 sink
          → 判定 source 到 sink 的数据流是否过滤
    │
    ↓  输出是否经过编码 / 过滤？
    │
    ├── 实体编码（`&lt;`）→ 输出层已防护，假设被削弱
    ├── 黑名单过滤 → 变形绕过（大小写、编码、不常见标签事件）
    └── CSP（Content-Security-Policy 响应头）
          → 判定 CSP 策略：`unsafe-inline`？`nonce`？`script-src` 域白名单？
          → 绕过方向：可控的 JSONP 端点、CDN 白名单滥用、`nonce` 泄漏、`base` 标签
    │
    ↓  是否为 Node.js / 前端框架环境？
    │
    └── 是 → 考虑原型链污染（prototype pollution）作为前置链
```

## 四、输出格式

```markdown
## XSS 专项分析结果

### 回显上下文判定
<HTML / JavaScript / 属性 / DOM> ← <Evidence>

### 编码与过滤情况
<实体编码 / 黑名单 / CSP / 无> ← <Evidence>

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
CTF-Analyzer-V2.0 判定 Web 题且存在 XSS 方向
  → 调用 xss playbook
  → 本 Skill 返回上下文判定与绕过方向
  → Analyzer 汇总（更新假设、计算优先级、生成 Smart Next Step）
```

**协作纪律：** 只被 Analyzer 调用；不修改 Session State；**核心判断是「输出上下文」而非「输入了什么」**——同样的输入在不同上下文下危险性完全不同；未观察到实际渲染结果前，不宣称「XSS 可触发」。

## 六、安全边界

所有操作仅限 CTF / 靶场 / 本地实验 / 明确授权环境。XSS payload 演示只在自建环境或平台自己的 sandbox 中进行。不针对未授权真实网站提供 payload。越界时转为 XSS 原理讲解、输出编码加固建议（上下文相关编码）、CSP 部署指南。

## 七、自检规则

- [ ] 是否先判定了输出**上下文**，再谈 payload？
- [ ] 是否检查了输出层编码 / CSP（而不只是输入过滤）？
- [ ] 是否编造了浏览器渲染结果？
- [ ] 未观察到渲染结果前，是否把结论保持为「存在风险的可能性」？

## 八、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `client-side.md`、`client-side-advanced.md`、`node-and-prototype.md`（Node.js 与原型链）
**hack-skills 语料（VillanCh，MIT）：** `xss-cross-site-scripting/`
**相关跨目录引用：** `../Web-Analyzer/`（通用 Web 流程）、`../Web-Analyzer/prototype-pollution-advanced/`（原型链深度利用）、`../Web-Analyzer/waf-bypass-techniques/`（绕过技术）
