---
name: huorong-edr
description: 读取本机火绒(Huorong)安全软件的事件日志、隔离区与信任区，用于红蓝对抗演练的检测对照。调用该skill查询火绒检测/拦截记录。
---

# Huorong EDR（火绒联动 · 只读为主）

读取本机火绒个人版（6.x）的本地 SQLite 数据，把火绒的**检测/拦截事件**接入 ARTEX 任务流程：

- 演练前：查当前防护状态（病毒库版本、事件基线）
- 演练中：按时间窗查火绒的 `filemon/sysprot/quaran/scan` 事件，对照 ARTEX 攻击动作
- 演练后：生成「哪些攻击被火绒检出、哪些绕过」的检测覆盖率报告

**火绒个人版无公开 API/CLI，本 skill 只做只读查询；信任区写入见 `hr_trust.py`（默认 dry-run）。**

## 融合进 ARTEX agent（HTTP API + 自定义工具）

火绒数据在 Windows 主机，ARTEX agent 可能跑在 Docker（Linux）里，直接读不了。
因此提供 HTTP 桥接 + 自定义工具注册，agent 可以**自动调用**：

```
1. 主机后台跑 API（首次自动生成 token 到 data\huorong\api_token.txt）:
   python scripts/hr_api.py            # 监听 0.0.0.0:18888
   GET /api/{status|events|quarantine|trust}?since=60m&types=filemon&limit=50
   鉴权: 请求头 X-HR-Token: <token>

2. ARTEX 启动后，注册自定义工具（一次）:
   powershell -File scripts/register_tools.ps1 -ApiKey <ARTEX的API Key>
   → 向 POST /api/tools/custom 注册 huorong_query（kind=http）
   → agent 工具列表出现 huorong_query，任务里说"查火绒拦截事件"即可自动调用

3. Docker 部署时 agent 通过 http://host.docker.internal:18888 访问；
   本机部署时把 register_tools.ps1 的 -HrBase 改成 http://127.0.0.1:18888
```

## 脚本位置与调用

```bash
# 事件查询（核心）
python scripts/hr_query.py events --since 60m --format markdown
python scripts/hr_query.py events --types filemon,quaran --since 24h --format json

# 状态/隔离区/信任区
python scripts/hr_query.py status
python scripts/hr_query.py quarantine --format markdown
python scripts/hr_query.py trust

# 实时采集（长任务用）
python scripts/hr_collector.py --interval 10 --out data/huorong/events.jsonl

# 信任区管理（默认 dry-run，输出将执行的 SQL；--apply 才写）
python scripts/hr_trust.py list
python scripts/hr_trust.py add --path "E:\ARTEX-main\data\labs" --apply
```

所有查询对 `C:\ProgramData\Huorong\Sysdiag\*.db` 先做「复制 + WAL 合并」再只读打开，避免与火绒守护进程抢锁。

## 事件类型对照（log.db / HrLogV3_60）

| fname | 含义 | detail 关键字段 |
|---|---|---|
| filemon | 文件实时监控（病毒拦截） | recname(威胁名), path, cmdline, p_procname, token(处置) |
| sysprot | 系统防护（漏洞利用/行为拦截） | 按规则不同字段不同 |
| quaran | 隔离区操作 | type(QUARAN_add/restore), success/fail 文件列表 |
| scan | 病毒扫描 | threats, objects, files, duration, threat_list |
| c_malinfo | 恶意网址拦截 | url 等 |
| instmon | 软件安装拦截 | 安装包信息 |
| sysrepair | 漏洞修复 | scannedcount, repairedcount |
| update | 病毒库更新 | fetched 列表 |
| microphone | 麦克风保护 | 进程信息 |

## 边界与禁止

| 允许 | 禁止 |
|---|---|
| 只读查询本机火绒日志/隔离区/信任区 | 关闭/篡改火绒防护、清空日志、删除隔离区文件 |
| 将检测事件作为演练证据写入报告 | 直接修改火绒配置库（hr_trust.py --apply 除外，且需先备份） |
| 时间窗过滤、结构化导出 | 把本机真实拦截数据外传（webhook 仅限内网/自用地址） |

## 演练检测对照流程（推荐用法）

1. `hr_query.py status` + `events --since 24h` 建立基线；
2. ARTEX 执行攻击任务，记录起止时间 T1–T2；
3. `hr_query.py events --since-ts T1 --until-ts T2 --format markdown` 拉出该窗口的火绒事件；
4. 将火绒 `recname/path` 与 ARTEX 动作对照，得到：检出项（火绒有记录）与绕过项（ARTEX 有动作、火绒无记录）。
