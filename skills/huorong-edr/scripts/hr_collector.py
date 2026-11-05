#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""hr_collector.py —— 火绒事件实时采集器（B: 实时联动）。

轮询 log.db 的新事件，增量写入 JSONL；可选 POST 到 webhook。
中断后从 state 文件续传，不丢事件。

用法:
  python hr_collector.py --interval 10 --out E:\\ARTEX-main\\data\\huorong\\events.jsonl
  python hr_collector.py --interval 10 --webhook http://127.0.0.1:8000/ingest/hr
  python hr_collector.py --once    # 只拉取一次增量
"""
import argparse
import json
import os
import sys
import time
import urllib.request

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import hrlib  # noqa: E402

DEFAULT_OUT = r"E:\ARTEX-main\data\huorong\events.jsonl"
DEFAULT_STATE = r"E:\ARTEX-main\data\huorong\collector_state.json"


def load_state(path):
    if os.path.exists(path):
        try:
            with open(path, "r", encoding="utf-8") as f:
                return json.load(f)
        except Exception:
            pass
    return {"last_id": 0, "last_ts": 0}


def save_state(path, st):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        json.dump(st, f)


def append_events(out_path, events):
    if not events:
        return
    os.makedirs(os.path.dirname(out_path), exist_ok=True)
    with open(out_path, "a", encoding="utf-8") as f:
        for e in events:
            f.write(json.dumps(e, ensure_ascii=False) + "\n")


def post_webhook(url, events):
    data = json.dumps(events, ensure_ascii=False).encode("utf-8")
    req = urllib.request.Request(url, data=data, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            return resp.status
    except Exception as e:
        print(f"  [webhook] 失败: {e}", file=sys.stderr)
        return None


def poll_once(state, out_path, webhook):
    con = hrlib.connect("log.db")
    events = hrlib.query_events(con, since_ts=None, until_ts=None, types=None, limit=1000)
    con.close()
    # 新事件 = id > state.last_id，按 id 升序
    fresh = [e for e in events if e["id"] > state["last_id"]]
    fresh.sort(key=lambda e: e["id"])
    if fresh:
        state["last_id"] = fresh[-1]["id"]
        state["last_ts"] = fresh[-1]["ts"]
        append_events(out_path, fresh)
        print(f"[{time.strftime('%H:%M:%S')}] +{len(fresh)} 条事件 (last_id={state['last_id']})")
        for e in fresh:
            print(f"    {hrlib.FNAME_LABEL.get(e['type'], e['type'])} | {e['threat']} | {e['path']}")
        if webhook:
            post_webhook(webhook, fresh)
    else:
        print(f"[{time.strftime('%H:%M:%S')}] 无新事件")
    return fresh


def main():
    ap = argparse.ArgumentParser(description="火绒事件实时采集器")
    ap.add_argument("--interval", type=int, default=10, help="轮询间隔秒（默认10）")
    ap.add_argument("--out", default=DEFAULT_OUT, help="JSONL 输出路径")
    ap.add_argument("--state", default=DEFAULT_STATE, help="续传状态文件")
    ap.add_argument("--webhook", default=None, help="可选 webhook URL，逐批 POST")
    ap.add_argument("--once", action="store_true", help="只拉一次增量")
    args = ap.parse_args()

    state = load_state(args.state)
    print(f"huorong collector: out={args.out} state.last_id={state['last_id']}")
    poll_once(state, args.out, args.webhook)
    save_state(args.state, state)

    while not args.once:
        try:
            time.sleep(args.interval)
            poll_once(state, args.out, args.webhook)
            save_state(args.state, state)
        except KeyboardInterrupt:
            save_state(args.state, state)
            print("已停止，状态已保存")
            break
        except Exception as e:
            print(f"轮询出错（跳过本轮）: {e}", file=sys.stderr)
            time.sleep(args.interval * 2)


if __name__ == "__main__":
    main()
