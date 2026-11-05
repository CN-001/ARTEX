# 02 — Asymmetric Encryption & Key Derivation

## RSA Key Pair Generation (BCrypt)
```c
BOOL GenerateRsaKeyPair(IN DWORD dwKeyBits,
    OUT PBYTE* ppPubKey, OUT PDWORD pdwPubKeySize,
    OUT PBYTE* ppPrivKey, OUT PDWORD pdwPrivKeySize)
{
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_KEY_HANDLE hKey = NULL;
    BOOL bResult = FALSE;

    if (!NT_SUCCESS(BCryptOpenAlgorithmProvider(&hAlg, BCRYPT_RSA_ALGORITHM, NULL, 0))) goto _END;
    if (!NT_SUCCESS(BCryptGenerateKeyPair(hAlg, &hKey, dwKeyBits, 0))) goto _END;
    if (!NT_SUCCESS(BCryptFinalizeKeyPair(hKey, 0))) goto _END;

    BCryptExportKey(hKey, NULL, BCRYPT_RSAPUBLIC_BLOB, NULL, 0, pdwPubKeySize, 0);
    *ppPubKey = (PBYTE)LocalAlloc(LPTR, *pdwPubKeySize);
    BCryptExportKey(hKey, NULL, BCRYPT_RSAPUBLIC_BLOB, *ppPubKey, *pdwPubKeySize, pdwPubKeySize, 0);

    BCryptExportKey(hKey, NULL, BCRYPT_RSAFULLPRIVATE_BLOB, NULL, 0, pdwPrivKeySize, 0);
    *ppPrivKey = (PBYTE)LocalAlloc(LPTR, *pdwPrivKeySize);
    BCryptExportKey(hKey, NULL, BCRYPT_RSAFULLPRIVATE_BLOB, *ppPrivKey, *pdwPrivKeySize, pdwPrivKeySize, 0);

    bResult = TRUE;
_END:
    if (hKey) BCryptDestroyKey(hKey);
    if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bResult;
}
```

## RSA Encrypt with Public Key
```c
BOOL RsaEncrypt(IN PBYTE pPubKeyBlob, IN DWORD dwPubKeySize,
    IN PBYTE pPlain, IN DWORD dwPlainSize,
    OUT PBYTE* ppCipher, OUT PDWORD pdwCipherSize)
{
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_KEY_HANDLE hKey = NULL;
    BOOL bResult = FALSE; PBYTE pCipher = NULL; DWORD dwOut = 0;

    if (!NT_SUCCESS(BCryptOpenAlgorithmProvider(&hAlg, BCRYPT_RSA_ALGORITHM, NULL, 0))) goto _END;
    if (!NT_SUCCESS(BCryptImportKeyPair(hAlg, NULL, BCRYPT_RSAPUBLIC_BLOB, &hKey, pPubKeyBlob, dwPubKeySize, 0))) goto _END;

    BCryptEncrypt(hKey, pPlain, dwPlainSize, NULL, NULL, 0, NULL, 0, &dwOut, BCRYPT_PAD_PKCS1);
    pCipher = (PBYTE)LocalAlloc(LPTR, dwOut);
    if (!NT_SUCCESS(BCryptEncrypt(hKey, pPlain, dwPlainSize, NULL, NULL, 0, pCipher, dwOut, &dwOut, BCRYPT_PAD_PKCS1))) goto _END;

    *ppCipher = pCipher; *pdwCipherSize = dwOut; pCipher = NULL; bResult = TRUE;
_END:
    SAFE_FREE(pCipher); if (hKey) BCryptDestroyKey(hKey); if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bResult;
}
```

## RSA Decrypt with Private Key
```c
BOOL RsaDecrypt(IN PBYTE pPrivKeyBlob, IN DWORD dwPrivKeySize,
    IN PBYTE pCipher, IN DWORD dwCipherSize,
    OUT PBYTE* ppPlain, OUT PDWORD pdwPlainSize)
{
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_KEY_HANDLE hKey = NULL;
    BOOL bResult = FALSE; PBYTE pPlain = NULL; DWORD dwOut = 0;

    if (!NT_SUCCESS(BCryptOpenAlgorithmProvider(&hAlg, BCRYPT_RSA_ALGORITHM, NULL, 0))) goto _END;
    if (!NT_SUCCESS(BCryptImportKeyPair(hAlg, NULL, BCRYPT_RSAFULLPRIVATE_BLOB, &hKey, pPrivKeyBlob, dwPrivKeySize, 0))) goto _END;

    BCryptDecrypt(hKey, pCipher, dwCipherSize, NULL, NULL, 0, NULL, 0, &dwOut, BCRYPT_PAD_PKCS1);
    pPlain = (PBYTE)LocalAlloc(LPTR, dwOut);
    if (!NT_SUCCESS(BCryptDecrypt(hKey, pCipher, dwCipherSize, NULL, NULL, 0, pPlain, dwOut, &dwOut, BCRYPT_PAD_PKCS1))) goto _END;

    *ppPlain = pPlain; *pdwPlainSize = dwOut; pPlain = NULL; bResult = TRUE;
_END:
    SAFE_FREE(pPlain); if (hKey) BCryptDestroyKey(hKey); if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bResult;
}
```

## RSA via NCrypt (alternative, hardware-backed key storage)
```c
#include <ncrypt.h>
#pragma comment(lib, "ncrypt.lib")

BOOL NcryptGenerateRsaKeyPair(IN DWORD dwKeyBits,
    OUT PBYTE* ppPubKey, OUT PDWORD pdwPubSize,
    OUT PBYTE* ppPrivKey, OUT PDWORD pdwPrivSize)
{
    NCRYPT_PROV_HANDLE hProv = 0; NCRYPT_KEY_HANDLE hKey = 0;
    DWORD dwExportPolicy = NCRYPT_ALLOW_PLAINTEXT_EXPORT_FLAG;
    BOOL bResult = FALSE;

    if (NCryptOpenStorageProvider(&hProv, MS_KEY_STORAGE_PROVIDER, 0) != 0) goto _END;
    if (NCryptCreatePersistedKey(hProv, &hKey, BCRYPT_RSA_ALGORITHM, NULL, 0, 0) != 0) goto _END;
    NCryptSetProperty(hKey, NCRYPT_LENGTH_PROPERTY, (PBYTE)&dwKeyBits, sizeof(DWORD), 0);
    NCryptSetProperty(hKey, NCRYPT_EXPORT_POLICY_PROPERTY, (PBYTE)&dwExportPolicy, sizeof(DWORD), 0);
    if (NCryptFinalizeKey(hKey, 0) != 0) goto _END;

    NCryptExportKey(hKey, 0, BCRYPT_RSAPUBLIC_BLOB, NULL, NULL, 0, pdwPubSize, 0);
    *ppPubKey = (PBYTE)LocalAlloc(LPTR, *pdwPubSize);
    NCryptExportKey(hKey, 0, BCRYPT_RSAPUBLIC_BLOB, NULL, *ppPubKey, *pdwPubSize, pdwPubSize, 0);

    NCryptExportKey(hKey, 0, BCRYPT_RSAFULLPRIVATE_BLOB, NULL, NULL, 0, pdwPrivSize, 0);
    *ppPrivKey = (PBYTE)LocalAlloc(LPTR, *pdwPrivSize);
    NCryptExportKey(hKey, 0, BCRYPT_RSAFULLPRIVATE_BLOB, NULL, *ppPrivKey, *pdwPrivSize, pdwPrivSize, 0);

    bResult = TRUE;
_END:
    if (hKey) NCryptFreeObject(hKey); if (hProv) NCryptFreeObject(hProv);
    return bResult;
}
```

---

## 7 Random Key Generation Methods

### 1. BCryptGenRandom (CSPRNG, recommended)
```c
PBYTE GenKey_BCrypt(DWORD dwLen) {
    PBYTE p = (PBYTE)LocalAlloc(LPTR, dwLen);
    if (!p) return NULL;
    if (BCryptGenRandom(NULL, p, dwLen, BCRYPT_USE_SYSTEM_PREFERRED_RNG) != 0) { LocalFree(p); return NULL; }
    return p;
}
```

### 2. CryptoAPI (CryptGenRandom)
```c
#include <wincrypt.h>
PBYTE GenKey_CryptoAPI(DWORD dwLen) {
    PBYTE p = (PBYTE)LocalAlloc(LPTR, dwLen); HCRYPTPROV hProv = 0;
    if (!p) return NULL;
    if (!CryptAcquireContextA(&hProv, NULL, NULL, PROV_RSA_AES, CRYPT_VERIFYCONTEXT)) { LocalFree(p); return NULL; }
    BOOL ok = CryptGenRandom(hProv, dwLen, p);
    CryptReleaseContext(hProv, 0);
    if (!ok) { LocalFree(p); return NULL; }
    return p;
}
```

### 3. RtlGenRandom (SystemFunction036)
```c
typedef BOOLEAN(APIENTRY* fnRtlGenRandom)(PVOID, ULONG);
PBYTE GenKey_RtlGenRandom(DWORD dwLen) {
    PBYTE p = (PBYTE)LocalAlloc(LPTR, dwLen); if (!p) return NULL;
    HMODULE h = LoadLibraryA("Advapi32");
    fnRtlGenRandom fn = (fnRtlGenRandom)GetProcAddress(h, "SystemFunction036");
    if (!fn || !fn(p, dwLen)) { LocalFree(p); FreeLibrary(h); return NULL; }
    FreeLibrary(h); return p;
}
```

### 4. RDRAND (hardware CPU entropy)
```c
#include <intrin.h>
PBYTE GenKey_RDRAND(DWORD dwLen) {
    PBYTE p = (PBYTE)LocalAlloc(LPTR, dwLen); if (!p) return NULL;
    unsigned int v = 0;
    for (DWORD i = 0; i < dwLen; i++) {
        if ((i % 4) == 0) { int r = 10; while (r-- > 0 && !_rdrand32_step(&v)); }
        p[i] = (BYTE)(v >> ((i % 4) * 8));
    }
    return p;
}
```

### 5. GetTickCount64 + rand (weak, educational)
```c
PBYTE GenKey_TickCount(DWORD dwLen) {
    PBYTE p = (PBYTE)LocalAlloc(LPTR, dwLen); if (!p) return NULL;
    srand((unsigned)(GetTickCount64() ^ time(NULL)));
    USHORT usAddr = (USHORT)((ULONG_PTR)p & 0xFFFF);
    for (DWORD i = 0; i < dwLen; i++) {
        p[i] = (BYTE)((GetTickCount64() * rand() + usAddr) % 0xFF);
        srand((unsigned)(GetTickCount64() + rand()));
    }
    return p;
}
```

### 6. __TIME__ compile-time macro (weak)
```c
PBYTE GenKey_CompileTime(DWORD dwLen) {
    PBYTE p = (PBYTE)LocalAlloc(LPTR, dwLen); if (!p) return NULL;
    USHORT usAddr = (USHORT)((ULONG_PTR)p & 0xFFFF);
    for (DWORD i = 0; i < dwLen; i++) {
        srand(__TIME__[i % 7] + rand());
        p[i] = (BYTE)((__TIME__[i % 7] * rand() + usAddr) % 0xFF);
    }
    return p;
}
```

### 7. rand() + time(NULL) (weakest — trivially brute-forceable)
```c
PBYTE GenKey_Rand(DWORD dwLen) {
    PBYTE p = (PBYTE)LocalAlloc(LPTR, dwLen); if (!p) return NULL;
    srand((unsigned)time(NULL));
    for (DWORD i = 0; i < dwLen; i++) p[i] = (BYTE)(rand() % 0xFF);
    return p;
}
```

---

## HMAC-SHA256
```c
#define HMAC_SHA256_LEN 32

BOOL HmacSha256(IN PBYTE pKey, IN DWORD dwKeyLen, IN PBYTE pData, IN DWORD dwDataLen,
    OUT BYTE pOutput[32])
{
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_HASH_HANDLE hHash = NULL;
    PBYTE pHashObj = NULL; DWORD dwHashObjSize = 0, dwResult = 0; BOOL bOk = FALSE;

    if (!NT_SUCCESS(BCryptOpenAlgorithmProvider(&hAlg, BCRYPT_SHA256_ALGORITHM, NULL, BCRYPT_ALG_HANDLE_HMAC_FLAG))) goto _END;
    BCryptGetProperty(hAlg, BCRYPT_OBJECT_LENGTH, (PBYTE)&dwHashObjSize, sizeof(DWORD), &dwResult, 0);
    pHashObj = (PBYTE)LocalAlloc(LPTR, dwHashObjSize); if (!pHashObj) goto _END;
    if (!NT_SUCCESS(BCryptCreateHash(hAlg, &hHash, pHashObj, dwHashObjSize, pKey, dwKeyLen, 0))) goto _END;
    if (!NT_SUCCESS(BCryptHashData(hHash, pData, dwDataLen, 0))) goto _END;
    if (!NT_SUCCESS(BCryptFinishHash(hHash, pOutput, HMAC_SHA256_LEN, 0))) goto _END;
    bOk = TRUE;
_END:
    if (hHash) BCryptDestroyHash(hHash); SAFE_FREE(pHashObj);
    if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bOk;
}
```

## HKDF (Extract + Expand)
```c
#define HKDF_HASH_LEN 32

BOOL HkdfExtract(IN PBYTE pIKM, IN DWORD dwIKMLen, IN OPTIONAL PBYTE pSalt, IN DWORD dwSaltLen,
    OUT BYTE pPRK[32])
{
    BYTE zeroSalt[HKDF_HASH_LEN] = { 0 };
    PBYTE pActualSalt = pSalt ? pSalt : zeroSalt;
    DWORD dwActualSaltLen = pSalt ? dwSaltLen : HKDF_HASH_LEN;
    return HmacSha256(pActualSalt, dwActualSaltLen, pIKM, dwIKMLen, pPRK);
}

BOOL HkdfExpand(IN BYTE pPRK[32], IN OPTIONAL PBYTE pInfo, IN DWORD dwInfoLen,
    IN DWORD dwOKMLen, OUT PBYTE* ppOKM)
{
    DWORD dwBlocks = (dwOKMLen + HKDF_HASH_LEN - 1) / HKDF_HASH_LEN;
    if (dwBlocks > 255) return FALSE;

    PBYTE pOKM = (PBYTE)LocalAlloc(LPTR, dwOKMLen);
    if (!pOKM) return FALSE;

    BYTE prev[HKDF_HASH_LEN] = { 0 }, block[HKDF_HASH_LEN];
    for (BYTE i = 1; i <= (BYTE)dwBlocks; i++) {
        DWORD dwBufLen = ((i > 1) ? HKDF_HASH_LEN : 0) + dwInfoLen + 1;
        PBYTE pBuf = (PBYTE)LocalAlloc(LPTR, dwBufLen);
        DWORD off = 0;
        if (i > 1) { memcpy(pBuf, prev, HKDF_HASH_LEN); off = HKDF_HASH_LEN; }
        if (pInfo && dwInfoLen) { memcpy(pBuf + off, pInfo, dwInfoLen); off += dwInfoLen; }
        pBuf[off] = i;

        HmacSha256(pPRK, HKDF_HASH_LEN, pBuf, dwBufLen, block);
        LocalFree(pBuf);

        DWORD dwCopy = ((DWORD)i == dwBlocks) ? (dwOKMLen - (i - 1) * HKDF_HASH_LEN) : HKDF_HASH_LEN;
        memcpy(pOKM + (i - 1) * HKDF_HASH_LEN, block, dwCopy);
        memcpy(prev, block, HKDF_HASH_LEN);
    }
    *ppOKM = pOKM;
    return TRUE;
}
```

## SHA-1 Hash (for key fingerprinting)
```c
#define SHA1_LEN 20
BOOL Sha1Hash(IN PBYTE pData, IN DWORD dwLen, OUT BYTE pHash[20]) {
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_HASH_HANDLE hHash = NULL;
    PBYTE pObj = NULL; DWORD dwObjSize = 0, dwResult = 0; BOOL bOk = FALSE;
    if (!NT_SUCCESS(BCryptOpenAlgorithmProvider(&hAlg, BCRYPT_SHA1_ALGORITHM, NULL, 0))) goto _END;
    BCryptGetProperty(hAlg, BCRYPT_OBJECT_LENGTH, (PBYTE)&dwObjSize, sizeof(DWORD), &dwResult, 0);
    pObj = (PBYTE)LocalAlloc(LPTR, dwObjSize);
    if (!NT_SUCCESS(BCryptCreateHash(hAlg, &hHash, pObj, dwObjSize, NULL, 0, 0))) goto _END;
    if (!NT_SUCCESS(BCryptHashData(hHash, pData, dwLen, 0))) goto _END;
    if (!NT_SUCCESS(BCryptFinishHash(hHash, pHash, SHA1_LEN, 0))) goto _END;
    bOk = TRUE;
_END:
    if (hHash) BCryptDestroyHash(hHash); SAFE_FREE(pObj);
    if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0); return bOk;
}
```
