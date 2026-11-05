---
name: ctf-crypto
description: CTF Crypto 方向自主解题：编码与古典密码识别、RSA 各类攻击（小指数、共模、Wiener、Fermat、广播）、AES 分组模式与 padding oracle、哈希爆破、z3 与格求解。适用于 Crypto 题与密码学分析。
license: MIT
compatibility: 仅限 CTF / 靶场 / 授权环境
---

# CTF Crypto 自主解题

## 1. 第一步永远是分类（不要跳）

```bash
python3 - <<'PY'
import base64, binascii
s = open('/app/data/tasks/TID/data.txt').read().strip()
print('len', len(s), 'charset', sorted(set(s))[:40])
for name, fn in [('b64', lambda: base64.b64decode(s + '===')),
                 ('b32', lambda: base64.b32decode(s + '===')),
                 ('hex', lambda: binascii.unhexlify(s))]:
    try: print(name, '->', fn()[:80])
    except Exception as e: print(name, 'x', type(e).__name__)
PY
```

**严格区分**：编码（无密钥、可逆）/ 加密（有密钥）/ 哈希（单向，只能爆破或查表）。
不要把 Base64 说成加密，不要说"解密 MD5"。

## 2. 编码与古典密码

```bash
ciphey -t "PAYLOAD" 2>/dev/null || true      # 多层编码链自动解
# 凯撒 / ROT 暴力
python3 -c "
s='PAYLOAD'
for k in range(26):
    print(k, ''.join(chr((ord(c)-65-k)%26+65) if c.isupper() else c for c in s))"
```

## 3. RSA（最高频）

**先算三件事**：`n` 的位数、`e` 的大小、是否给了多组 `(n,e,c)`。

| 条件 | 攻击 | 做法 |
|---|---|---|
| `e=3` 且 `m^3 < n` | 小指数开方 | `gmpy2.iroot(c, 3)` |
| `n` 太小 / 因子接近 | 分解 | `sympy.factorint` / Fermat 差平方 |
| 同一 `n` 两组 `(e1,e2)` 互素 | 共模攻击 | 扩展欧几里得求 `a,b`，算 `c1^a*c2^b` |
| 多组广播、`e` 小 | Håstad 广播 | CRT 合并后开 `e` 次方 |
| `e` 很大（接近 n） | Wiener / Boneh-Durfee | 连分数逼近 `d` |
| 已知 `p` 高位 / 低位 | Coppersmith | sage 或 z3 |

```python
from Crypto.Util.number import long_to_bytes
from gmpy2 import iroot
import sympy
m, exact = iroot(c, 3)
print(long_to_bytes(int(m)) if exact else 'no')
print(sympy.factorint(n))
```

```bash
sage -c "print(factor(2^64+1))"          # 若已装 sage（Coppersmith / 格）
python3 -c "from z3 import *; x=BitVec('x',32); s=Solver(); s.add(x*x==12345); print(s.check())"
```

## 4. AES / 分组密码

- ECB + 可控明文 → **分块重放 / 字节翻转**；
- CBC + 可控制密文 + 报错差异 → **padding oracle**；
- CBC + 固定 key 复用 → **bit-flipping**；
- 弱 key（短口令）→ 字典爆破。

## 5. 哈希

```bash
hashcat -m 0 -a 0 HASH ${WORDLISTS:-/opt/ctf/wordlists}/rockyou.txt --quiet 2>/dev/null \
  || john --wordlist=${WORDLISTS:-/opt/ctf/wordlists}/rockyou.txt hash.txt
```

## 6. 纪律

- 解出来先**验证**：明文是否含 flag 格式，或用题目的加密函数能否复现密文。
- 多层编码要解到底：每层结果都检查是否又是编码（可能 5 层 base64）。
- 数学攻击前先写清已知量与未知量，条件写错会白跑很久。
