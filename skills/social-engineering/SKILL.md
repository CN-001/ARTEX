---
name: social-engineering
description: 社工（社会工程学）作战能力：目标画像与 OSINT、钓鱼邮件/短信战役、Vishing 电话社工、凭据收集、物理社工与投递载体、检测规避与报告。用于授权红队/攻防演练/企业钓鱼演练，全程要求明确书面授权与范围。
license: MIT
compatibility: 仅限明确授权的红队演练、企业钓鱼演练、CTF/攻防题场景；严禁用于未授权真实人员
---

# 社工（社会工程学）作战

**前置铁律（不满足只做被动收集）**：
1. 任务 scope 里必须有明确授权：组织名、域名、可投递的人员范围、时间窗、禁止项（如禁止真实投递、禁止物理接触）。
2. 投递类操作（真实发邮件/发短信/打电话/寄 U 盘）**必须有客户/赛方书面授权**，且提前约定"演练标识"（如邮件页脚水印、专用演练域名）。
3. 收集到的个人信息只在任务内使用，不进任何云端 API（敏感任务切本地模型）；任务结束按约定销毁。

## 0. 图书馆路由（深度资料在这里）

先用 grep 查 `/app/skills/sec-skills-library/sec_skills_index.md`，以下技能已预置：

| 需求 | 技能 |
|---|---|
| 完整社工渗透（phish/vish/smish/物理） | `conducting-social-engineering-penetration-test` |
| 定向钓鱼模拟战役 | `conducting-spearphishing-simulation-campaign` |
| 批量钓鱼演练执行 | `executing-phishing-simulation-campaign` |
| 电话社工（vishing）剧本 | `conducting-social-engineering-pretext-call` |
| 组织外部资产 OSINT | `conducting-external-reconnaissance-with-osint` |
| OSINT 收集方法 | `collecting-open-source-intelligence` |
| 邮件头分析与伪造判定 | `analyzing-email-headers-for-phishing-investigation` |
| 相似域名（typosquat）生成 | `analyzing-typosquatting-domains-with-dnstwist` |
| 二维码钓鱼（quishing） | `detecting-qr-code-phishing-with-email-security` |
| 钓鱼应急响应（蓝队视角，攻防战防守用） | `conducting-phishing-incident-response` |

## 1. 目标画像（OSINT，只做被动）

```bash
mkdir -p /app/data/tasks/TID/se && cd /app/data/tasks/TID/se
# ① 组织外部资产
subfinder -d example.com -silent -all > subs.txt
httpx -l subs.txt -silent -title -tech-detect -status-code > httpx.txt
# ② 邮箱/人员线索（theHarvester 走搜索引擎，注意限速）
theHarvester -d example.com -b bing,duckduckgo -l 200 -f harvest
# ③ 技术栈与供应商（登录页、邮箱服务商、OA/HR 系统版本）
grep -iE 'outlook|exchange|zimbra|confluence|jira|oa|hr' httpx.txt
# ④ 泄露检查（只查公开泄露，需账号/API 的经授权）
#   https://haveibeenpwned.com/API/v3  （需 API key）
```

产出：**组织画像表**（部门结构调整、员工姓名/职位/邮箱格式、上下游供应商、常用内部系统、活跃社交平台）。

邮箱格式推断：`python /app/skills/social-engineering/scripts/osint_email_format.py --domain example.com --names 张三,李四`

## 2. 钓鱼战役搭建（授权演练）

```
① 演练域名: 注册相似域名（公司名+常见后缀混写），必须挂"演练"标识页
② 邮件服务: 用独立 SMTP（演练专用），配置 SPF/DKIM/DMARC（否则进垃圾箱，也测不准）
③ 投递平台: Gophish（开源，自带统计/落地页/邮件模板）
   - 装法: 下载 gophish release zip → ./gophish → 后台 https://localhost:3333
   - 管理后台只监听内网，暴露到公网前必须改密码+加认证
④ 落地页: 仿登录页（凭证收集仅演练域，页面顶部标注演练）+ 跟踪像素
⑤ 统计: 打开率/点击率/提交率/报告率，按部门聚合（不给个人点名，除非约定）
```

**邮件模板**见 `references/phishing-templates.md`（IT 通知/密码过期/供应商发票/HR 通知/会议邀请 5 类）。
**投递前自检**：SPF/DKIM/DMARC 是否配置正确（`mail-tester.com` 打分 ≥8），邮件正文是否含演练标识，落地页是否可关停。

## 3. Vishing 电话社工（授权）

剧本结构（详见 `references/vishing-scripts.md`）：
```
开场（验证身份钩子）→ 建立紧迫感 → 索取信息/指引操作 → 收尾（防怀疑）
原则: 不录音冒充真人领导；不索要超出范围的凭据；必须准备"演练结束声明"话术
```

## 4. 凭据收集与后利用（仅授权范围）

| 手法 | 工具 | 说明 |
|---|---|---|
| 反向代理钓鱼（抓 MFA 会话） | Evilginx2 | 需授权；练熟"合法登录页 + 代理"链路 |
| 二维码钓鱼 | 自建短链 + 二维码 | 绕过邮件网关，练"扫码后落地页" |
| 鱼叉邮件附件 | 演练专用 payload | **仅靶场/演练环境**，禁止真实客户机器 |
| 物理载体 | USB/Dropbox 标签诱饵 | 需物理授权，练"插U盘后自动回连" |

后利用：凭据进入内网后按 `pentest-ops` 技能走；发现物用 ARTEX finding 记录（附证据截图/时间戳）。

## 5. 检测规避（红队视角）与防守对照（蓝队视角）

红队：SPF/DKIM/DMARC 配好、域名年龄>30天、链接信誉预热、附件加密+口令分渠道给。
蓝队（攻防战防守方用）：见图书馆 `conducting-phishing-incident-response`、`hunting-for-spearphishing-indicators`、`detecting-qr-code-phishing-with-email-security`。

## 6. 报告与清理（必做）

报告要点：演练目标与范围、时间线、各渠道转化率（打开/点击/提交/报告）、暴露的流程缺陷、员工培训建议、**证据链（时间戳+截图+邮件原文）**。
清理：关停演练域名/落地页、删除收集的凭据、销毁本地资料、交付报告。

## 7. 工具就绪检查（已预装，20261019 起）

```bash
# 全部已装好，无需再装：
which gophish evilginx2 dnstwist holehe sherlock linkedin2username \
      theHarvester subfinder httpx nuclei amass katana gau waybackurls 2>/dev/null
```

- `gophish`：钓鱼平台（/opt/gophish，首次运行生成 admin 密码，改配置再对外）
- `evilginx2`：反向代理钓鱼（/usr/local/bin/evilginx2）
- `dnstwist`：相似域名生成（支持 fuzz/字典/正则模式）
- `holehe`：邮箱注册检测；`sherlock`：账号追踪；`linkedin2username`：员工列表生成

若新容器缺失（回滚到旧镜像时），按需重装：
`pip install dnstwist holehe sherlock-project`，gophish 用 release zip，evilginx2 用 go build。
