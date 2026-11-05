#!/usr/bin/env python3
"""CTF-AI-Skill-Pack 发布门禁（一键运行全部自研检查器）

  L1 check_skills.py    格式合规：13 个可发现技能 / 大脑 24 章 / 三章节 / 无悬空引用
  L1 check_content.py   危险内容：凭据泄漏 / 破坏性技术授权横幅 / 散文 eval-exec
  L2 check_dispatch.py  调度契约：trigger 可分离性 / 路由覆盖分级
  L4 check_contract.py  协作契约：专项 Skill 不重建推理系统

未覆盖（需人工/模型/多实例执行，见 test-plan.md）：
  L2 语义路由准确性、L3 推理行为、L5 跨实例稳定性

用法：python3 tests/run_checks.py [pack_root]
退出码：0 可发布 / 1 任一检查失败
"""

import io
import subprocess
import sys
from pathlib import Path

if hasattr(sys.stdout, "buffer"):
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding="utf-8")
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding="utf-8")

CHECKS = ["check_skills.py", "check_content.py", "check_dispatch.py", "check_contract.py"]


def main():
    here = Path(__file__).resolve().parent
    root = Path(sys.argv[1]) if len(sys.argv) > 1 else here.parent
    py = sys.executable or "python3"

    print("=" * 60)
    print(f"CTF-AI-Skill-Pack 发布门禁 — {root}")
    print("=" * 60)

    failed = []
    for name in CHECKS:
        script = here / name
        print(f"\n▶ {name}")
        print("-" * 60)
        r = subprocess.run([py, str(script), str(root)])
        print("-" * 60)
        if r.returncode != 0:
            failed.append(name)

    print("\n" + "=" * 60)
    for name in CHECKS:
        mark = "❌ 失败" if name in failed else "✅ 通过"
        print(f"  {mark}  {name}")
    print("=" * 60)
    if failed:
        print(f"❌ {len(failed)} 项检查失败：{failed} —— 不可发布")
        return 1
    print("✅ 全部自研检查通过（L1 格式 + L1 危险内容 + L2 调度契约 + L4 协作契约）")
    print("   仍未覆盖：L2 语义路由、L3 推理行为、L5 跨实例稳定性")
    print("   见 tests/test-plan.md，这些项需人工/模型/多实例执行。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
