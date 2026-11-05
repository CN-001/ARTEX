# Win-Offensive-Dev — Windows 免杀开发技能系统

> 信益安信息安全研究院出品 | XinYiAn Security Research Institute

Windows 攻击性安全开发的完整知识库与自动化工具链。包含 shellcode 加载、进程注入、API 隐藏、syscall、防御规避、内存逃逸等 13 个技术模块，以及自动化 loader 生成（/malfuzz）和 Sysinternals 分析（/reverse）两个核心命令。

---

## 快速开始

### 环境要求

- [Claude Code](https://claude.ai/code) CLI 或 IDE 扩展
- Windows 10/11 x64
- MSVC 编译器（Visual Studio Build Tools 或完整 VS）
- 项目目录下的 `SysinternalsSuite/` 文件夹（已包含）
- 项目目录下的 `calc.bin`（x64 calc shellcode，用于验证）

### 安装 Skill

Skills 通过 `.claude/skills/` 目录自动加载，无需手动注册。

```
win-offensive-dev/
├── .claude/
│   └── skills/
│       ├── win-offensive-dev/    # 主知识库 skill（自动触发）
│       │   └── SKILL.md
│       ├── malfuzz/              # /malfuzz 命令
│       │   └── SKILL.md
│       ├── reverse/              # /reverse 命令
│       │   └── SKILL.md
│       └── edrfuzz/              # /edrfuzz 命令
│           └── SKILL.md
├── edr_profiles/                 # EDR 检测特征知识库（自动生成）
├── references/                   # 技术参考文件
│   ├── execution.md
│   ├── iat_hiding.md
│   ├── syscalls.md
│   ├── defense_evasion.md
│   ├── process_manipulation.md
│   ├── memory_evasion.md
│   ├── credential_access.md
│   ├── persistence.md
│   ├── networking.md
│   ├── reflective_loading.md
│   ├── token_privilege.md
│   ├── malfuzz.md                # /malfuzz 技术池与模板
│   └── reverse_analysis.md       # /reverse 分析流程
├── SysinternalsSuite/            # Sysinternals 工具集
├── calc.bin                      # 验证用 calc shellcode
├── SKILL.md                      # 主 skill（兼容旧版加载）
└── README.md
```

**加载方式**：在此项目目录下启动 Claude Code 即可。Skills 会被自动发现和加载。

```bash
# 进入项目目录
cd C:\Users\你的用户名\Desktop\AdaptixC2\win-offensive-dev

# 启动 Claude Code
claude
```

启动后输入 `/` 可以看到可用的 skill 列表，包括 `malfuzz` 和 `reverse`。

---

## 命令使用

### /malfuzz — 批量生成免杀 Loader

自动生成多个技术组合不同的 shellcode loader EXE，用于测试哪种组合能绕过目标环境的检测。

#### 语法

```
/malfuzz <语言> <数量> <shellcode路径> [sign] [sandbox] ["描述"]
```

| 参数 | 说明 |
|------|------|
| 语言 | `C`、`C++` 或 `Rust` |
| 数量 | 生成 loader 的数量（建议 3-10） |
| shellcode路径 | .bin 文件的绝对路径或相对路径 |
| sign | 可选。启用签名掠夺（从 Program Files 复制签名/图标/版本信息） |
| sandbox | 可选。启用反沙箱+反调试技术 |
| 描述 | 可选。用引号包裹，指定技术偏好 |

#### 示例

```
# 基础：生成 5 个 C++ loader
/malfuzz C++ 5 ./payload.bin

# 带签名掠夺：看起来像合法软件
/malfuzz C++ 5 ./payload.bin sign

# 全套：签名 + 反沙箱
/malfuzz C 8 ./payload.bin sign sandbox

# 指定偏好：只用 callback 执行方式
/malfuzz C++ 3 ./payload.bin "只用callback执行，不要fiber"

# Rust 版本
/malfuzz Rust 4 ./payload.bin sandbox
```

#### 执行流程

1. **生成阶段**：使用 `calc.bin` 作为 payload 生成所有 loader 源码并编译
2. **自动验证**：逐个运行 EXE，自动检测 calc.exe 进程是否启动（3-5秒超时）
   - ✅ 检测到 calc.exe → 验证通过，自动 kill
   - ❌ 未检测到 → 报告失败原因，提供修复
3. **替换阶段**：验证通过的 loader 自动替换为真实 shellcode，重新加密和编译
4. **输出**：所有文件在 `fuzz_output_<时间戳>/` 目录下，附带验证结果

#### 输出结构

```
fuzz_output_20260507_143022/
├── README.md                    # 汇总表格
├── common/
│   └── peb_resolve.h            # 共享 PEB walk 头文件
├── loader_01_fiber_xor.c        # 源码
├── loader_01_fiber_xor.exe      # 编译产物
├── loader_01.md                 # 技术说明
├── loader_02_callback_aes.c
├── loader_02_callback_aes.exe
├── loader_02.md
└── ...
```

#### 技术组合维度

每个 loader 从 7 个维度随机选取技术：

| 维度 | 可选项 |
|------|--------|
| 加密 | XOR (16-32 byte), AES-256-CBC, RC4 |
| 执行 | Fiber, Callback (5种), APC self-inject, Mapping injection |
| API解析 | Djb2 hash, CRC32, Jenkins |
| 字符串 | Stack strings, XOR encrypted |
| 防御规避 | ETW HBP bypass, AMSI HBP bypass, 无 |
| 反分析 | Anti-debug, Anti-VM, Delay, API hammer, 无 |
| 熵值 | 英文填充, 字典填充, 无 |

#### 硬性约束（所有 loader 必须满足）

- IAT 只有 3-6 个无害导入函数
- 二进制中零明文字符串
- 所有敏感 API 通过 PEB walk + hash 动态解析
- 内存 RW→RX 分离（不使用 RWX）
- 无 CRT 依赖
- Shannon 熵值 < 6.5

---

### /reverse — 恶意软件分析与自动增强

使用 Sysinternals 工具对 EXE 进行静态+动态分析，评估检测风险，自动生成改进方案。

#### 语法

```
/reverse <路径>
```

- 路径为单个 .exe 文件：直接分析该文件
- 路径为文件夹：扫描文件夹下所有 .exe 逐个分析

#### 示例

```
# 分析单个 loader
/reverse ./fuzz_output_20260507_143022/loader_01_fiber_xor.exe

# 分析整个 fuzz 输出目录
/reverse ./fuzz_output_20260507_143022/

# 分析任意 EXE
/reverse C:\tools\my_implant.exe
```

#### 分析流程

**Phase 1 — 静态分析：**
- `strings64.exe`：提取可读字符串，检查是否有敏感 API/DLL 名泄露
- `sigcheck64.exe`：检查签名状态、计算熵值、读取版本信息
- PE 导入表解析：验证 IAT 是否干净
- 按 section 计算 Shannon 熵值

**Phase 2 — 动态分析：**
- `Procmon64.exe`：自动启动监控 → 运行目标 EXE → 捕获行为 → 导出分析
- 检查：注册表操作、文件写入、网络连接、子进程创建、内存保护变更

**Phase 3 — 报告生成：**
- 每个 EXE 生成独立的 markdown 分析报告
- 每项检查标记 ✅ PASS / ❌ FAIL / ⚠️ WARNING
- 计算整体风险评分

**Phase 4 — 自动增强：**
- 如果发现问题（❌ FAIL），自动映射到修复方案
- 进入 plan mode 展示具体代码修改建议
- 用户确认后自动修复、重新编译、再次分析

#### 检测规则示例

| 发现 | 严重性 | 自动修复 |
|------|--------|----------|
| strings 检出 API 名 | ❌ Critical | 改用 stack strings |
| 熵值 > 6.5 | ❌ Fail | 添加英文填充段落 |
| IAT 有 VirtualAlloc | ❌ Fail | 改用 PEB walk 解析 |
| IAT > 6 个导入 | ⚠️ Warning | 移除多余导入 |
| Procmon 检测到网络连接 | ⚠️ Warning | 预期行为（staged loader） |

---

### /edrfuzz — EDR 检测特征收集

在 `/malfuzz` 生成 loader 并在目标环境测试后，使用 `/edrfuzz` 收集 EDR 检测结果，建立 EDR 检测特征知识库。

#### 语法

```
/edrfuzz [EDR名称]
```

#### 交互流程

1. 自动扫描最近的 `fuzz_output_*` 目录，列出所有 loader
2. 询问用户哪些 loader 未被 EDR 检测到（输入编号）
3. 询问 EDR 产品名称
4. 自动分析绕过/检测的技术共性
5. 保存到 `edr_profiles/<edr_name>/profile.md`

#### 示例

```
/edrfuzz
# 交互式：列出 loader → 输入未检测编号 → 输入 EDR 名称

/edrfuzz CrowdStrike
# 跳过 EDR 名称询问，直接使用 CrowdStrike
```

#### 输出结构

```
edr_profiles/
├── README.md                    # 所有 EDR 汇总表
├── crowdstrike/
│   └── profile.md              # CrowdStrike 检测特征
├── defender/
│   └── profile.md              # Defender 检测特征
└── sentinelone/
    └── profile.md              # SentinelOne 检测特征
```

#### 知识库用途

- 多次测试后自动分析 EDR 检测趋势
- 推荐针对特定 EDR 最有效的技术组合
- 为后续 `/malfuzz` 提供技术偏好建议

---

## 主知识库 (win-offensive-dev)

主 skill 在对话中自动触发，当你讨论以下话题时会自动加载相关技术参考：

- Shellcode 加载与执行
- API 解析与 IAT 隐藏
- Syscall 与 Hook 规避
- 进程操纵（hollowing, injection, spoofing）
- 内存逃逸（sleep obfuscation, anti-debug）
- 防御规避（ETW/AMSI bypass）
- 持久化机制
- C2 通信
- 反射加载
- Token 操纵与提权

每个技术都包含 OPSEC 元数据（检测风险评级、Sysmon Event ID、ETW Provider、EDR Hook 点、缓解建议）。

---

## 典型工作流

```
1. /malfuzz C++ 5 ./beacon.bin sign sandbox
   → 生成 5 个不同技术组合的 loader，自动验证 calc

2. /reverse ./fuzz_output_20260507/
   → 对所有 loader 进行静态+动态分析

3. 查看报告 → 自动生成增强方案 → 确认修复 → 重新编译

4. 将通过 /reverse 的 loader 部署到目标环境测试

5. /edrfuzz
   → 收集 EDR 检测结果，记录哪些技术组合绕过了目标 EDR

6. 基于 edr_profiles 知识库，下次 /malfuzz 时指定有效技术偏好:
   /malfuzz C++ 5 ./beacon.bin sign sandbox "只用Fiber和Mapping，RC4加密"
```

---

## 编译环境配置

确保 MSVC 编译器在 PATH 中：

```powershell
# 方法 1：使用 Developer Command Prompt
# 从开始菜单打开 "x64 Native Tools Command Prompt for VS 2022"

# 方法 2：手动加载环境
& "C:\Program Files\Microsoft Visual Studio\2022\Community\VC\Auxiliary\Build\vcvars64.bat"

# 验证
cl /?
```

Rust 环境（如果使用 Rust loader）：

```powershell
# 安装 Rust
winget install Rustlang.Rust.MSVC

# 验证
rustc --version
cargo --version
```

---

## 注意事项

- 所有生成的 loader 仅用于授权的安全测试和研究
- `calc.bin` 验证阶段确保 loader 逻辑正确后再替换真实 payload
- 动态分析（Procmon）需要管理员权限
- 签名掠夺生成的是无效签名（证书不匹配），仅用于绕过简单的签名检查
- 建议在隔离的测试环境中运行生成的 EXE

---

> 信益安信息安全研究院出品 | XinYiAn Security Research Institute
