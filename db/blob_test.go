package db

import (
	"bytes"
	"encoding/base64"
	"sync"
	"testing"
)

// resetMasterKeyCache 重置 loadMasterKey 的 sync.Once 缓存，让测试能反复切换主密钥。
func resetMasterKeyCache() {
	masterOnce = sync.Once{}
	masterKey = nil
	masterErr = nil
}

// TestEncryptDecryptBlob 验证二进制静态加密（流量归档等）：
//   - 有主密钥：可往返、密文不含明文、篡改后解密失败；
//   - 无主密钥：fail-closed（返回 error）。
func TestEncryptDecryptBlob(t *testing.T) {
	plain := []byte("Cookie: session=secret-value; <html>body</html>")

	resetMasterKeyCache()
	t.Setenv("ARTEX_KEY_SECRET", base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x5a}, 32)))
	t.Cleanup(resetMasterKeyCache)

	sealed, err := EncryptBlob(plain)
	if err != nil {
		t.Fatalf("EncryptBlob: %v", err)
	}
	if bytes.Contains(sealed, []byte("session=secret-value")) {
		t.Fatal("密文中出现明文")
	}
	got, err := DecryptBlob(sealed)
	if err != nil {
		t.Fatalf("DecryptBlob: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("往返不一致: got %q", got)
	}

	// 篡改任意一个字节后必须解密失败（GCM 认证）。
	bad := append([]byte(nil), sealed...)
	bad[len(bad)-1] ^= 0xff
	if _, err := DecryptBlob(bad); err == nil {
		t.Error("篡改后的密文应当解密失败")
	}

	// 无主密钥 → fail-closed。
	resetMasterKeyCache()
	t.Setenv("ARTEX_KEY_SECRET", "")
	t.Setenv("ARTEX_KEY_SECRET_FILE", "")
	if _, err := EncryptBlob(plain); err == nil {
		t.Error("无主密钥时 EncryptBlob 应返回错误（fail-closed）")
	}
}
