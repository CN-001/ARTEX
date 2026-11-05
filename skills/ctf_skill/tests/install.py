#!/usr/bin/env python3
"""CTF-AI-Skill-Pack 安装器（展平分类层）

**为什么需要展平：** 本 Pack 的仓库布局是分类分组（skills/01-web/SQLi/SKILL.md），
但 Claude Code 等加载器只识别 <skills-root>/<skill-name>/SKILL.md（直接子目录）。
直接 `cp -r skills/* ~/.claude/skills/` 会装成 00-core/CTF-Analyzer-V2.0/SKILL.md
这种三层路径，加载器发现不了任何一个技能。

本脚本把 13 个技能目录展平复制到目标目录，并验证安装结果确实可被发现。

用法：
  python3 tests/install.py ~/.claude/skills        # 安装
  python3 tests/install.py /tmp/verify --dry-run   # 只模拟不写盘
退出码：0 安装成功且验证通过 / 1 失败
"""

import io
import shutil
import sys
from pathlib import Path

if hasattr(sys.stdout, "buffer"):
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding="utf-8")

EXPECTED = {
    "CTF-Analyzer-V2.0", "Web-Analyzer", "SQLi", "XSS", "SSRF",
    "Crypto-Analyzer", "Pwn-Analyzer", "Reverse-Analyzer",
    "Forensics-Analyzer", "Misc-Analyzer", "Linux-Security",
    "Writeup-Generator", "CTF-Review",
}


def discover(src: Path) -> list[Path]:
    """skills/<分类>/<技能>/ 三层布局下，所有含 SKILL.md 的技能目录。"""
    found = []
    for p in sorted(src.rglob("SKILL.md")):
        rel = p.relative_to(src)
        if len(rel.parts) != 3:  # <分类>/<技能>/SKILL.md
            print(f"  ⚠️ 跳过非标准深度：{rel}")
            continue
        found.append(p.parent)
    return found


def main():
    args = sys.argv[1:]
    dry = "--dry-run" in args
    args = [a for a in args if a != "--dry-run"]
    if not args:
        print("用法：install.py <目标 skills 目录> [--dry-run]")
        return 1

    src = Path(__file__).resolve().parent.parent / "skills"
    dst = Path(args[0]).expanduser()
    if not src.is_dir():
        print(f"❌ 找不到源目录：{src}")
        return 1

    skills = discover(src)
    print("=" * 60)
    print(f"CTF-AI-Skill-Pack 安装器{'（dry-run，不写盘）' if dry else ''}")
    print(f"源：{src}")
    print(f"目标：{dst}")
    print("=" * 60)

    names = {p.name for p in skills}
    if names != EXPECTED:
        print(f"❌ 技能集合与预期不符")
        print(f"   缺失：{sorted(EXPECTED - names)}")
        print(f"   多余：{sorted(names - EXPECTED)}")
        return 1

    print(f"\n将安装 {len(skills)} 个技能（展平分类层）：")
    for p in skills:
        print(f"  {p.relative_to(src).as_posix():<32} → {p.name}/")

    if dry:
        print("\n[dry-run] 未写入磁盘。去掉 --dry-run 执行实际安装。")
        return 0

    dst.mkdir(parents=True, exist_ok=True)
    for p in skills:
        target = dst / p.name
        if target.exists():
            shutil.rmtree(target)
        shutil.copytree(p, target)
    print(f"\n✅ 已复制到 {dst}")

    # 验证：目标目录的直接子目录里恰好 13 个 SKILL.md，且无更深层嵌套
    installed = sorted(q.parent.name for q in dst.glob("*/SKILL.md"))
    nested = list(dst.rglob("SKILL.md"))
    nested = [q for q in nested if len(q.relative_to(dst).parts) != 2]
    print("\n--- 安装验证 ---")
    if len(installed) != 13:
        print(f"  ❌ 目标目录只发现 {len(installed)} 个直接子技能（应为 13）")
        return 1
    if nested:
        print(f"  ❌ 目标目录存在深层嵌套 SKILL.md：{nested[:3]}")
        return 1
    missing = EXPECTED - set(installed)
    if missing:
        print(f"  ❌ 缺失技能：{sorted(missing)}")
        return 1
    print(f"  ✅ 加载器将发现恰好 {len(installed)} 个技能：{', '.join(installed)}")
    print("  ✅ 无深层嵌套 SKILL.md（第三方语料的 PLAYBOOK.md 不参与发现）")
    print("\n安装完成。重启 Claude Code 会话后技能生效。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
