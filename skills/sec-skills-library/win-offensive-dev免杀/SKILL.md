---
name: win-offensive-dev
description: Windows offensive security development skill for C/C++ and Rust. Use when writing shellcode loaders, process injection, API hashing, IAT hiding, syscall stubs, NTDLL unhooking, sleep obfuscation, reflective loaders, BOF loaders, ETW/AMSI bypass, or any Windows red team tooling code. Also use when user says "单文件嵌入" (single-file embedded shellcode). Trigger on: shellcode, injection, evasion, loader, syscall, unhooking, PEB walk, export table, API hash, IAT hiding, process hollowing, APC injection, fiber execution, callback execution, Early Bird, PPID spoofing.
---

# Windows Offensive Development

Complete knowledge base for writing Windows offensive security tools in C/C++ and Rust. All techniques are self-contained implementations — no external references needed.

## Workflow

### 1. Understand the Request
Identify which technique category the user needs:
- Shellcode loading & execution
- Process injection (remote)
- API resolution & IAT hiding
- Syscalls & hook evasion
- Anti-analysis & anti-debugging
- Memory evasion (sleep obfuscation, encryption)
- Credential access & enumeration

### 2. Select Technique
Read the appropriate reference file from `references/` for implementation details.

### 3. Generate Code
- **C/C++**: Use MSVC-compatible code with `windows.h`, `winternl.h`, `tlhelp32.h`
- **Rust**: Use `windows-sys` crate with `unsafe` FFI blocks
- Apply IAT hiding by default for sensitive APIs
- XOR-encrypt embedded payloads with random keys

### 4. Generate README
**After every batch of loaders/tools generated**, create a `README.md` in the same output directory containing:
- Table of all generated .exe files with technique name
- Common features shared by all
- Key differences between each variant
- C++ vs Rust comparison if both languages used

README template:
```markdown
# [Directory Name] - [Category]

> Generated: [date] | Payload: [name] ([size]) | Encryption: [method]

## Common Features (All N)
- [Execution method]
- [Encryption method]
- [Module resolution method]

## Loaders

| File | Language | Technique | IAT Content | Key Difference |
|------|----------|-----------|-------------|----------------|
| name.exe | C++ | ... | ... | ... |

## Key Differences
...
```

## Topic Index

### Shellcode Loading & Execution
- Local execution: VirtualAlloc + CreateThread, Fiber, Callback, APC self-inject
- Remote injection: Early Bird APC, Classic injection, Mapping injection
- **Reference**: `references/execution.md`

### API Resolution & IAT Hiding
- PEB walk (custom GetModuleHandle)
- Export table walk (custom GetProcAddress)
- Hash algorithms: Djb2, CRC32, JenkinsOneAtATime
- Compile-time constexpr hashing
- Stack strings, XOR-encrypted strings, dispatch tables
- **Reference**: `references/iat_hiding.md`

### Syscalls & Hook Evasion
- Direct syscalls, SSN extraction, Hell's Gate
- Indirect syscalls (HellsHall)
- NTDLL unhooking (disk, KnownDlls, suspended process)
- **Reference**: `references/syscalls.md`

### Process Manipulation
- Process hollowing, ghost injection, herpaderping
- Threadless injection, module stomping/overloading
- PPID spoofing, argument spoofing, Block DLL policy
- **Reference**: `references/process_manipulation.md`

### Memory Evasion
- Sleep obfuscation: Ekko, Zilean, Foliage
- PE fluctuation, heap encryption
- Payload encryption: XOR, RC4, AES (BCrypt)
- **Reference**: `references/memory_evasion.md`

### Defense Evasion
- ETW bypass (byte patching, HBP patchless, session hijacking)
- AMSI bypass (byte patching, HBP patchless)
- Binary metadata modification, entropy reduction
- CRT removal, IAT camouflage
- **Reference**: `references/defense_evasion.md`

### Credential Access
- LSASS dumping (MiniDump, handle duplication, SilentProcessExit)
- SAM dumping (local, remote, from disk)
- Browser credential extraction (Chrome/Firefox cookies & passwords)
- **Reference**: `references/credential_access.md`

## Single-File Loader Pattern

When user says "单文件嵌入" or "single file embedded":
1. Read the .bin shellcode file
2. Generate random XOR key (16-32 bytes)
3. Encrypt shellcode with key
4. Embed both encrypted bytes and key as C/Rust arrays in source
5. At runtime: decrypt → allocate → execute
6. Output is one standalone .exe with no external dependencies

## Code Style

### C/C++ Conventions
```c
typedef BOOL(WINAPI* tCreateProcessA)(LPCSTR,LPSTR,...);
HMODULE hK32 = PebMod(HASH_KERNEL32);
auto fCPA = (tCreateProcessA)ProcH(hK32, djb2c("CreateProcessA"));
if (!fCPA) return 1;
```

### Rust Conventions
```rust
type CPA = unsafe extern "system" fn(*const u8, *mut u8, ...) -> i32;
let f: CPA = std::mem::transmute(proc_djb2(base, djb2(b"CreateProcessA")));
```
