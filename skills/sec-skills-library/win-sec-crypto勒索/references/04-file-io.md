# 04 — File I/O & IOCP

## File Pointer Helpers
```c
BOOL SetFilePointerFwd(HANDLE hFile, DWORDLONG dwOffset, DWORD dwMoveMethod) {
    LARGE_INTEGER li = { 0 }; li.QuadPart = (LONGLONG)dwOffset;
    return SetFilePointerEx(hFile, li, NULL, dwMoveMethod);
}
BOOL SetFilePointerBack(HANDLE hFile, DWORDLONG dwOffset, DWORD dwMoveMethod) {
    LARGE_INTEGER li = { 0 }; li.QuadPart = -((LONGLONG)dwOffset);
    return SetFilePointerEx(hFile, li, NULL, dwMoveMethod);
}
```

## Sync File Read/Write Helpers
```c
BOOL ReadFileFromDiskW(IN LPCWSTR szPath, OUT PBYTE* ppBuf, OUT PDWORD pdwSize) {
    HANDLE hFile = CreateFileW(szPath, GENERIC_READ, FILE_SHARE_READ, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
    if (hFile == INVALID_HANDLE_VALUE) return FALSE;
    *pdwSize = GetFileSize(hFile, NULL); DWORD dwRead = 0;
    *ppBuf = (PBYTE)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, *pdwSize);
    BOOL ok = ReadFile(hFile, *ppBuf, *pdwSize, &dwRead, NULL) && (dwRead == *pdwSize);
    CloseHandle(hFile);
    if (!ok) { HeapFree(GetProcessHeap(), 0, *ppBuf); *ppBuf = NULL; }
    return ok;
}

BOOL WriteFileToDiskW(IN LPCWSTR szPath, IN PBYTE pBuf, IN DWORD dwSize) {
    HANDLE hFile = CreateFileW(szPath, GENERIC_READ | GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, NULL);
    if (hFile == INVALID_HANDLE_VALUE) return FALSE;
    DWORD dwWritten = 0;
    BOOL ok = WriteFile(hFile, pBuf, dwSize, &dwWritten, NULL) && (dwWritten == dwSize);
    CloseHandle(hFile);
    return ok;
}
```

---

## Async Overlapped I/O with IOCP

### ASYNC_CHUNK Structure
```c
#define MAX_IN_FLIGHT 16

typedef enum { OP_READ = 1, OP_WRITE = 2 } CHUNK_OP;

typedef struct _ASYNC_CHUNK {
    OVERLAPPED  Overlapped;
    PBYTE       pBuffer;
    DWORD       dwBytesXfer;
    BOOL        bInUse;
    BOOL        bEncrypted;
    INT         iChunkIndex;
    CHUNK_OP    opType;
    ULONGLONG   u64Offset;
} ASYNC_CHUNK, *PASYNC_CHUNK;
```

### Post Async Read
```c
BOOL PostReadOp(INT idx, PASYNC_CHUNK pArr, PDWORD pdwActive,
    PULONGLONG pullOffset, ULONGLONG ullFileLen, HANDLE hFile)
{
    if (*pullOffset >= ullFileLen) return FALSE;
    PASYNC_CHUNK pChunk = &pArr[idx];
    pChunk->bInUse = TRUE;
    pChunk->opType = OP_READ;
    pChunk->u64Offset = *pullOffset;

    DWORD dwToRead = (DWORD)min(CHUNK_SIZE, ullFileLen - *pullOffset);
    pChunk->dwBytesXfer = dwToRead;

    RtlSecureZeroMemory(&pChunk->Overlapped, sizeof(OVERLAPPED));
    pChunk->Overlapped.Offset     = (DWORD)(pChunk->u64Offset & 0xFFFFFFFF);
    pChunk->Overlapped.OffsetHigh = (DWORD)(pChunk->u64Offset >> 32);

    *pullOffset += dwToRead;
    (*pdwActive)++;

    ReadFile(hFile, pChunk->pBuffer, dwToRead, NULL, &pChunk->Overlapped);
    return TRUE;
}
```

### IOCP Main Loop Pattern
```c
DWORD HandleFileEncryptionAsync(IN LPWSTR szFilePath) {
    HANDLE hFile = CreateFileW(szFilePath, GENERIC_READ | GENERIC_WRITE, 0, NULL,
        OPEN_EXISTING, FILE_FLAG_OVERLAPPED, NULL);

    HANDLE hIOCP = CreateIoCompletionPort(INVALID_HANDLE_VALUE, NULL, 0, 0);
    CreateIoCompletionPort(hFile, hIOCP, (ULONG_PTR)hFile, 0);

    LARGE_INTEGER liSize; GetFileSizeEx(hFile, &liSize);
    ULONGLONG ullFileLen = liSize.QuadPart, ullOffset = 0;
    DWORD dwActive = 0;

    ASYNC_CHUNK ChunkArr[MAX_IN_FLIGHT] = { 0 };
    for (int i = 0; i < MAX_IN_FLIGHT; i++)
        ChunkArr[i].pBuffer = (PBYTE)LocalAlloc(LPTR, CHUNK_SIZE);

    for (int i = 0; i < MAX_IN_FLIGHT && ullOffset < ullFileLen; i++)
        PostReadOp(i, ChunkArr, &dwActive, &ullOffset, ullFileLen, hFile);

    while (dwActive > 0) {
        DWORD dwBytes; ULONG_PTR ulKey; LPOVERLAPPED pOv;
        GetQueuedCompletionStatus(hIOCP, &dwBytes, &ulKey, &pOv, INFINITE);
        PASYNC_CHUNK pC = CONTAINING_RECORD(pOv, ASYNC_CHUNK, Overlapped);

        if (pC->opType == OP_READ) {
            pC->opType = OP_WRITE;
            RtlSecureZeroMemory(&pC->Overlapped, sizeof(OVERLAPPED));
            pC->Overlapped.Offset     = (DWORD)(pC->u64Offset & 0xFFFFFFFF);
            pC->Overlapped.OffsetHigh = (DWORD)(pC->u64Offset >> 32);
            WriteFile(hFile, pC->pBuffer, dwBytes, NULL, &pC->Overlapped);
        } else {
            pC->bInUse = FALSE;
            dwActive--;
            for (int i = 0; i < MAX_IN_FLIGHT; i++) {
                if (!ChunkArr[i].bInUse && ullOffset < ullFileLen) {
                    PostReadOp(i, ChunkArr, &dwActive, &ullOffset, ullFileLen, hFile);
                    break;
                }
            }
        }
    }

    for (int i = 0; i < MAX_IN_FLIGHT; i++) LocalFree(ChunkArr[i].pBuffer);
    CloseHandle(hIOCP); CloseHandle(hFile);
    return 0;
}
```

---

## BackupRead/BackupWrite (timestamp-preserving I/O)

```c
#define BACKUP_CHUNK (64 * 1024)
#define STREAM_HDR_FIXED (FIELD_OFFSET(WIN32_STREAM_ID, cStreamName))

BOOL EnableBackupPrivileges(VOID) {
    HANDLE hToken; LUID luid; TOKEN_PRIVILEGES tp;
    LPCWSTR privs[] = { SE_BACKUP_NAME, SE_RESTORE_NAME };
    for (int i = 0; i < 2; i++) {
        OpenProcessToken((HANDLE)-1, TOKEN_ADJUST_PRIVILEGES, &hToken);
        LookupPrivilegeValueW(NULL, privs[i], &luid);
        tp.PrivilegeCount = 1;
        tp.Privileges[0].Luid = luid;
        tp.Privileges[0].Attributes = SE_PRIVILEGE_ENABLED;
        AdjustTokenPrivileges(hToken, FALSE, &tp, sizeof(tp), NULL, NULL);
        CloseHandle(hToken);
    }
    return TRUE;
}

BOOL ReadBackupFile(HANDLE hFile, LARGE_INTEGER liSize, PVOID* ppBuf, PULONGLONG pullRead) {
    PBYTE pTmp = (PBYTE)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, (SIZE_T)(liSize.QuadPart + BACKUP_CHUNK));
    DWORD dwTotal = 0, dwRead = 0; LPVOID lpCtx = NULL;

    while (TRUE) {
        BackupRead(hFile, pTmp + dwTotal, BACKUP_CHUNK, &dwRead, FALSE, TRUE, &lpCtx);
        if (dwRead == 0) break;
        dwTotal += dwRead;
    }
    PBYTE ptr = pTmp;
    while (ptr < pTmp + dwTotal) {
        WIN32_STREAM_ID* pS = (WIN32_STREAM_ID*)ptr;
        PBYTE pData = ptr + STREAM_HDR_FIXED + pS->dwStreamNameSize;
        if (pS->dwStreamId == BACKUP_DATA && pS->dwStreamNameSize == 0) {
            *ppBuf = HeapAlloc(GetProcessHeap(), 0, (SIZE_T)pS->Size.QuadPart);
            memcpy(*ppBuf, pData, (SIZE_T)pS->Size.QuadPart);
            if (pullRead) *pullRead = pS->Size.QuadPart;
            break;
        }
        ptr = pData + (SIZE_T)pS->Size.QuadPart;
    }
    BackupRead(hFile, NULL, 0, &dwRead, TRUE, FALSE, &lpCtx);
    HeapFree(GetProcessHeap(), 0, pTmp);
    return (*ppBuf != NULL);
}

BOOL WriteBackupFile(HANDLE hFile, PVOID pData, LARGE_INTEGER liSize) {
    LPVOID lpCtx = NULL; DWORD dwWritten = 0;
    WIN32_STREAM_ID hdr = { .dwStreamId = BACKUP_DATA, .Size = liSize, .dwStreamNameSize = 0 };
    BackupWrite(hFile, (PBYTE)&hdr, STREAM_HDR_FIXED, &dwWritten, FALSE, TRUE, &lpCtx);

    ULONGLONG remain = liSize.QuadPart; PBYTE p = (PBYTE)pData;
    while (remain > 0) {
        DWORD dw = (DWORD)min(BACKUP_CHUNK, remain);
        BackupWrite(hFile, p, dw, &dwWritten, FALSE, TRUE, &lpCtx);
        p += dwWritten; remain -= dwWritten;
    }
    BackupWrite(hFile, NULL, 0, &dwWritten, TRUE, FALSE, &lpCtx);
    return TRUE;
}
```
