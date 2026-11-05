---
name: crypto-analyzer
version: 1.0.0
description: CTF Crypto 题专项 playbook —— 在 CTF-Analyzer-V2.0 判定题目为 Crypto 后被调用。提供编码/加密/哈希严格区分、密码学攻击决策树与 ~150 种具名攻击索引（含出处 CTF 与年份）。本 Skill 是知识层，不实现自己的推理系统；假设、证据、置信度、提示等级一律由 CTF-Analyzer-V2.0 统一维护。
trigger:
  - 由 CTF-Analyzer-V2.0 在判定题目类型为 Crypto 后调度调用
  - 用户要求分析密码学/编码内容（Base64、古典密码、RSA、AES、哈希、格密码等）
  - 用户提供密文、密钥参数（n/e/c）、哈希值并希望深入攻击方向
---

# ctf-crypto —— CTF Crypto 题专项 Playbook

## 一、核心功能

| 职责 | 归属 |
|---|---|
| 判断题目类型、维护 Session State、Level 1/2/3、假设追踪与优先级 | **CTF-Analyzer-V2.0（大脑）** |
| 编码/加密/哈希判定、具体攻击技术、数学原理、工具命令 | **本 Skill（playbook）** |
| 决定下一步验证什么、解读测试结果 | **CTF-Analyzer-V2.0** |
| 被询问具体攻击原理时提供纵深技术内容 | **本 Skill** |

**本 Skill 不输出题目类型判断、不维护假设状态、不自行宣布 Level 等级。**

## 二、输入识别

| 输入类型 | 用途 |
|---|---|
| 密文字符串 | 字符集/长度/格式分析 |
| 参数三元组（n, e, c）等 | RSA 与公钥密码方向 |
| 哈希值 | 定长摘要识别 → 哈希破解方向 |
| 源码中的加密/编码逻辑 | 自实现算法的缺陷分析 |
| 已知明文/密文对 | 已知明文攻击、密钥恢复 |
| 题目提示（如「因子接近」） | 结构性缺陷线索 |

## 三、核心纪律：编码 / 加密 / 哈希必须严格区分

**这是 Crypto 题最常见的初学者错误，也是本 Skill 的第一条硬规则。**

| 概念 | 特征 | 是否可逆 | 正确表述 |
|---|---|---|---|
| **编码** | 无密钥（Base64 / Hex / URL / ROT13 / Morse） | 可逆，无需密钥 | 「解码」 |
| **加密** | 有密钥（RSA / AES / DES） | 可逆，需密钥 | 「解密」 |
| **哈希** | 单向（MD5 / SHA-1 / SHA-256） | **不可逆** | **「哈希破解 / 彩虹表 / 撞库」——禁止说「解密」** |

## 四、分析流程：Crypto 决策树

```text
输入数据（字符串 / 密文 / 参数组）
    │
    ↓  字符集 / 长度 / 格式分析
    │
    ↓  是否像编码？（无密钥、可打印、字符集匹配标准编码）
    │
    ├── 无证据 → Unknown
    │
    ├── 是 → 编码识别（Base64 / Base32 / Hex / URL / ROT13 / Morse / Base85）
    │     │
    │     ↓  解码后是否得到可读结果？
    │     │
    │     ├── 是 → 是否已是 flag？若不是，检查第二层编码
    │     └── 否 → 回到格式分析，或进入「非编码」分支
    │
    └── 否 → 判断数据性质
          │
          ├── 随机二进制 → 加密数据 / 压缩数据
          ├── 提供了 n、e、c → RSA 方向
          │     └── 详见下方「RSA 缺陷速查表」
          ├── 固定块长度 + IV / 模式特征 → 对称加密（AES / DES / 流密码）
          │     └── ECB 模式重放 / Padding Oracle / IV 重用 / 弱密钥
          ├── 定长摘要（16/20/32/64 字节 hex）→ Hash 方向
          │     └── 彩虹表 / 撞库 / 长度扩展 / 盐值缺失
          ├── 保留自然语言统计特征 → 古典密码
          │     └── 词频 / 重合指数（IC） / Kasiski 测试
          └── 证据不足 → Unknown
```

### RSA 缺陷速查表（按可利用性排序）

| 缺陷条件 | 攻击方向 | 参考文档 |
|---|---|---|
| 两素因子接近 | Fermat 分解 | `rsa-attacks.md` |
| 小 e + 短消息 | 低指数攻击 / 立方根 | `rsa-attacks.md` |
| 多组密文同 e 同明文 | Hastad 广播 | `rsa-attacks.md` |
| d 较小 | Wiener 攻击 | `rsa-attacks-2.md` |
| 共享因子（多组 n） | GCD 分解 | `rsa-attacks.md` |
| 共模（同 n 同明文不同 e） | 共模攻击 | `rsa-attacks.md` |
| 无填充 / 可预测填充 | 选择密文、Franklin-Reiter | `rsa-attacks-2.md` |
| CRT 实现缺陷 | 故障注入 / 位翻转恢复 | `rsa-attacks-2.md`（CSAW CTF 2016） |

### 对称密码缺陷速查

| 缺陷 | 信号 |
|---|---|
| ECB 模式 | 相同明文块 → 相同密文块；模式可见 |
| Padding Oracle | 解密失败与填充错误的响应可区分 |
| IV 重用 | 相同明文 → 相同密文 |
| 弱密钥 / 硬编码密钥 | 源码中直接出现密钥常量 |

## 五、输出格式

返回给 CTF-Analyzer-V2.0 的结果（**不是面向用户的最终格式**）：

```markdown
## Crypto 专项分析结果

### 编码/加密/哈希判定
<类别> + <判定依据的 Evidence>

### 决策树命中分支
<节点> ← <Evidence>

### 特征观察
- 字符集：<>
- 长度：<>
- 结构特征：<>

### 建议假设（供 Analyzer 纳入）
- <假设名>：<技术依据>

### 可选攻击方法（按验证成本从低到高）
1. <方法> —— 成本：<>, 预期成功信号：<>
2. ...

### 风险与不确定性
- <哪些是推断，哪些缺证据>

### 参考语料
- <相关 .md>
```

## 六、与 CTF-Analyzer-V2.0 的协作方式

```text
用户输入
  → CTF-Analyzer-V2.0（识别类型、提取 Evidence、建立假设）
  → 判定为 Crypto
  → 调用 ctf-crypto playbook
  → 本 Skill 返回判定结果与攻击方法索引
  → Analyzer 汇总（更新假设、计算优先级、生成 Smart Next Step）
  → Level 1/2/3 输出
```

**协作纪律：**

1. 只被 Analyzer 调用，不反向调用。
2. 不修改 Session State。
3. 编码/加密/哈希的判定结论必须附带 Evidence。
4. 证据不足以走决策树时返回 Unknown + 所需信息清单，**不编造明文或 flag**。
5. 未实际完成计算前，不宣称攻击成功。

## 七、安全边界

- 本 Skill 内容用于 CTF、靶场、密码学课程与授权研究。
- 对真实系统凭证（如真实服务器的密钥、证书）的破解指引不予提供。
- 越界时转为密码学原理讲解、靶场复现方案与防御建议（如正确选择填充模式、密钥管理规范）。

## 八、自检规则

- [ ] 编码 / 加密 / 哈希是否严格区分？是否出现「把 MD5 说成解密」这类错误？
- [ ] 决策树分支走法是否有 Evidence 支撑？
- [ ] 是否编造了明文、密钥、flag 或计算结果？
- [ ] 攻击方法是否按成本排序并给出成功信号？
- [ ] 是否返回假设供 Analyzer 纳入，而非自己宣布结论？

## 九、参考语料索引

**ctf-skills 语料（Lukasz Jagiello，MIT）：** `rsa-attacks*.md`、`modern-ciphers*.md`、`classic-ciphers.md`、`ecc-attacks.md`、`dh-attacks.md`、`lattice-and-lwe.md`、`post-quantum.md`、`prng*.md`、`stream-ciphers.md`、`exotic-crypto*.md`、`zkp-and-advanced.md`、`advanced-math.md`、`historical.md`

**hack-skills 语料（VillanCh，MIT）：** `classical-cipher-analysis/`（含 Kasiski / IC / chi-squared 等 ~15 个可运行 Python 实现）、`rsa-attack-techniques/`、`hash-attack-techniques/`、`symmetric-cipher-attacks/`、`lattice-crypto-attacks/`

## 十、本地工具资产（tools/）

`tools/` 内置 RSA 攻击脚本库与密码工具（用户本机 CTF 工具集）。**RSA 题命中攻击方向后，先查 `tools/README.md` 取对应脚本**；脚本内 `n/e/c/dp/p/q` 为硬编码示例，将题目参数替换进变量后 `python 脚本名.py` 运行（依赖 pycryptodome / gmpy2 / libnum）。**工具权威副本固定位于 `E:\ctf\工具，脚本`（本机稳定路径）；若 `tools/` 目录缺失或被应用周期性清空，直接从该路径调用；并可运行 `python E:\ctf\restore_ctf_tools.py` 一键恢复全部工具到实时技能。**

| 攻击类型 | 触发条件 | 工具/脚本 |
|---|---|---|
| 大数分解 | n 过大 / 需分解 | `tools/yafu-1.34(分解N)/yafu-x64.exe`（`factor(N)`；大数 `yafu-x64 "factor(@)" -batchfile n.txt`） |
| dp,dq 泄露 | 已知 `p,q,dp,dq,c` | `tools/dp,dq/1.py` |
| dp 泄露 | 已知 `n,e,dp,c` | `tools/dp泄露/1.py` |
| gcd(e,φ)>1 且 e 非素数 | `gmpy2.gcd(e,phi_n)>1` | `tools/e和ph_ni不互素，e不是素数/1.py` |
| n 与 c 不互素 | `gcd(n,c)=p` | `tools/nc不互素/1.py` |
| n 多因子（p·q·r…） | 多素数乘积 | `tools/n分解多个不一样的/例题/baigei2.py` |
| n=p^r | n 为素数次幂 | `tools/n是p的r次方/` |
| Rabin（e=2） | e=2 | `tools/rabin(e=2)/全给的.py` |
| Wiener | d 小 / e 很大 | `tools/winner维纳（e很大）脚本n÷e约等于1/1.py` |
| 小 e 攻击 | e=3 或很小 | `tools/低加密指数攻击（e很小或e=3）/1.py` |
| Hastad 广播 | 多组 `(n,c)` 同小 e | `tools/低加密指数广播攻击（…中国剩余定理求解）/1.py` |
| 共模攻击 | 同 n，e1,e2 互质 | `tools/共模脚本/e1e2.py` |
| 共享素数 | 多 n 共享素数 | `tools/共享素数脚本/2个c.py`、`3个n,e,c.py` |
| PEM 公钥读取 | 题目给 `.pem` | `tools/公私钥脚本/公钥读取.py` |
| 模逆元求 d | 已知 e、φ | `tools/求模逆元d/` |
| Coppersmith 类 | d 低位 / m 高位 / p 高位 / 维纳变形 / 广播 | `tools/sage脚本_*/`（需 SageMath 环境，见 README） |
| 编码/古典密码 | 需要快速试编码 | `tools/CaptfEncoder-3.1.2-bin/CaptfEncoder-win-x64-3.1.2.exe` |
| 凯撒爆破 | 凯撒类题 | `tools/KaisaPro.exe`（GUI） |
| XOR 爆破 | XOR 类题 | `tools/XorBrute.exe`（GUI） |
