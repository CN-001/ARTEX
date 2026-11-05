# 01 — Symmetric Encryption

## Shared Constants & Macros
```c
#include <windows.h>
#include <bcrypt.h>
#pragma comment(lib, "bcrypt.lib")

#define AES_KEY_SIZE    32
#define AES_IV_SIZE     16
#define AES_BLOCK_SIZE  16
#define NT_SUCCESS(S)   (((NTSTATUS)(S)) >= 0x00)
#define SAFE_FREE(p)    do { if (p) { LocalFree(p); (p) = NULL; } } while(0)
```

## DuplicateData — CBC mode consumes IV, must duplicate before each call
```c
PBYTE DuplicateData(IN PBYTE pData, IN DWORD dwLen) {
    if (!pData || dwLen == 0) return NULL;
    PBYTE p = (PBYTE)LocalAlloc(LPTR, dwLen);
    if (p) memcpy(p, pData, dwLen);
    return p;
}
```

## AES Internal Setup Helper
```c
static BOOL AesSetup(LPCWSTR szMode, PBYTE pKey,
    BCRYPT_ALG_HANDLE* phAlg, BCRYPT_KEY_HANDLE* phKey,
    PBYTE* ppKeyObj, PDWORD pdwKeyObjSize) {
    DWORD dwResult = 0, dwBlockSize = 0;
    if (!NT_SUCCESS(BCryptOpenAlgorithmProvider(phAlg, BCRYPT_AES_ALGORITHM, NULL, 0))) return FALSE;
    BCryptGetProperty(*phAlg, BCRYPT_OBJECT_LENGTH, (PBYTE)pdwKeyObjSize, sizeof(DWORD), &dwResult, 0);
    BCryptGetProperty(*phAlg, BCRYPT_BLOCK_LENGTH, (PBYTE)&dwBlockSize, sizeof(DWORD), &dwResult, 0);
    if (dwBlockSize != AES_BLOCK_SIZE) return FALSE;
    if (!NT_SUCCESS(BCryptSetProperty(*phAlg, BCRYPT_CHAINING_MODE, (PBYTE)szMode,
        (DWORD)(wcslen(szMode) + 1) * sizeof(WCHAR), 0))) return FALSE;
    *ppKeyObj = (PBYTE)LocalAlloc(LPTR, *pdwKeyObjSize);
    if (!*ppKeyObj) return FALSE;
    return NT_SUCCESS(BCryptGenerateSymmetricKey(*phAlg, phKey, *ppKeyObj, *pdwKeyObjSize, pKey, AES_KEY_SIZE, 0));
}
```

---

## AES-256-CBC Encrypt
```c
BOOL Aes256CbcEncrypt(
    IN  PBYTE pPlainText, IN DWORD dwPlainSize,
    IN  PBYTE p32Key, IN PBYTE p16Iv,
    OUT PBYTE* ppCipherText, OUT PDWORD pdwCipherSize)
{
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_KEY_HANDLE hKey = NULL;
    PBYTE pKeyObj = NULL, pIvDup = NULL, pCipher = NULL;
    DWORD dwKeyObjSize = 0, dwCipherSize = 0;
    BOOL bResult = FALSE;

    if (!AesSetup(BCRYPT_CHAIN_MODE_CBC, p32Key, &hAlg, &hKey, &pKeyObj, &dwKeyObjSize)) goto _END;

    DWORD dwPadFlag = ((dwPlainSize % AES_BLOCK_SIZE) == 0) ? 0 : BCRYPT_BLOCK_PADDING;

    pIvDup = DuplicateData(p16Iv, AES_IV_SIZE);
    if (!pIvDup) goto _END;
    if (!NT_SUCCESS(BCryptEncrypt(hKey, pPlainText, dwPlainSize, NULL, pIvDup, AES_IV_SIZE,
        NULL, 0, &dwCipherSize, dwPadFlag))) goto _END;

    pCipher = (PBYTE)LocalAlloc(LPTR, dwCipherSize);
    SAFE_FREE(pIvDup);
    pIvDup = DuplicateData(p16Iv, AES_IV_SIZE);
    if (!pIvDup) goto _END;
    if (!NT_SUCCESS(BCryptEncrypt(hKey, pPlainText, dwPlainSize, NULL, pIvDup, AES_IV_SIZE,
        pCipher, dwCipherSize, &dwCipherSize, dwPadFlag))) goto _END;

    *ppCipherText = pCipher; *pdwCipherSize = dwCipherSize;
    pCipher = NULL; bResult = TRUE;

_END:
    SAFE_FREE(pIvDup); SAFE_FREE(pCipher); SAFE_FREE(pKeyObj);
    if (hKey) BCryptDestroyKey(hKey);
    if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bResult;
}
```

## AES-256-CBC Decrypt
```c
BOOL Aes256CbcDecrypt(
    IN  PBYTE pCipherText, IN DWORD dwCipherSize,
    IN  PBYTE p32Key, IN PBYTE p16Iv,
    OUT PBYTE* ppPlainText, OUT PDWORD pdwPlainSize)
{
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_KEY_HANDLE hKey = NULL;
    PBYTE pKeyObj = NULL, pIvDup = NULL, pPlain = NULL;
    DWORD dwKeyObjSize = 0, dwPlainSize = 0;
    BOOL bResult = FALSE;

    if (!AesSetup(BCRYPT_CHAIN_MODE_CBC, p32Key, &hAlg, &hKey, &pKeyObj, &dwKeyObjSize)) goto _END;
    DWORD dwPadFlag = ((dwCipherSize % AES_BLOCK_SIZE) == 0) ? BCRYPT_BLOCK_PADDING : 0;

    pIvDup = DuplicateData(p16Iv, AES_IV_SIZE);
    if (!pIvDup) goto _END;
    if (!NT_SUCCESS(BCryptDecrypt(hKey, pCipherText, dwCipherSize, NULL, pIvDup, AES_IV_SIZE,
        NULL, 0, &dwPlainSize, dwPadFlag))) goto _END;
    pPlain = (PBYTE)LocalAlloc(LPTR, dwPlainSize);
    SAFE_FREE(pIvDup);
    pIvDup = DuplicateData(p16Iv, AES_IV_SIZE);
    if (!pIvDup) goto _END;
    if (!NT_SUCCESS(BCryptDecrypt(hKey, pCipherText, dwCipherSize, NULL, pIvDup, AES_IV_SIZE,
        pPlain, dwPlainSize, &dwPlainSize, dwPadFlag))) goto _END;

    *ppPlainText = pPlain; *pdwPlainSize = dwPlainSize;
    pPlain = NULL; bResult = TRUE;
_END:
    SAFE_FREE(pIvDup); SAFE_FREE(pPlain); SAFE_FREE(pKeyObj);
    if (hKey) BCryptDestroyKey(hKey); if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bResult;
}
```

---

## AES-256-CFB Encrypt (stream-like, output == input size, no padding)
```c
BOOL Aes256CfbEncrypt(
    IN PBYTE pPlain, IN DWORD dwSize, IN PBYTE p32Key, IN PBYTE p16Iv,
    OUT PBYTE* ppOut, OUT PDWORD pdwOutSize)
{
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_KEY_HANDLE hKey = NULL;
    PBYTE pKeyObj = NULL, pIvDup = NULL, pCipher = NULL;
    DWORD dwKeyObjSize = 0, dwOutSize = 0; BOOL bResult = FALSE;

    if (!AesSetup(BCRYPT_CHAIN_MODE_CFB, p32Key, &hAlg, &hKey, &pKeyObj, &dwKeyObjSize)) goto _END;
    pIvDup = DuplicateData(p16Iv, AES_IV_SIZE); if (!pIvDup) goto _END;
    if (!NT_SUCCESS(BCryptEncrypt(hKey, pPlain, dwSize, NULL, pIvDup, AES_IV_SIZE, NULL, 0, &dwOutSize, 0))) goto _END;
    pCipher = (PBYTE)LocalAlloc(LPTR, dwOutSize);
    SAFE_FREE(pIvDup); pIvDup = DuplicateData(p16Iv, AES_IV_SIZE); if (!pIvDup) goto _END;
    if (!NT_SUCCESS(BCryptEncrypt(hKey, pPlain, dwSize, NULL, pIvDup, AES_IV_SIZE, pCipher, dwOutSize, &dwOutSize, 0))) goto _END;
    *ppOut = pCipher; *pdwOutSize = dwOutSize; pCipher = NULL; bResult = TRUE;
_END:
    SAFE_FREE(pIvDup); SAFE_FREE(pCipher); SAFE_FREE(pKeyObj);
    if (hKey) BCryptDestroyKey(hKey); if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bResult;
}
```

---

## AES-256-GCM Encrypt (AEAD: confidentiality + integrity)

```c
#define AES_GCM_IV_SIZE   12
#define AES_GCM_TAG_SIZE  16

BOOL Aes256GcmEncrypt(
    IN  PBYTE pPlain, IN DWORD dwSize, IN PBYTE p32Key, IN PBYTE p12Nonce,
    OUT PBYTE* ppCipher, OUT PDWORD pdwCipherSize, OUT BYTE pTag[16])
{
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_KEY_HANDLE hKey = NULL;
    PBYTE pKeyObj = NULL, pCipher = NULL;
    DWORD dwResult = 0, dwKeyObjSize = 0, dwCipherSize = 0;
    BOOL bResult = FALSE;

    if (!NT_SUCCESS(BCryptOpenAlgorithmProvider(&hAlg, BCRYPT_AES_ALGORITHM, NULL, 0))) goto _END;
    if (!NT_SUCCESS(BCryptSetProperty(hAlg, BCRYPT_CHAINING_MODE, (PBYTE)BCRYPT_CHAIN_MODE_GCM,
        sizeof(BCRYPT_CHAIN_MODE_GCM), 0))) goto _END;
    BCryptGetProperty(hAlg, BCRYPT_OBJECT_LENGTH, (PBYTE)&dwKeyObjSize, sizeof(DWORD), &dwResult, 0);
    pKeyObj = (PBYTE)LocalAlloc(LPTR, dwKeyObjSize); if (!pKeyObj) goto _END;
    if (!NT_SUCCESS(BCryptGenerateSymmetricKey(hAlg, &hKey, pKeyObj, dwKeyObjSize, p32Key, AES_KEY_SIZE, 0))) goto _END;

    BCRYPT_AUTHENTICATED_CIPHER_MODE_INFO AuthInfo = { 0 };
    BCRYPT_INIT_AUTH_MODE_INFO(AuthInfo);
    AuthInfo.pbNonce = p12Nonce;  AuthInfo.cbNonce = AES_GCM_IV_SIZE;
    AuthInfo.pbTag   = pTag;      AuthInfo.cbTag   = AES_GCM_TAG_SIZE;
    AuthInfo.pbAuthData = NULL;   AuthInfo.cbAuthData = 0;

    pCipher = (PBYTE)LocalAlloc(LPTR, dwSize); if (!pCipher) goto _END;
    if (!NT_SUCCESS(BCryptEncrypt(hKey, pPlain, dwSize, &AuthInfo, NULL, 0,
        pCipher, dwSize, &dwCipherSize, 0))) goto _END;

    *ppCipher = pCipher; *pdwCipherSize = dwCipherSize; pCipher = NULL; bResult = TRUE;
_END:
    SAFE_FREE(pCipher); SAFE_FREE(pKeyObj);
    if (hKey) BCryptDestroyKey(hKey); if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bResult;
}
```

## AES-256-CCM (same as GCM but IV = 11 bytes)
```c
#define AES_CCM_IV_SIZE  11
```

---

## RC4 — Custom Implementation (KSA + PRGA)
```c
#define SBOX_SIZE 256

BOOL Rc4Encrypt(IN PBYTE pData, IN DWORD dwDataSize, IN PBYTE pKey, IN DWORD dwKeySize,
                OUT PBYTE* ppOut, OUT PDWORD pdwOutSize) {
    PBYTE pOut = (PBYTE)LocalAlloc(LPTR, dwDataSize);
    if (!pOut) return FALSE;
    memcpy(pOut, pData, dwDataSize);

    BYTE S[SBOX_SIZE];
    for (int i = 0; i < SBOX_SIZE; i++) S[i] = (BYTE)i;
    DWORD j = 0;
    for (int i = 0; i < SBOX_SIZE; i++) {
        j = (j + S[i] + pKey[i % dwKeySize]) % SBOX_SIZE;
        BYTE t = S[i]; S[i] = S[j]; S[j] = t;
    }

    DWORD ii = 0, jj = 0;
    for (DWORD n = 0; n < dwDataSize; n++) {
        ii = (ii + 1) % SBOX_SIZE;
        jj = (jj + S[ii]) % SBOX_SIZE;
        BYTE t = S[ii]; S[ii] = S[jj]; S[jj] = t;
        pOut[n] ^= S[(S[ii] + S[jj]) % SBOX_SIZE];
    }
    *ppOut = pOut; *pdwOutSize = dwDataSize;
    return TRUE;
}
```

## RC4 — Via SystemFunction032 (undocumented Advapi32)
```c
typedef struct { DWORD Length; DWORD MaximumLength; PVOID Buffer; } USTRING;
typedef NTSTATUS(NTAPI* fnSystemFunction032)(USTRING*, USTRING*);

BOOL Rc4SysFunc032(IN OUT PBYTE pData, IN DWORD dwDataSize, IN PBYTE pKey, IN DWORD dwKeySize) {
    HMODULE hLib = LoadLibraryA("Advapi32");
    if (!hLib) return FALSE;
    fnSystemFunction032 pFunc = (fnSystemFunction032)GetProcAddress(hLib, "SystemFunction032");
    if (!pFunc) { FreeLibrary(hLib); return FALSE; }
    USTRING Data = { dwDataSize, dwDataSize, pData };
    USTRING Key  = { dwKeySize,  dwKeySize,  pKey };
    NTSTATUS s = pFunc(&Data, &Key);
    FreeLibrary(hLib);
    return (s == 0);
}
```

## RC4 — Via BCrypt
```c
BOOL Rc4BcryptEncrypt(IN OUT PBYTE pData, IN DWORD dwSize, IN PBYTE pKey, IN DWORD dwKeySize) {
    BCRYPT_ALG_HANDLE hAlg = NULL; BCRYPT_KEY_HANDLE hKey = NULL;
    DWORD dwResult = 0; BOOL bOk = FALSE;
    if (!NT_SUCCESS(BCryptOpenAlgorithmProvider(&hAlg, BCRYPT_RC4_ALGORITHM, NULL, 0))) goto _END;
    if (!NT_SUCCESS(BCryptGenerateSymmetricKey(hAlg, &hKey, NULL, 0, pKey, dwKeySize, 0))) goto _END;
    if (!NT_SUCCESS(BCryptEncrypt(hKey, pData, dwSize, NULL, NULL, 0, pData, dwSize, &dwResult, 0))) goto _END;
    bOk = TRUE;
_END:
    if (hKey) BCryptDestroyKey(hKey); if (hAlg) BCryptCloseAlgorithmProvider(hAlg, 0);
    return bOk;
}
```

---

## ChaCha20 — Complete Custom Implementation (no Windows API)
```c
#define CHACHA_KEY_LEN    32
#define CHACHA_NONCE_LEN  12
#define CHACHA_ROUNDS     20
#define CHACHA_BLOCK      64
#define CHACHA_STATE_DW   16
#define CHACHA_CONST      "expand 32-byte k"

#define ROTL32(v, n) (((v) << (n)) | ((v) >> (32 - (n))))
#define QR(a,b,c,d) \
    a += b; d ^= a; d = ROTL32(d,16); \
    c += d; b ^= c; b = ROTL32(b,12); \
    a += b; d ^= a; d = ROTL32(d, 8); \
    c += d; b ^= c; b = ROTL32(b, 7);

typedef struct {
    DWORD state[CHACHA_STATE_DW];
    BYTE  leftover[CHACHA_BLOCK];
    DWORD leftoverUsed;
} CHACHA20_CTX;

static void ChaCha20Block(DWORD out[16], const DWORD in[16]) {
    DWORD x[16]; memcpy(x, in, 64);
    for (int i = 0; i < CHACHA_ROUNDS; i += 2) {
        QR(x[0],x[4],x[ 8],x[12]); QR(x[1],x[5],x[ 9],x[13]);
        QR(x[2],x[6],x[10],x[14]); QR(x[3],x[7],x[11],x[15]);
        QR(x[0],x[5],x[10],x[15]); QR(x[1],x[6],x[11],x[12]);
        QR(x[2],x[7],x[ 8],x[13]); QR(x[3],x[4],x[ 9],x[14]);
    }
    for (int i = 0; i < 16; i++) out[i] = x[i] + in[i];
}

void ChaCha20Init(CHACHA20_CTX* ctx, const BYTE key[32], const BYTE nonce[12], DWORD counter) {
    const DWORD* c = (const DWORD*)CHACHA_CONST;
    ctx->state[0]=c[0]; ctx->state[1]=c[1]; ctx->state[2]=c[2]; ctx->state[3]=c[3];
    memcpy(&ctx->state[4], key, 32);
    ctx->state[12] = counter;
    memcpy(&ctx->state[13], nonce, 12);
    ctx->leftoverUsed = CHACHA_BLOCK;
}

void ChaCha20Crypt(CHACHA20_CTX* ctx, BYTE* data, DWORD len) {
    for (DWORD i = 0; i < len; i++) {
        if (ctx->leftoverUsed >= CHACHA_BLOCK) {
            DWORD block[16];
            ChaCha20Block(block, ctx->state);
            memcpy(ctx->leftover, block, CHACHA_BLOCK);
            ctx->state[12]++;
            ctx->leftoverUsed = 0;
        }
        data[i] ^= ctx->leftover[ctx->leftoverUsed++];
    }
}
```
