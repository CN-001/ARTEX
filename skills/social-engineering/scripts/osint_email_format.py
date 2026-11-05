#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
osint_email_format.py —— 授权演练：从姓名+域名推断企业邮箱格式并生成候选

用法:
  python osint_email_format.py --domain example.com --names "zhang.san,li.si"
  python osint_email_format.py --domain example.com --names-file names.txt --patterns

输出:
  - 域名 MX 检查（判断是否真实收信域）
  - 每个姓名的候选邮箱（常见企业格式）
  - emails.txt 候选列表

授权提醒: 仅在演练授权范围内对目标组织使用；生成结果用于社工演练/资产核验，
         不得用于批量投递骚扰或撞库。
"""
import argparse
import itertools
import re
import subprocess
import sys

DEFAULT_PATTERNS = [
    "{f}.{l}", "{f}{l}", "{f}", "{l}", "{f}_{l}",
    "{f[0]}{l}", "{f[0]}.{l}", "{l}.{f}", "{f[0]}{l}",
    "{f}.{l}1", "{f}{l}1",
]


def norm(name: str) -> str:
    """把各种姓名字符串归一成 (first, last)。
    'zhang.san' → ('zhang','san')；'zhang san' → 同；'san.zhang' 按常见姓在后处理。
    """
    name = name.strip().lower()
    name = re.sub(r"[_\-]+", " ", name)
    parts = [p for p in re.split(r"[\s.]", name) if p]
    if not parts:
        return None
    if len(parts) == 1:
        return parts[0], ""
    return parts[0], parts[-1]


def candidates(domain: str, first: str, last: str, patterns):
    out = []
    for pat in patterns:
        try:
            local = pat.format(f=first, l=last, F=first.capitalize(), L=last.capitalize())
        except Exception:
            continue
        if local and local not in out:
            out.append(local)
    return ["%s@%s" % (c, domain) for c in out]


def check_mx(domain: str) -> list:
    try:
        out = subprocess.run(["nslookup", "-type=MX", domain],
                             capture_output=True, text=True, timeout=10,
                             errors="ignore").stdout
        mxs = re.findall(r"mail exchanger = ([^\r\n]+)", out)
        return [m.strip().rstrip(".") for m in mxs]
    except Exception:
        return []


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--domain", required=True, help="目标企业域名（如 example.com）")
    ap.add_argument("--names", default="", help="逗号分隔的姓名/拼音（zhang.san,li.si）")
    ap.add_argument("--names-file", default="", help="每行一个姓名")
    ap.add_argument("--patterns", action="store_true", help="输出全部候选格式")
    ap.add_argument("--out", default="emails.txt")
    args = ap.parse_args()

    names = []
    if args.names:
        names.extend([n for n in args.names.split(",") if n.strip()])
    if args.names_file:
        with open(args.names_file, encoding="utf-8", errors="ignore") as f:
            names.extend([ln.strip() for ln in f if ln.strip()])
    if not names:
        print("请用 --names 或 --names-file 提供姓名"); sys.exit(1)

    print("目标域名:", args.domain)
    mxs = check_mx(args.domain)
    if mxs:
        print("MX 记录:", ", ".join(mxs))
        print("→ 该域真实收信，候选邮箱可作为演练目标（需授权）")
    else:
        print("MX: 未查到（可能非邮件域/被限制查询）")

    all_out = []
    for n in names:
        pair = norm(n)
        if not pair:
            continue
        first, last = pair
        if not last:  # 单段名：只生成单段格式
            locals_ = [first, first + "1"]
            cands = ["%s@%s" % (c, args.domain) for c in locals_]
        else:
            pats = DEFAULT_PATTERNS if args.patterns else DEFAULT_PATTERNS[:6]
            cands = candidates(args.domain, first, last, pats)
        all_out.extend(cands)
        print("\n姓名: %s  →  first=%s last=%s" % (n, first, last or "-"))
        for c in cands:
            print("   ", c)

    with open(args.out, "w", encoding="utf-8") as f:
        f.write("\n".join(all_out) + "\n")
    print("\n候选数: %d  →  已写入 %s" % (len(all_out), args.out))


if __name__ == "__main__":
    main()
