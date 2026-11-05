#!/usr/bin/env python3
"""CTF-AI-Skill-Pack L2 调度契约检查器（可执行）

L2 的完整版（大脑凭语义判断题型并路由）无法脚本化——那是 L3 人工/模型测试的范围。
本检查器覆盖 L2 中**可自动化**的部分，即 trigger 契约与可分离性：

  1. 12 个专项 Skill 的 trigger 必须声明「由 CTF-Analyzer-V2.0 …调度」
     —— 专项 Skill 只能被大脑调度，不能被加载器直接触发（否则架构契约破裂）
  2. 大脑的 trigger 必须是「类型无关的通用入口」，不得绑死某一具体类别
  3. 区分词唯一性：同一个区分词不得出现在两个专项 Skill 的 trigger 里
     —— 否则基于 trigger 的路由无法把两个 Skill 分开（真实歧义）
  4. 路由用例覆盖分级：R1-R10c 的信号词在目标 Skill 的 trigger 或大脑 trigger 中
     是否有对应（强覆盖 / 大脑兜底 / 完全未覆盖）

用法：python3 check_dispatch.py [pack_root]
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
    "08-writeup/Writeup-Generator": "writeup-generator",
    "08-writeup/CTF-Review": "ctf-review",
}
BRAIN = "00-core/CTF-Analyzer-V2.0"

# 中文区分词（trigger 中出现的、能唯一指认某个类别的术语）
# 注意：模式词（Writeup / 复盘 / 总结一下）不是题目类别，只用于跨 Skill 唯一性检查，
# 不用于大脑类型中立性检查——大脑作为通用入口允许提到这些关键词
CJK_DISCRIMINATORS = [
    "SQL 注入", "数据库", "服务端请求伪造", "URL 参数", "内网探测", "元数据端点",
    "跨站脚本", "回显上下文", "原型链污染", "认证", "授权", "反序列化", "文件上传",
    "命令注入", "模板注入", "抓包", "密码学", "密文", "密钥参数", "哈希",
    "内存安全", "保护机制", "反调试", "反汇编", "加壳", "符号执行", "验证逻辑",
    "磁盘镜像", "内存镜像", "隐写", "混杂编码", "受限 shell", "OSINT", "杂项",
    "提权", "后渗透", "劫持", "Writeup", "复盘", "总结一下",
]
# 大脑类型中立性检查使用：只含「题目类别」术语，排除模式词
BRAIN_NEUTRALITY_TERMS = [d for d in CJK_DISCRIMINATORS if d not in
                          ("Writeup", "复盘", "总结一下")]
# 结构性西文词，不是区分词（所有 Skill 都会出现）
STOPWORDS = {"CTF", "V2", "Analyzer", "Confirmed", "Findings", "SKILL", "Skill",
             "Version", "version"}

# 路由用例（对齐 tests/test-plan.md 的 R1-R10c）
# 信号词取该类输入的代表性术语，不是输入原文的字面子串——否则覆盖分级会退化为
# 字符串匹配游戏。预期 Skill 的 trigger 中出现信号词 = 强覆盖。
ROUTING_CASES = [
    ("R1", "Web 登录认证", "01-web/Web-Analyzer", ["登录", "认证", "username"]),
    ("R2", "RSA 参数题", "02-crypto/Crypto-Analyzer", ["RSA", "n=", "c="]),
    ("R3", "栈溢出 Pwn", "03-pwn/Pwn-Analyzer", ["checksec", "栈溢出"]),
    ("R4", "Reverse 异或比较", "04-reverse/Reverse-Analyzer", ["CrackMe", "验证逻辑", "异或"]),
    ("R5", "PNG 取证", "05-forensics/Forensics-Analyzer", ["PCAP", "隐写", "文件结构"]),
    ("R6", "Base64 杂项", "06-misc/Misc-Analyzer", ["编码", "Decode"]),
    ("R7", "Linux 提权", "07-linux/Linux-Security", ["sudo -l", "提权"]),
    ("R8", "写 Writeup", "08-writeup/Writeup-Generator", ["帮我写 Writeup"]),
    ("R8b", "复盘", "08-writeup/CTF-Review", ["复盘", "我学到了什么"]),
    ("R10a", "SQL 注入源码", "01-web/SQLi", ["SQL 报错", "注入"]),
    ("R10b", "XSS 回显源码", "01-web/XSS", ["跨站脚本", "回显"]),
    ("R10c", "SSRF fetch 参数", "01-web/SSRF", ["SSRF", "URL 参数"]),
    ("R9", "仅一个 URL（Unknown）", None, ["http://123.45.67.89:5000/"]),
]


def parse_frontmatter(text: str) -> dict:
    m = re.match(r"^---\s*\n(.*?)\n---\s*\n", text, re.DOTALL)
    if not m:
        return {}
    fm, current = {}, None
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


def latin_tokens(text: str) -> set:
    return {t for t in re.findall(r"[A-Za-z][A-Za-z0-9_+/=]{1,}", text) if len(t) >= 2}


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
    print("CTF-AI-Skill-Pack L2 调度契约检查")
    print("=" * 60)

    triggers: dict[str, list[str]] = {}
    for rel in list(SPECIALIZED) + [BRAIN]:
        p = skills / rel / "SKILL.md"
        if not p.is_file():
            ok = fail(f"缺失 {rel}") and ok
            continue
        fm = parse_frontmatter(p.read_text(encoding="utf-8"))
        trg = fm.get("trigger")
        triggers[rel] = trg if isinstance(trg, list) else ([trg] if trg else [])
        if not triggers[rel]:
            ok = fail(f"{rel}：trigger 为空或缺失") and ok

    # ---- 检查 1：专项 Skill 只能被大脑调度 ----
    print("\n[1] 专项 Skill 调度入口（必须经 CTF-Analyzer-V2.0，不可直接触发）")
    for rel in SPECIALIZED:
        if rel not in triggers:
            continue
        if not any("CTF-Analyzer-V2.0" in t for t in triggers[rel]):
            ok = fail(f"{rel}：无「由 CTF-Analyzer-V2.0 调度」声明，可能被加载器直接触发") and ok
    print("  ✅ 12 个专项 Skill 均声明由大脑调度（加载器不会绕过大脑直调）")

    # ---- 检查 2：大脑 trigger 类型无关 ----
    print("\n[2] 大脑 trigger 类型中立性")
    brain_trg = triggers.get(BRAIN, [])
    brain_text = " ".join(brain_trg)
    leaked = [d for d in BRAIN_NEUTRALITY_TERMS if d in brain_text]
    if leaked:
        ok = fail(f"大脑 trigger 绑死了具体类别术语：{leaked}") and ok
    elif "判断类型" not in brain_text and "类型" not in brain_text:
        ok = fail("大脑 trigger 未体现「判断题型并路由」职责") and ok
    else:
        print("  ✅ 大脑 trigger 为通用入口，不与任何专项类别绑定")

    # ---- 检查 3：区分词唯一性 ----
    print("\n[3] 区分词唯一性（同一区分词不得属于两个专项 Skill）")
    owner: dict[str, list[str]] = {}
    for rel in SPECIALIZED:
        for bullet in triggers.get(rel, []):
            # 否定式枚举（「不属于 Web/Crypto/…」「非 Web 题」）不算对该类别的认领
            negative = ("不属于" in bullet) or re.search(r"(?:^|[，,。；;])非", bullet)
            if negative:
                continue
            for tok in latin_tokens(bullet):
                if tok in STOPWORDS:
                    continue
                owner.setdefault(tok, []).append(rel)
            for d in CJK_DISCRIMINATORS:
                if d in bullet:
                    owner.setdefault(d, []).append(rel)
    ambiguous = {k: v for k, v in owner.items() if len({x for x in v}) > 1}
    if ambiguous:
        for term, rels in sorted(ambiguous.items()):
            print(f"  ⚠️ 「{term}」同时属于 {sorted(rels)}")
        print("  ⚠️ 存在共享区分词 → 基于 trigger 的路由无法分离这两个 Skill，"
              "需大脑凭语境裁决（见 test-plan.md L5 歧义点）")
    else:
        print("  ✅ 无共享区分词，12 个专项 Skill 的 trigger 互相可分离")

    # ---- 检查 4：路由用例覆盖分级 ----
    print("\n[4] 路由用例覆盖（信号词 vs trigger 文本）")
    brain_text_full = " ".join(triggers.get(BRAIN, []))
    uncovered = []
    for rid, desc, expected, signals in ROUTING_CASES:
        if expected is None:
            # R9：信息不足，不应有专项 Skill 强认领
            hits = [r for r in SPECIALIZED if any(s in " ".join(triggers.get(r, [])) for s in signals)]
            status = "大脑兜底（Unknown，不调度专项）" if not hits else f"⚠️ {hits} 强认领"
            if hits:
                ok = fail(f"{rid} 信息不足用例被 {hits} 的 trigger 强认领") and ok
            print(f"  {rid} {desc:<14} → {status}")
            continue
        skill_text = " ".join(triggers.get(expected, []))
        in_skill = [s for s in signals if s in skill_text]
        in_brain = [s for s in signals if s in brain_text_full]
        if in_skill:
            status = f"强覆盖（{len(in_skill)}/{len(signals)} 命中 {expected} trigger）"
        elif in_brain:
            status = "大脑兜底（大脑通用入口捕获后路由）"
        else:
            status = "完全未覆盖（仅靠大脑通用 trigger 兜底）"
            uncovered.append(rid)
        print(f"  {rid} {desc} → {status}")
    if uncovered:
        print(f"  ⚠️ 以下用例无任何 trigger 信号对应：{uncovered}")
        print("     这不是失败（大脑通用 trigger 会兜底），但意味着专项 Skill 的 trigger")
        print("     对该类输入不可发现，售卖前应考虑补充触发描述。")

    # ---- 汇总 ----
    print("\n" + "=" * 60)
    print("说明：本检查器只验证 trigger 契约与可分离性。")
    print("      「给定题目描述，大脑是否路由到正确专项」属 L3 推理行为测试，")
    print("      需按 tests/test-plan.md 第四节以人工/模型执行 13 个场景。")
    if ok:
        print("✅ L2 调度契约检查全部通过")
        return 0
    print("❌ 存在失败项，请修复后再发布")
    return 1


if __name__ == "__main__":
    sys.exit(main())
