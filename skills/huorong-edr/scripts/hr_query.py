#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""hr_query.py —— 查询本机火绒的事件 / 状态 / 隔离区 / 信任区。

用法:
  python hr_query.py status
  python hr_query.py events --since 60m [--types filemon,quaran] [--format markdown|json|table]
  python hr_query.py events --since-ts 1792596043 --until-ts 1792599000 --format markdown
  python hr_query.py quarantine [--format markdown]
  python hr_query.py trust [--format markdown]
"""
import argparse
import json
import os
import sys
import time

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import hrlib  # noqa: E402


def parse_window(since=None, until=None, since_ts=None, until_ts=None):
    now = int(time.time())
    since_ts = since_ts or (now - parse_dur(since) if since else None)
    until_ts = until_ts or (now - parse_dur(until) if until else None)
    return since_ts, until_ts


def parse_dur(s: str):
    """把 60m / 24h / 3d 转成秒。"""
    if not s:
        return 0
    s = s.strip().lower()
    unit = {"s": 1, "m": 60, "h": 3600, "d": 86400}.get(s[-1])
    if unit is None:
        raise ValueError(f"无法解析时长: {s}（用 60m / 24h / 3d）")
    return int(float(s[:-1]) * unit)


def ts_str(ts):
    return time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(ts))


def fmt_events(events, fmt):
    if fmt == "json":
        return json.dumps(events, ensure_ascii=False, indent=2)
    if fmt == "markdown":
        lines = ["| 时间 | 类型 | 威胁名 | 路径 | 进程 | 处置 |",
                 "|---|---|---|---|---|---|"]
        for e in events:
            lines.append(
                f"| {ts_str(e['ts'])} | {e['type_label']} | {e['threat']} | "
                f"`{e['path'] or e['url']}` | `{e['proc'] or e['parent_proc']}` | {e['token']} |"
            )
        return "\n".join(lines)
    # table
    rows = [["时间", "类型", "威胁名", "路径", "进程", "处置"]]
    for e in events:
        rows.append([ts_str(e["ts"]), e["type_label"], e["threat"],
                     (e["path"] or e["url"] or "")[:60], (e["proc"] or e["parent_proc"])[:40], str(e["token"])])
    widths = [max(len(r[i]) for r in rows) for i in range(len(rows[0]))]
    out = []
    for r in rows:
        out.append("  ".join(r[i].ljust(widths[i]) for i in range(len(r))))
    return "\n".join(out)


def cmd_status():
    out = {}
    con = hrlib.connect("log.db")
    # 病毒库版本/产品版本从最近一条 scan/update 事件的 detail 取
    row = con.execute(
        "SELECT detail FROM HrLogV3_60 WHERE fname IN ('scan','update') ORDER BY id DESC LIMIT 1"
    ).fetchone()
    ver = {}
    if row:
        d = hrlib.parse_detail(row["detail"])
        ver = d.get("version", {}) if isinstance(d, dict) else {}
        if isinstance(ver, dict):
            ver = {k: v for k, v in ver.items() if k in ("product", "dbtime")}
    counts = {r["fname"]: r["n"] for r in con.execute(
        "SELECT fname, COUNT(*) AS n FROM HrLogV3_60 GROUP BY fname")}
    total = sum(counts.values())
    out["product_version"] = ver.get("product", "未知")
    out["virdb_time"] = ts_str(ver["dbtime"]) if ver.get("dbtime") else "未知"
    out["event_total"] = total
    out["event_by_type"] = {hrlib.FNAME_LABEL.get(k, k): v for k, v in counts.items()}
    con.close()
    return json.dumps(out, ensure_ascii=False, indent=2)


def cmd_events(args):
    since_ts, until_ts = parse_window(args.since, args.until, args.since_ts, args.until_ts)
    types = [t.strip() for t in args.types.split(",")] if args.types else None
    con = hrlib.connect("log.db")
    events = hrlib.query_events(con, since_ts=since_ts, until_ts=until_ts,
                                types=types, limit=args.limit)
    con.close()
    print(fmt_events(events, args.format))


def cmd_quarantine(args):
    con = hrlib.connect("QuarantineEx.db")
    rows = hrlib.query_quarantine(con, limit=args.limit)
    con.close()
    if args.format == "json":
        print(json.dumps(rows, ensure_ascii=False, indent=2))
        return
    print("| id | 时间 | 威胁名 | 路径 | sha1 |")
    print("|---|---|---|---|---|")
    for r in rows:
        print(f"| {r['id']} | {ts_str(r['ts'])} | {r['threat']} | `{r['path']}` | {r['sha1'][:12]}… |")


def cmd_trust(args):
    con = hrlib.connect("wlfile.db")
    rows = hrlib.query_trust(con, limit=args.limit)
    con.close()
    if args.format == "json":
        print(json.dumps(rows, ensure_ascii=False, indent=2))
        return
    print(f"共 {len(rows)} 条信任项：")
    for r in rows:
        print(f"  [{r['id']}] type={r['type']} {r['fn']}")


def main():
    ap = argparse.ArgumentParser(description="查询本机火绒数据")
    sub = ap.add_subparsers(dest="cmd", required=True)

    p = sub.add_parser("status", help="防护状态/病毒库版本/事件统计")

    p = sub.add_parser("events", help="安全事件（log.db HrLogV3_60）")
    p.add_argument("--since", help="相对时长，如 60m/24h/3d")
    p.add_argument("--until", help="相对时长，如 10m")
    p.add_argument("--since-ts", type=int, help="起始 unix 秒")
    p.add_argument("--until-ts", type=int, help="结束 unix 秒")
    p.add_argument("--types", help="逗号分隔: filemon,sysprot,quaran,scan,c_malinfo,instmon")
    p.add_argument("--format", choices=["markdown", "json", "table"], default="markdown")
    p.add_argument("--limit", type=int, default=200)

    p = sub.add_parser("quarantine", help="隔离区（QuarantineEx.db）")
    p.add_argument("--format", choices=["markdown", "json"], default="markdown")
    p.add_argument("--limit", type=int, default=200)

    p = sub.add_parser("trust", help="信任区（wlfile.db TrustRegion_60）")
    p.add_argument("--format", choices=["markdown", "json"], default="markdown")
    p.add_argument("--limit", type=int, default=500)

    args = ap.parse_args()
    if args.cmd == "status":
        print(cmd_status())
    elif args.cmd == "events":
        cmd_events(args)
    elif args.cmd == "quarantine":
        cmd_quarantine(args)
    elif args.cmd == "trust":
        cmd_trust(args)


if __name__ == "__main__":
    main()
