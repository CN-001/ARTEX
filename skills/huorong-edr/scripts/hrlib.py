#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""hrlib.py —— 火绒本地 SQLite 只读访问公共库。

原理：先把 *.db 与 -wal/-shm 复制到临时目录（sqlite 会重放 WAL），
再以 mode=ro 打开，避免与火绒守护进程抢锁。
仅使用标准库。
"""
import json
import os
import shutil
import sqlite3
import tempfile
import uuid

HR_DIR = r"C:\ProgramData\Huorong\Sysdiag"

FNAME_LABEL = {
    "filemon": "文件实时监控",
    "sysprot": "系统防护",
    "quaran": "隔离区操作",
    "scan": "病毒扫描",
    "c_malinfo": "恶意网址拦截",
    "instmon": "软件安装拦截",
    "sysrepair": "漏洞修复",
    "update": "病毒库更新",
    "microphone": "麦克风保护",
    "netatt": "网络攻击拦截",
    "hips": "自定义防护",
}


def snapshot(name: str) -> str:
    """把 C:\\ProgramData\\Huorong\\Sysdiag\\<name> 及 WAL 复制到临时目录，返回主库路径。"""
    src = os.path.join(HR_DIR, name)
    if not os.path.exists(src):
        raise FileNotFoundError(f"火绒数据文件不存在: {src}")
    dst = os.path.join(tempfile.gettempdir(), f"hr_edr_{uuid.uuid4().hex[:8]}")
    os.makedirs(dst, exist_ok=True)
    for suffix in ("", "-wal", "-shm"):
        p = src + suffix
        if os.path.exists(p):
            shutil.copy2(p, os.path.join(dst, name + suffix))
    return os.path.join(dst, name)


def connect(name: str) -> sqlite3.Connection:
    db = snapshot(name)
    con = sqlite3.connect(f"file:{db}?mode=ro", uri=True)
    con.row_factory = sqlite3.Row
    return con


def parse_detail(raw) -> dict:
    if not raw:
        return {}
    if isinstance(raw, dict):
        return raw
    try:
        return json.loads(raw)
    except Exception:
        return {"_raw": raw}


def flatten_event(row) -> dict:
    """把 HrLogV3_60 行规整成统一结构。"""
    d = parse_detail(row["detail"])
    inner = d.get("detail", {}) if isinstance(d, dict) else {}
    if not isinstance(inner, dict):
        inner = {}
    ev = {
        "id": row["id"],
        "type": row["fname"],
        "type_label": FNAME_LABEL.get(row["fname"], row["fname"]),
        "ts": row["ts"],
        "guid": d.get("guid") if isinstance(d, dict) else None,
        "threat": inner.get("recname", ""),
        "path": inner.get("pathname", inner.get("path", "")),
        "cmdline": inner.get("cmdline", ""),
        "parent_proc": inner.get("p_procname", ""),
        "proc": inner.get("procname", inner.get("proc", "")),
        "token": inner.get("token", ""),
        "treatment": inner.get("treatment", ""),
        "result": inner.get("result", ""),
        "md5": inner.get("md5", ""),
        "sha1": inner.get("sha1", ""),
        "sha256": inner.get("sha256", ""),
        "url": inner.get("url", ""),
        "task": inner.get("taskname", inner.get("task", "")),
    }
    ev["raw"] = d
    return ev


def query_events(con, since_ts=None, until_ts=None, types=None, limit=200):
    sql = "SELECT id,fid,fname,ts,detail,guid FROM HrLogV3_60"
    conds, args = [], []
    if since_ts is not None:
        conds.append("ts>=?"); args.append(since_ts)
    if until_ts is not None:
        conds.append("ts<=?"); args.append(until_ts)
    if types:
        q = ",".join("?" * len(types))
        conds.append(f"fname IN ({q})"); args.extend(types)
    if conds:
        sql += " WHERE " + " AND ".join(conds)
    sql += " ORDER BY id DESC LIMIT ?"
    args.append(limit)
    return [flatten_event(r) for r in con.execute(sql, args)]


def query_quarantine(con, limit=200):
    try:
        rows = con.execute(
            "SELECT id,fn,sha1,vn,fid,ts,rid,info,type FROM FilesV3_60 ORDER BY id DESC LIMIT ?",
            (limit,),
        ).fetchall()
    except sqlite3.OperationalError:
        return []
    out = []
    for r in rows:
        out.append({
            "id": r["id"], "path": r["fn"], "sha1": r["sha1"],
            "threat": r["vn"], "ts": r["ts"], "type": r["type"],
            "info": parse_detail(r["info"]),
        })
    return out


def query_trust(con, limit=500):
    try:
        rows = con.execute(
            "SELECT id,fn,hash,type,memo FROM TrustRegion_60 ORDER BY id LIMIT ?", (limit,)
        ).fetchall()
    except sqlite3.OperationalError:
        return []
    return [dict(r) for r in rows]
