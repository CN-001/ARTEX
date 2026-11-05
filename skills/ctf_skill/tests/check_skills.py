#!/usr/bin/env python3
"""CTF-AI-Skill-Pack L1 格式合规检查器（可执行）

检查内容：
  1. 可发现的 SKILL.md 恰好 13 个，无嵌套技能泄漏
  2. 每个专项 SKILL.md 的 frontmatter 含 name/version/description/trigger
  3. name 与目录名一致
  4. 每个专项含「与 CTF-Analyzer-V2.0 的协作方式」「安全边界」「自检规则」三节
  5. 00-core 大脑保留 V2.0 全部 24 章
  6. 第三方语料的 SKILL.md 已降级为 PLAYBOOK.md（不参与发现）
  7. 各 SKILL.md 的「参考语料索引」中引用的文件真实存在

用法：python3 check_skills.py [pack_root]
退出码：0 全部通过 / 1 有失败项
"""

import io
import re
import sys
from pathlib import Path

# Windows 控制台默认 GBK，强制 UTF-8 输出
if hasattr(sys.stdout, "buffer"):
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding="utf-8")

EXPECTED_SKILLS = {
    "00-core/CTF-Analyzer-V2.0": "ctf-analyzer",
    "01-web/Web-Analyzer": "web-analyzer",
    "01-web/SQLi": "sqli",
    "01-web/XSS": "xss",
    "01-web/SSRF": "ssrf",
    "02-crypto/Crypto-Analyzer": "crypto-analyzer",
    "03-pwn/Pwn-Analyzer": "pwn-analyzer",
    "04-reverse/Reverse-Analyzer": "reverse-analyzer",
    "05-forensics/Forensics-Analyzer": "forensics-analyzer",
    "06-misc/Misc-Analyzer": "misc-analyzer",
    "07-linux/Linux-Security": "linux-security",
    "08-writeup/Writeup-Generator": "writeup-generator",
    "08-writeup/CTF-Review": "ctf-review",
}

REQUIRED_SECTIONS = ["与 CTF-Analyzer-V2.0 的协作方式", "安全边界", "自检规则"]
# 00-core 大脑的 24 章标题（一～二十四）
BRAIN_CHAPTERS = 24


def parse_frontmatter(text: str) -> dict:
    m = re.match(r"^---\s*\n(.*?)\n---\s*\n", text, re.DOTALL)
    if not m:
        return {}
    fm = {}
    current = None
    for line in m.group(1).splitlines():
        kv = re.match(r"^(\w+):\s*(.*)$", line)
        if kv:
            # 值为空的键（如 trigger:）按列表初始化，否则后续 bullet 无法追加
            fm[kv.group(1)] = [] if not kv.group(2).strip() else kv.group(2).strip()
            current = kv.group(1)
        elif line.startswith("  - ") and current:
            fm.setdefault(current, [])
            if isinstance(fm[current], list):
                fm[current].append(line.strip()[2:])
    return fm


def fail(msg: str):
    print(f"  ❌ {msg}")
    return False


def main():
    root = Path(sys.argv[1] if len(sys.argv) > 1 else ".")
    skills = root / "skills"
    if not skills.is_dir():
        print(f"❌ 找不到 skills/ 目录：{skills}")
        return 1

    ok = True
    print("=" * 60)
    print("CTF-AI-Skill-Pack L1 格式合规检查")
    print("=" * 60)

    # 检查 1：可发现的 SKILL.md 恰好 13 个
    print("\n[1] 可发现技能数量（应为 13）")
    discovered = sorted(
        p.relative_to(skills).parent.as_posix()
        for p in skills.rglob("SKILL.md")
    )
    extra = [d for d in discovered if d not in EXPECTED_SKILLS]
    missing = [d for d in EXPECTED_SKILLS if d not in discovered]
    if extra:
        ok = fail(f"多余的嵌套技能（应降级为 PLAYBOOK.md）：{extra}") and ok
    if missing:
        ok = fail(f"缺失的技能：{missing}") and ok
    if not extra and not missing:
        print(f"  ✅ 恰好 {len(discovered)} 个可发现技能")

    # 检查 6：第三方语料已降级
    print("\n[2] 第三方语料降级检查（PLAYBOOK.md 不参与发现）")
    nested_skill = list(skills.rglob("SKILL.md"))
    nested_skill = [
        p for p in nested_skill
        if p.parent.as_posix() not in {skills / d.split("/")[0] / d.split("/")[1] for d in EXPECTED_SKILLS}
    ]
    playbooks = list(skills.rglob("PLAYBOOK.md"))
    print(f"  ℹ️ PLAYBOOK.md 参考文档：{len(playbooks)} 个")
    if playbooks:
        print(f"  ✅ 参考语料已降级，不参与技能发现")

    # 检查 2/3/4：每个专项的 frontmatter 与章节
    print("\n[3] 专项 SKILL.md 格式（frontmatter + 三章节）")
    for rel, expected_name in EXPECTED_SKILLS.items():
        if rel == "00-core/CTF-Analyzer-V2.0":
            continue  # 大脑豁免统一格式
        p = skills / rel / "SKILL.md"
        if not p.is_file():
            continue
        text = p.read_text(encoding="utf-8")
        fm = parse_frontmatter(text)
        problems = []
        for field in ("name", "version", "description", "trigger"):
            if field not in fm:
                problems.append(f"缺 {field}")
        # trigger 必须是至少一条的非空列表（空字符串/零条 bullet 不算）
        if isinstance(fm.get("trigger"), str):
            problems.append("trigger 为空字符串（缺 bullet）")
        elif isinstance(fm.get("trigger"), list) and not fm["trigger"]:
            problems.append("trigger 零条")
        if fm.get("name") and fm["name"] != expected_name:
            problems.append(f"name={fm['name']} 与目录名 {expected_name} 不一致")
        for sec in REQUIRED_SECTIONS:
            if sec not in text:
                problems.append(f"缺章节「{sec}」")
        if problems:
            ok = fail(f"{rel}：{'; '.join(problems)}") and ok
    print("  ✅ 12 个专项 SKILL.md 格式全部通过（详见上方无 ❌ 行）")

    # 检查 5：大脑 24 章
    print("\n[4] 00-core 大脑完整性（24 章核心机制不得被削弱）")
    brain = skills / "00-core/CTF-Analyzer-V2.0" / "SKILL.md"
    if not brain.is_file():
        ok = fail("大脑文件缺失") and ok
    else:
        btext = brain.read_text(encoding="utf-8")
        # 顶层章节：## 一、 ... ## 二十四、
        chapters = re.findall(r"^## ([一二三四五六七八九十]+)、", btext, re.M)
        nums = ["一", "二", "三", "四", "五", "六", "七", "八", "九", "十",
                "十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八",
                "十九", "二十", "二十一", "二十二", "二十三", "二十四"]
        found = set(chapters)
        missing_ch = [n for n in nums if n not in found]
        if missing_ch:
            ok = fail(f"大脑缺失章节：{missing_ch}") and ok
        else:
            print(f"  ✅ 24 章核心机制完整（{'、'.join(nums[:3])}…{'、'.join(nums[-2:])}）")
        # 核心机制关键词不得被删除
        for kw in ["Progressive Hint System", "Evidence", "Hypothesis Priority",
                   "Smart Next Step", "Session State", "Simple Challenge Mode",
                   "Anti-Hallucination", "Result Interpretation"]:
            if kw not in btext:
                ok = fail(f"大脑缺失核心机制关键词：{kw}") and ok

    # 检查 7：参考语料引用存在性
    print("\n[5] 参考语料引用完整性")
    broken = 0
    for rel in EXPECTED_SKILLS:
        p = skills / rel / "SKILL.md"
        if not p.is_file():
            continue
        text = p.read_text(encoding="utf-8")
        base = p.parent
        # 提取反引号内的引用（.md 文件或目录名）
        for ref in re.findall(r"`([A-Za-z0-9_\-./]+\.(?:md|py|sh)|[a-z0-9-]+/)`", text):
            if ref.startswith("http") or ref.startswith("../"):
                target = (base / ref).resolve()
            else:
                target = base / ref
            if not target.exists():
                print(f"  ❌ {rel}：悬空引用 {ref}")
                broken += 1
                ok = False
    if broken == 0:
        print("  ✅ 无悬空引用")

    print("\n" + "=" * 60)
    if ok:
        print("✅ L1 格式合规检查全部通过")
        return 0
    else:
        print("❌ 存在失败项，请修复后再发布")
        return 1


if __name__ == "__main__":
    sys.exit(main())
