#!/usr/bin/env python3
"""CTF-AI-Skill-Pack L1 危险内容扫描（可执行）

复用 ctf-skills skill_security_auditor.py 的思路（本 Pack 自研实现，无第三方依赖）。

**核心设计：区分「危险内容」与「CTF 教学性提及」。**
本 Pack 是安全教学产品，sqlmap / nmap / hydra / 反弹 shell 的**原理性提及**是
预期内容，不是危险内容。只有以下两类才算真问题：

  CRITICAL（直接失败）：
    - 破坏性命令（rm -rf /、fork bomb、curl|sh、dd of=/dev/…）
    - 硬编码凭据（AKIA*、私钥块、sk-*、GitHub PAT）

  HIGH（需在该文件内找到正当理由标注，否则失败）：
    - eval / exec 执行字符串字面量
    - 远程下载并执行

  EDUCATIONAL（允许，仅统计）：
    - 工具名教学性提及（sqlmap / nmap / hydra / nikto / dirb …）
    - 反弹 shell / 提权技术的原理讲解（在 CTF/靶场/授权语境下）

用法：python3 check_content.py [pack_root]
退出码：0 无 CRITICAL 且 HIGH 全部有正当理由 / 1 有失败项
"""

import io
import re
import sys
from pathlib import Path

if hasattr(sys.stdout, "buffer"):
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding="utf-8")

CRITICAL = {
    "破坏性命令 rm -rf": r"rm\s+-rf\s+/(?:\s|$|[*\s])",
    "fork bomb": r":\(\)\s*\{\s*:\|\s*:\s*&\s*\}\s*;:",
    "远程下载并执行": r"(?:curl|wget)\b[^\n|]*\|\s*(?:sh|bash|python)",
    "dd 写裸设备": r"dd\s+[^-\n]*of\s*=\s*/dev/(?:sd|nvme|hd)",
    "覆盖系统文件": r">\s*/etc/(?:passwd|shadow|sudoers)\b",
}

CREDENTIALS = {
    "AWS Access Key": r"AKIA[0-9A-Z]{16}",
    "私钥块": r"-----BEGIN\s+(?:RSA\s+|EC\s+|OPENSSH\s+)?PRIVATE KEY-----",
    "OpenAI 风格密钥": r"sk-[A-Za-z0-9]{20,}",
    "GitHub PAT": r"gh[pousr]_[A-Za-z0-9]{36,}",
}

HIGH = {
    "eval 执行字符串": r"\beval\s*\(",
    "exec 执行字符串": r"\bexec\s*\(",
}

# HIGH 命中时，同文件内出现这些标注即视为已说明正当理由
JUSTIFICATION = ("教学", "演示", "原理", "示例", "仅为演示", "CTF", "靶场", "授权")

# 授权边界横幅（出现在文件开头）：标明该文件为「授权受限内容」
# 有此横幅时，破坏性攻击技术（curl|sh、改写系统文件等）记录为「已标注」而非失败
AUTH_BANNER = "授权边界"


def split_code_fences(text: str):
    """把文本拆成 (prose_segments, in_fence_count)。

    eval/exec 出现在代码块里，通常是被讲解的漏洞代码本身（如 SSTI 的 eval(request.body)），
    属于教学内容；出现在散文里才像是「指示读者执行」。两者必须区分，否则误报淹没真问题。
    """
    prose, in_fence = [], False
    buf = []
    for line in text.splitlines():
        if line.strip().startswith("```"):
            if in_fence:
                in_fence = False
            else:
                in_fence = True
                if buf:
                    prose.append("\n".join(buf))
                    buf = []
            continue
        if not in_fence:
            buf.append(line)
    if buf:
        prose.append("\n".join(buf))
    # 行内代码 span（`exec()`、`Runtime.exec()` 等）是函数名提及，不是可执行指令
    return [re.sub(r"`[^`\n]*`", " ", seg) for seg in prose]

# 教学性提及（统计用，不失败）
EDUCATIONAL = r"\b(?:sqlmap|nmap|hydra|nikto|dirb|gobuster|wfuzz|hashcat|john|metasploit|msfconsole|burp|wireshark|tcpdump|gdb|objdump|radare2|ghidra|frida|angr)\b"
EDUCATIONAL_CN = r"反弹\s*shell|提权|越权|注入|后门|木马|字典攻击|暴力破解"


def scan_file(path: Path):
    text = path.read_text(encoding="utf-8")
    banner = AUTH_BANNER in text[:400]  # 横幅只在文件开头识别
    crit, cred, high_unjust, high_ok, edu = [], [], [], 0, set()

    # 破坏性/攻击性技术：有授权横幅 → 记录为已标注；无横幅 → 失败
    for label, pat in CRITICAL.items():
        for m in re.finditer(pat, text):
            if banner:
                high_ok += 1
            else:
                crit.append((label, m.group(0)))
    # 凭据泄漏：无论是否有横幅都失败——真实凭据不存在「教学豁免」
    for label, pat in CREDENTIALS.items():
        for m in re.finditer(pat, text):
            cred.append((label, m.group(0)))
    # eval/exec：只在代码块外的散文中命中才算问题
    for seg in split_code_fences(text):
        for label, pat in HIGH.items():
            for m in re.finditer(pat, seg):
                ctx = seg[max(0, m.start() - 200): m.end() + 200]
                if any(j in ctx for j in JUSTIFICATION):
                    high_ok += 1
                else:
                    high_unjust.append((label, m.group(0)))
    for m in re.finditer(EDUCATIONAL, text):
        edu.add(m.group(0))
    for m in re.finditer(EDUCATIONAL_CN, text):
        edu.add(m.group(0))
    return crit, cred, high_unjust, high_ok, edu


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
    print("CTF-AI-Skill-Pack L1 危险内容扫描")
    print("=" * 60)

    targets = sorted(skills.rglob("SKILL.md")) + sorted(skills.rglob("PLAYBOOK.md"))
    total_crit = total_cred = total_high_unjust = total_high_ok = 0
    all_edu: set[str] = set()

    for p in targets:
        crit, cred, high_unjust, high_ok, edu = scan_file(p)
        total_crit += len(crit)
        total_cred += len(cred)
        total_high_unjust += len(high_unjust)
        total_high_ok += high_ok
        all_edu |= edu
        rel = p.relative_to(root).as_posix()
        for label, hit in cred:
            ok = fail(f"{rel}：凭据泄漏 {label} → {hit}（不可豁免）") and ok
        for label, hit in crit:
            ok = fail(f"{rel}：破坏性技术无授权横幅 {label} → {hit}") and ok
        for label, hit in high_unjust:
            ok = fail(f"{rel}：散文 eval/exec 无理由 {label} → {hit}") and ok

    print(f"\n  扫描文件：{len(targets)} 个（SKILL.md + PLAYBOOK.md）")
    print(f"  凭据泄漏（AKIA / 私钥 / sk- / PAT）：{total_cred}  —— 要求 0，不可豁免")
    print(f"  破坏性技术无授权横幅：{total_crit}  —— 要求 0（有横幅的计入下方「已标注」）")
    print(f"  散文中 eval/exec 无理由：{total_high_unjust}  —— 要求 0")
    print(f"  已标注 / 授权受限 / 代码块内教学性：{total_high_ok} 处（允许）")
    print(f"  教学性提及（允许）：{len(all_edu)} 类工具/概念"
          f"{' —— ' + '、'.join(sorted(all_edu)[:8]) if all_edu else ''}")

    print("\n" + "=" * 60)
    if ok and total_crit == 0 and total_cred == 0 and total_high_unjust == 0:
        print("✅ L1 危险内容扫描通过")
        print("   无凭据泄漏；破坏性技术全部带授权横幅；散文无未标注的 eval/exec。")
        print("   代码块内的漏洞代码与教学性工具提及属预期内容，不告警失败。")
        return 0
    print("❌ 存在失败项，请修复后再发布")
    return 1


if __name__ == "__main__":
    sys.exit(main())
