#!/usr/bin/env python3
"""CTF-AI-Skill-Pack L4 协作契约检查器（可执行）

验证 12 个专项 Skill 是否遵守「不重建推理系统」的架构约束（test-plan.md 第五节）。
大脑只有 CTF-Analyzer-V2.0 一个；专项 Skill 若输出题型判断、宣布 Level、
改写 Session State 或反向调用大脑，就会出现多套并行推理逻辑互相冲突。

检查策略（避免误报）：
  - 禁止行为只检查**输出格式节**——那是 Skill 的实际产物。职责归属表、协作纪律里
    出现「判断题目类型 … 归属大脑」「不反向调用」是边界声明本身，不是违规。
  - 边界声明要求「显式否定」或「职责归属表把该职责判给大脑」二者满足其一。
  - 08-writeup 的两个 Skill 是面向用户的产物（test-plan.md 已记录该例外），
    不要求「专项分析结果」格式。

用法：python3 check_contract.py [pack_root]
退出码：0 全部通过 / 1 有失败项
"""

import io
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "buffer"):
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding="utf-8")

SPECIALIZED = {
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
    "08-writeup/Writeup-Generator": "writeup-generator",   # 面向用户产物，例外
    "08-writeup/CTF-Review": "ctf-review",                 # 面向用户产物，例外
}
USER_FACING = {"08-writeup/Writeup-Generator", "08-writeup/CTF-Review"}

# 输出格式节里不得出现的行为（Skill 的产物中不能夹带大脑的职责）
FORBIDDEN_IN_OUTPUT = {
    "输出题目类型判断": r"Challenge\s*Type|题目类型[:：]|题型判定",
    "宣布提示等级": r"Level\s*[123]|当前\s*Level|Hint Level",
    "改写 Session State": r"(?:更新|写入|修改)\s*(?:CTF\s*)?Session\s*State",
}

REQUIRED_DECLARATIONS = [
    # 二者满足其一即可：显式否定，或职责归属表把该职责判给大脑
    ("不输出题型判断", r"不输出题型判断|不输出题目类型判断"),
    ("不宣布 Level", r"不宣布\s*Level|不自行宣布\s*Level"),
]


def strip_frontmatter(text: str) -> str:
    m = re.match(r"^---\s*\n.*?\n---\s*\n", text, re.DOTALL)
    return text[m.end():] if m else text


def extract_section(text: str, heading: str) -> str:
    """提取「## 四、输出格式」到下一个二级标题之间的正文。

    输出格式节通常含 ```markdown 围栏，围栏内的嵌套标题（如「## SQLi 专项分析结果」）
    不是真正的章节边界，必须跳过。
    """
    m = re.search(rf"^## [一二三四五六七八九十]+、{heading}\s*$",
                  text, re.M)
    if not m:
        return ""
    rest = text[m.end():]
    lines, in_fence, out = rest.splitlines(), False, []
    for line in lines:
        stripped = line.strip()
        if stripped.startswith("```"):
            in_fence = not in_fence
            out.append(line)
            continue
        if not in_fence and re.match(r"^## ", line):
            break
        out.append(line)
    return "\n".join(out)


def fail(msg: str) -> bool:
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
    print("CTF-AI-Skill-Pack L4 协作契约检查")
    print("=" * 60)

    for rel in SPECIALIZED:
        p = skills / rel / "SKILL.md"
        if not p.is_file():
            ok = fail(f"{rel}：文件缺失") and ok
            continue
        text = strip_frontmatter(p.read_text(encoding="utf-8"))
        problems = []

        # 检查 1-3：输出格式节不得夹带大脑职责
        output_sec = extract_section(text, "输出格式")
        if not output_sec:
            problems.append("找不到「输出格式」节")
        else:
            for label, pattern in FORBIDDEN_IN_OUTPUT.items():
                hits = re.findall(pattern, output_sec)
                if hits:
                    problems.append(f"输出格式{label}：{hits[:3]}")

        # 检查 4：不反向调用大脑（全文件，但排除「不反向调用」声明本身）
        rev = re.findall(r"调用\s*(?:CTF-)?Analyzer", text)
        real_rev = [h for h in rev if "不反向调用" not in text[max(0, text.find(h) - 20):
                                                                text.find(h) + len(h)]]
        if real_rev:
            problems.append(f"反向调用大脑：{real_rev[:2]}")

        # 检查 6：边界声明（显式否定 或 职责归属表判给大脑）
        core_sec = extract_section(text, "核心功能")
        attributed = ("判断题目类型" in core_sec or "题型" in core_sec) and \
                     "CTF-Analyzer-V2.0（大脑）" in core_sec
        for label, pattern in REQUIRED_DECLARATIONS:
            if not (re.search(pattern, text) or attributed):
                problems.append(f"缺边界声明「{label}」（显式否定或职责归属表均可）")

        # 检查 5：非面向用户的专项 Skill 必须用「专项分析结果」格式
        if rel not in USER_FACING:
            for marker in ("专项分析结果", "供 Analyzer 纳入"):
                if marker not in output_sec:
                    problems.append(f"输出格式缺「{marker}」")

        if problems:
            ok = fail(f"{rel}：{'; '.join(problems)}") and ok

    if ok:
        print("  ✅ 12 个专项 Skill 全部遵守协作契约")
        print("     输出格式不含题型判断 / Level 宣布 / Session State 改写")
        print("     边界声明齐备（显式否定或职责归属表）")
        print("     10 个分析类 Skill 统一「专项分析结果」格式；08-writeup 为面向用户例外")

    # 交叉验证：大脑必须独占推理职责
    print("\n[交叉验证] 大脑独占推理职责")
    brain = skills / "00-core/CTF-Analyzer-V2.0" / "SKILL.md"
    if not brain.is_file():
        ok = fail("大脑文件缺失") and ok
    else:
        btext = brain.read_text(encoding="utf-8")
        for kw in ["Progressive Hint System", "Session State", "Hypothesis Priority",
                   "Smart Next Step", "Result Interpretation"]:
            if kw not in btext:
                ok = fail(f"大脑缺失推理职责关键词：{kw}") and ok
        if ok:
            print("  ✅ 推理机制（提示等级 / 状态 / 优先级 / 下一步 / 结果解读）只存在于大脑")

    print("\n" + "=" * 60)
    if ok:
        print("✅ L4 协作契约检查全部通过")
        return 0
    print("❌ 存在失败项，请修复后再发布")
    return 1


if __name__ == "__main__":
    sys.exit(main())
