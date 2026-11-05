# SRC 提交格式（三种平台）

## ① 补天（butian.net）

```
【漏洞标题】某站SQL注入（需含资产名）
【资产】https://xxx.example.com/login
【漏洞类型】SQL注入
【危险等级】高危/中危/低危
【复现步骤】
1. 访问 https://xxx.example.com/login
2. 在 username 参数输入：' OR 1=1-- -
3. 观察返回：登录成功/报错回显
【证明截图】（登录成功界面或报错回显，带时间水印）
【原始请求】
POST /login HTTP/1.1
Host: xxx.example.com
...
【影响范围】可绕过认证/可读取数据库
【修复建议】参数化查询
【联系】手机/邮箱
```

## ② 漏洞盒子（vulbox.com）

```
【漏洞名称】xxx 站点存在任意文件读取
【URL】...
【漏洞等级】...
【漏洞描述】（背景 + 原理，2-3 段）
【复现过程】编号步骤 + 截图
【修复方案】
```

## ③ HackerOne（英文）

```
Title: [Severity] Vulnerability type at asset

Summary:
  - Brief description of the issue

Steps To Reproduce:
  1. ...
  2. ...
  3. ...

Supporting Material/References:
  - Request/response (redacted)
  - Screenshots

Impact:
  - What an attacker could do

Remediation:
  - How to fix
```

## 通用证据要求

- 截图必须**包含时间**和**目标 URL/资产标识**
- 不要上传完整数据库导出等敏感数据（打码/截断）
- 报告里**不要**写"我已获取全部用户数据"这类夸大表述
- 复测时保留原始证据，便于证明修复效果

## 常见坑

| 坑 | 后果 | 规避 |
|---|---|---|
| 报"信息泄露"但只是 robots.txt | 低价值/驳回 | 只看敏感文件/凭据类 |
| 重复漏洞 | 驳回 | 提交前搜平台历史 |
| 只给截图不给原始报文 | 审核慢 | 附 curl 命令或原始 HTTP |
| 夸大量级 | 扣分 | 按实际影响写 |
| 越界测试 | 封号 | 严格按 scope |
