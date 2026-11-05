#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""hr_api.py —— 火绒数据 HTTP API（供 ARTEX agent 自定义工具调用）。

把火绒本地 SQLite 数据以 HTTP 暴露给 ARTEX（无论 ARTEX 跑在 Docker 还是本机）：
  Docker 部署: agent 工具用 http://host.docker.internal:18888
  本机部署:    agent 工具用 http://127.0.0.1:18888

鉴权: 请求需带 X-HR-Token（首次启动自动生成 data\\huorong\\api_token.txt）。

端点:
  GET /healthz
  GET /api/status                防护状态/病毒库版本/事件统计
  GET /api/events?since=60m&until=&types=filemon,quaran&limit=50&format=markdown|json
  GET /api/quarantine            隔离区
  GET /api/trust                 信任区
"""
import json
import os
import sys
import time
import urllib.parse
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import hrlib  # noqa: E402

HOST = "0.0.0.0"
PORT = int(os.environ.get("HR_API_PORT", "18888"))
TOKEN_FILE = r"E:\ARTEX-main\data\huorong\api_token.txt"


def load_token() -> str:
    os.makedirs(os.path.dirname(TOKEN_FILE), exist_ok=True)
    if os.path.exists(TOKEN_FILE):
        tok = open(TOKEN_FILE, "r", encoding="utf-8").read().strip()
        if tok:
            return tok
    tok = __import__("secrets").token_hex(16)
    open(TOKEN_FILE, "w", encoding="utf-8").write(tok)
    return tok


TOKEN = load_token()


def parse_dur(s: str) -> int:
    if not s:
        return 0
    s = s.strip().lower()
    unit = {"s": 1, "m": 60, "h": 3600, "d": 86400}.get(s[-1])
    if unit is None:
        raise ValueError(f"无法解析时长: {s}")
    return int(float(s[:-1]) * unit)


def fmt_md_events(events):
    lines = ["| 时间 | 类型 | 威胁名 | 路径 | 进程 | 处置 |",
             "|---|---|---|---|---|---|"]
    for e in events:
        lines.append(
            f"| {time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(e['ts']))} | "
            f"{e['type_label']} | {e['threat']} | `{e['path'] or e['url']}` | "
            f"`{e['proc'] or e['parent_proc']}` | {e['token']} |")
    return "\n".join(lines)


class Handler(BaseHTTPRequestHandler):
    def log_message(self, *a):
        pass

    def _send(self, code, obj):
        body = json.dumps(obj, ensure_ascii=False).encode("utf-8")
        self.send_response(code)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def _auth(self) -> bool:
        return self.headers.get("X-HR-Token", "") == TOKEN

    def do_GET(self):
        parsed = urllib.parse.urlparse(self.path)
        path = parsed.path
        qs = urllib.parse.parse_qs(parsed.query)
        if path == "/healthz":
            self._send(200, {"ok": True})
            return
        if not self._auth():
            self._send(401, {"error": "unauthorized: 需要 X-HR-Token"})
            return
        try:
            if path == "/api/status":
                self._status()
            elif path == "/api/events":
                self._events(qs)
            elif path == "/api/quarantine":
                self._quarantine()
            elif path == "/api/trust":
                self._trust()
            else:
                self._send(404, {"error": "未知端点"})
        except Exception as e:
            self._send(500, {"error": str(e)})

    def _status(self):
        con = hrlib.connect("log.db")
        ver = {}
        row = con.execute(
            "SELECT detail FROM HrLogV3_60 WHERE fname IN ('scan','update') ORDER BY id DESC LIMIT 1"
        ).fetchone()
        if row:
            d = hrlib.parse_detail(row["detail"])
            if isinstance(d, dict):
                ver = d.get("version", {}) or {}
        counts = {r["fname"]: r["n"] for r in con.execute(
            "SELECT fname, COUNT(*) AS n FROM HrLogV3_60 GROUP BY fname")}
        con.close()
        self._send(200, {
            "product": ver.get("product", "未知"),
            "virdb_time": time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(ver["dbtime"])) if ver.get("dbtime") else "未知",
            "event_total": sum(counts.values()),
            "event_by_type": {hrlib.FNAME_LABEL.get(k, k): v for k, v in counts.items()},
        })

    def _events(self, qs):
        since = ((qs.get("since") or [None])[0] or "").strip() or None
        until = ((qs.get("until") or [None])[0] or "").strip() or None
        types = ((qs.get("types") or [None])[0] or "").strip() or None
        fmt = ((qs.get("format") or ["markdown"])[0] or "markdown").strip() or "markdown"
        try:
            limit = int(((qs.get("limit") or ["100"])[0]))
        except (TypeError, ValueError):
            limit = 100
        now = int(time.time())
        since_ts = now - parse_dur(since) if since else None
        until_ts = now - parse_dur(until) if until else None
        tlist = [t.strip() for t in types.split(",") if t.strip()] if types else None
        con = hrlib.connect("log.db")
        events = hrlib.query_events(con, since_ts=since_ts, until_ts=until_ts, types=tlist, limit=limit)
        con.close()
        if fmt == "json":
            self._send(200, {"count": len(events), "events": events})
        else:
            self._send(200, {"count": len(events), "text": fmt_md_events(events)})

    def _quarantine(self):
        con = hrlib.connect("QuarantineEx.db")
        rows = hrlib.query_quarantine(con, limit=200)
        con.close()
        self._send(200, {"count": len(rows), "items": rows})

    def _trust(self):
        con = hrlib.connect("wlfile.db")
        rows = hrlib.query_trust(con, limit=1000)
        con.close()
        self._send(200, {"count": len(rows), "items": rows})


def main():
    srv = ThreadingHTTPServer((HOST, PORT), Handler)
    print(f"hr_api 已启动: http://{HOST}:{PORT}  token={TOKEN}")
    print(f"  Docker agent 访问: http://host.docker.internal:{PORT}  本机: http://127.0.0.1:{PORT}")
    try:
        srv.serve_forever()
    except KeyboardInterrupt:
        pass


if __name__ == "__main__":
    main()
