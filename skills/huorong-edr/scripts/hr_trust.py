#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""hr_trust.py —— 火绒信任区管理（C: 信任区联动）。

list   : 只读列出信任区（安全）
add    : 默认 dry-run，打印将执行的 SQL；--apply 才写 wlfile.db
         写入前自动备份 wlfile.db* 到 .bak-<时间戳>
         若火绒守护进程占锁会重试 3s，仍失败则提示关闭火绒自保护后重试

用法:
  python hr_trust.py list
  python hr_trust.py add --path "E:\\ARTEX-main\\data\\labs"
  python hr_trust.py add --path "E:\\ARTEX-main\\data\\tools" --apply
  python hr_trust.py add --artex-dirs --apply   # 批量添加常见 ARTEX 目录（跳过已存在）
"""
import argparse
import json
import os
import shutil
import sqlite3
import sys
import time

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import hrlib  # noqa: E402

WLFILE = os.path.join(hrlib.HR_DIR, "wlfile.db")
BACKUP_DIR = r"E:\ARTEX-main\data\huorong\backup"

ARTEX_DIRS = [
    r"E:\ARTEX-main\data\tools",
    r"E:\ARTEX-main\data\seclists",
    r"E:\ARTEX-main\data\labs",
    r"E:\ARTEX-main\data\sessions",
    r"E:\ARTEX-main\data\tasks",
    r"E:\ARTEX-main\evidence",
    r"E:\ARTEX-main\screenshots",
    r"E:\ARTEX-main\intercept",
    r"E:\ARTEX-main\skills",
]


def backup():
    """备份 wlfile.db* 到 ARTEX 目录（火绒目录受自我保护保护，不能在其中新建文件）。"""
    os.makedirs(BACKUP_DIR, exist_ok=True)
    stamp = time.strftime("%Y%m%d_%H%M%S")
    for suffix in ("", "-wal", "-shm"):
        p = WLFILE + suffix
        if os.path.exists(p):
            dst = os.path.join(BACKUP_DIR, f"wlfile.db{suffix}.bak-{stamp}")
            shutil.copy2(p, dst)
            print(f"  已备份 {p} -> {dst}")


def existing_paths(con):
    try:
        rows = con.execute("SELECT fn FROM TrustRegion_60").fetchall()
        return {r["fn"] for r in rows}
    except sqlite3.OperationalError:
        return set()


def add_paths(paths, apply):
    for p in paths:
        p = p.strip().strip('"')
        if not p:
            continue
        fn = os.path.normpath(p)
        typ = 3 if os.path.isdir(fn) else 0   # type=3 目录信任, type=0 文件信任
        sql = ("INSERT INTO TrustRegion_60 (fn, hash, type, memo) "
               f"VALUES ('{fn}', '', {typ}, '');")
        if not apply:
            print(f"  [dry-run] 将添加(type={typ}): {fn}")
            print(f"            SQL: {sql}")
        else:
            print(f"  写入(type={typ}): {fn}")
    if not apply:
        print("\n未写入。确认无误后加 --apply 执行。")
        return
    backup()
    try:
        con = sqlite3.connect(WLFILE, timeout=3)
        con.row_factory = sqlite3.Row
    except Exception as e:
        print(f"打开 wlfile.db 失败: {e}", file=sys.stderr)
        sys.exit(1)
    exist = existing_paths(con)
    added, skipped = 0, 0
    try:
        for p in paths:
            fn = os.path.normpath(p.strip().strip('"'))
            if not fn:
                continue
            if fn in exist:
                skipped += 1
                continue
            typ = 3 if os.path.isdir(fn) else 0
            con.execute("INSERT INTO TrustRegion_60 (fn, hash, type, memo) VALUES (?, '', ?, '')", (fn, typ))
            added += 1
        con.commit()
        print(f"完成: 新增 {added} 条，跳过已存在 {skipped} 条")
        print("提示: 若火绒界面未立即显示，重启火绒或稍后刷新。")
    except (sqlite3.OperationalError, PermissionError) as e:
        print(f"写入失败（可能被火绒自我保护拦截）: {e}", file=sys.stderr)
        print("建议: 火绒设置 → 系统设置 → 临时关闭自我保护后再跑，或用火绒 UI 手动添加。", file=sys.stderr)
        sys.exit(1)
    finally:
        con.close()


def main():
    ap = argparse.ArgumentParser(description="火绒信任区管理")
    sub = ap.add_subparsers(dest="cmd", required=True)

    sub.add_parser("list", help="列出信任区")

    p = sub.add_parser("add", help="添加信任路径（默认 dry-run）")
    p.add_argument("--path", action="append", default=[], help="路径，可多次")
    p.add_argument("--artex-dirs", action="store_true", help="批量添加常见 ARTEX 目录")
    p.add_argument("--apply", action="store_true", help="实际写入（先自动备份）")

    args = ap.parse_args()

    if args.cmd == "list":
        con = hrlib.connect("wlfile.db")
        rows = hrlib.query_trust(con, limit=1000)
        con.close()
        print(f"信任区共 {len(rows)} 条：")
        for r in rows:
            print(f"  [{r['id']}] type={r['type']} {r['fn']}")
        return

    if args.cmd == "add":
        paths = list(args.path)
        if args.artex_dirs:
            paths.extend(ARTEX_DIRS)
        if not paths:
            print("请提供 --path 或 --artex-dirs")
            sys.exit(1)
        add_paths(paths, args.apply)


if __name__ == "__main__":
    main()
