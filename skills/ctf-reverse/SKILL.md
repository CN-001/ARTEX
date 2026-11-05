---
name: ctf-reverse
description: CTF Reverse 方向自主解题：ELF/PE 反汇编与反编译、算法识别与还原、动态调试与 patch、脱壳、Python/Java/Android 字节码还原、反调试绕过。适用于逆向题与二进制分析。
license: MIT
compatibility: 仅限 CTF / 靶场 / 授权环境
---

# CTF Reverse 自主解题

目标：**还原出正确输入**（通常就是 flag），不是写一份反汇编讲解。

## 1. 静态三件套

```bash
cd /app/data/tasks/TID/ && mkdir -p rev && cd rev
file ./chall; readelf -hW ./chall | head -14
strings -a -n 5 ./chall | grep -iE 'flag|correct|wrong|key|password|input' | head -30
objdump -d -M intel ./chall > disasm.asm; wc -l disasm.asm
r2 -A -q -c 'afl; izz~flag; q' ./chall
```

**先找比较点**：`strcmp/strncmp/memcmp/==` 附近，或"成功 / 失败"字符串的交叉引用。找到比较点就成功了一半。

## 2. 动态调试（比硬读反汇编快得多）

```bash
gdb -q ./chall -ex 'set pagination off' -ex 'b *main' -ex run -ex quit
gdb -q ./chall -ex 'b strcmp' -ex run -ex 'x/s $rdi' -ex 'x/s $rsi' -ex quit
ltrace ./chall 2>/dev/null | head -30
strace ./chall 2>/dev/null | head -30
```

- **期望值直接读寄存器 / 内存**。
- 遇到反调试（`ptrace`、`/proc/self/status` 检查）：patch 掉判断分支，或 `LD_PRELOAD` 挂钩。

## 3. 常见题型套路

| 特征 | 做法 |
|---|---|
| 逐字节与常量数组比较 | dump 数组，写逆变换还原 |
| 简单变换（异或 / 加减 / 查表） | 逆运算：`^k` 自逆，`+k` 用 `-k` |
| 标准算法变形（AES/RC4/Base64 变体） | 识别常量表（S-box、Base64 表），套对应实现 |
| 控制流平坦化 / 混淆 | 动态跑 + 断点抓中间值，别纯静态硬读 |
| `upx` / 加壳 | `upx -d ./chall`；其他壳先找 OEP dump |
| pyinstaller 打包 | `pyinstxtractor` + `decompyle3`/`uncompyle6` 还原 `.pyc` |
| Java / Android（jar/apk） | `jadx -d out ./app.apk` 读源码；`apktool d` 看资源与 smali |

```bash
decompyle3 ./chall.pyc 2>/dev/null || python3 -m uncompyle6 ./chall.pyc 2>/dev/null
chmod +x ./chall && ./chall
```

## 4. 还原后必须实证

```bash
echo "CANDIDATE_INPUT" | ./chall     # 本地必须返回成功 / flag
```

**只有程序亲自接受了输入，才算解出来。** 把「输入 → 原始回显」写进 `report_finding`。

## 5. 纪律

- 不逐行翻译反编译代码，只盯「输入 → 变换 → 比较」这条链。
- 算法不确定时先做**动态取值**：每步后断点读内存，比猜算法快。
- 拿到候选明文先按 flag 正则核对，再用程序复核。
