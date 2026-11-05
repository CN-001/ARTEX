# 03 — Thread Pool & Synchronization

## Private Thread Pool with Cleanup Group
```c
#include <windows.h>

PTP_POOL            g_Pool         = NULL;
TP_CALLBACK_ENVIRON g_CallBackEnv  = { 0 };
PTP_CLEANUP_GROUP   g_CleanupGroup = NULL;

typedef struct _WORK_ITEM {
    LPTHREAD_START_ROUTINE pRoutine;
    PVOID                  pContext;
} WORK_ITEM, *PWORK_ITEM;

VOID CALLBACK ThreadpoolCallback(PTP_CALLBACK_INSTANCE Inst, PVOID lpParam, PTP_WORK Work) {
    UNREFERENCED_PARAMETER(Inst);
    PWORK_ITEM pItem = (PWORK_ITEM)lpParam;
    pItem->pRoutine(pItem->pContext);
    CloseThreadpoolWork(Work);
    LocalFree(pItem);
}

BOOL InitializeThreadPool(VOID) {
    g_Pool = CreateThreadpool(NULL);
    if (!g_Pool) return FALSE;

    DWORD dwCpu = GetActiveProcessorCount(ALL_PROCESSOR_GROUPS);
    SetThreadpoolThreadMinimum(g_Pool, dwCpu);
    SetThreadpoolThreadMaximum(g_Pool, dwCpu * 2);

    InitializeThreadpoolEnvironment(&g_CallBackEnv);
    SetThreadpoolCallbackPool(&g_CallBackEnv, g_Pool);
    SetThreadpoolCallbackRunsLong(&g_CallBackEnv);

    g_CleanupGroup = CreateThreadpoolCleanupGroup();
    if (!g_CleanupGroup) return FALSE;
    SetThreadpoolCallbackCleanupGroup(&g_CallBackEnv, g_CleanupGroup, NULL);
    return TRUE;
}

BOOL EnqueueWork(IN LPTHREAD_START_ROUTINE pRoutine, IN PVOID pContext) {
    PWORK_ITEM pItem = (PWORK_ITEM)LocalAlloc(LPTR, sizeof(WORK_ITEM));
    if (!pItem) return FALSE;
    pItem->pRoutine = pRoutine;
    pItem->pContext  = pContext;
    PTP_WORK pWork = CreateThreadpoolWork(ThreadpoolCallback, pItem, &g_CallBackEnv);
    if (!pWork) { LocalFree(pItem); return FALSE; }
    SubmitThreadpoolWork(pWork);
    return TRUE;
}

VOID DestroyThreadPool(VOID) {
    CloseThreadpoolCleanupGroupMembers(g_CleanupGroup, FALSE, NULL);
    CloseThreadpoolCleanupGroup(g_CleanupGroup);
    DestroyThreadpoolEnvironment(&g_CallBackEnv);
    CloseThreadpool(g_Pool);
}
```

## Active Task Counter Pattern
```c
volatile LONG64 g_ActiveTasks = 0;

InterlockedIncrement64(&g_ActiveTasks);
InterlockedDecrement64(&g_ActiveTasks);
while (InterlockedCompareExchange64(&g_ActiveTasks, 0, 0) != 0) Sleep(250);
```

## Synchronization Primitives

### Interlocked Operations (atomic, no blocking)
```c
volatile LONG g_Counter = 0;
InterlockedIncrement(&g_Counter);
InterlockedDecrement(&g_Counter);
InterlockedAdd(&g_Counter, 5);
InterlockedExchange(&g_Counter, 100);
InterlockedCompareExchange(&g_Counter, 10, 5);
```

### Critical Section (mutual exclusion)
```c
CRITICAL_SECTION g_CS = { 0 };
InitializeCriticalSection(&g_CS);
EnterCriticalSection(&g_CS);
LeaveCriticalSection(&g_CS);
DeleteCriticalSection(&g_CS);
```

### SRW Lock (concurrent readers OR exclusive writer)
```c
SRWLOCK g_Lock = { 0 };
InitializeSRWLock(&g_Lock);

AcquireSRWLockExclusive(&g_Lock);
g_SharedData = newValue;
ReleaseSRWLockExclusive(&g_Lock);

AcquireSRWLockShared(&g_Lock);
printf("%d", g_SharedData);
ReleaseSRWLockShared(&g_Lock);
```

### WaitOnAddress (Windows 8+, efficient wait on memory change)
```c
#pragma comment(lib, "Synchronization.lib")
volatile LONG g_Signal = 0;

LONG compare = 0;
WaitOnAddress(&g_Signal, &compare, sizeof(LONG), INFINITE);

InterlockedExchange(&g_Signal, 1);
WakeByAddressSingle(&g_Signal);
```

## Custom Heap Management
```c
HANDLE g_Heap = NULL;
#define INITIAL_HEAP_SIZE 0x1000

BOOL CreateMemHeap(VOID) {
    g_Heap = HeapCreate(HEAP_GENERATE_EXCEPTIONS, INITIAL_HEAP_SIZE, 0);
    if (!g_Heap) return FALSE;
    HeapSetInformation(g_Heap, HeapEnableTerminationOnCorruption, NULL, 0);
    return TRUE;
}
PVOID AllocMem(DWORDLONG dwLen) { return HeapAlloc(g_Heap, HEAP_GENERATE_EXCEPTIONS | HEAP_ZERO_MEMORY, dwLen); }
PVOID ReAllocMem(PVOID p, DWORDLONG dwNew) { return HeapReAlloc(g_Heap, HEAP_GENERATE_EXCEPTIONS | HEAP_ZERO_MEMORY, p, dwNew); }
VOID FreeMem(PVOID* pp) { if (*pp) { HeapFree(g_Heap, 0, *pp); *pp = NULL; } }
BOOL DestroyMemHeap(VOID) { return HeapDestroy(g_Heap); }

LPWSTR HeapStrDupW(LPCWSTR szStr) {
    DWORD n = lstrlenW(szStr) + 1;
    LPWSTR p = (LPWSTR)AllocMem(n * sizeof(WCHAR));
    StringCchCopyW(p, n, szStr);
    return p;
}
```
