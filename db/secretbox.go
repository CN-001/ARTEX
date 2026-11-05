package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// LLM API keys are encrypted at rest with AES-256-GCM so that a database dump,
// a backup archive, or an accidental read of llm_profiles does not expose a live
// credential. The plaintext key exists only in memory while a provider is built.
//
// The master key is read once, from (in order):
//
//	ARTEX_KEY_SECRET       — base64 / hex / raw 32 bytes
//	ARTEX_KEY_SECRET_FILE  — path to a file containing one of the above
//
// 存储格式：enc:v1:<base64(nonce||ciphertext)>；不带该前缀的值按明文处理（向后兼容）。
//
// Encryption is fail-closed for writes: without a master key, new API keys must
// not be persisted in plaintext. Reads remain backward compatible so existing
// plaintext profiles can still be used while operators arrange a key.
const apiKeyEncPrefix = "enc:v1:"

var (
	masterOnce sync.Once
	masterKey  []byte
	masterErr  error
)

// loadMasterKey resolves and caches the master key. Callers must treat a nil key
// as "encryption unavailable", never as an error to abort startup with.
func loadMasterKey() ([]byte, error) {
	masterOnce.Do(func() {
		raw := strings.TrimSpace(os.Getenv("ARTEX_KEY_SECRET"))
		source := "ARTEX_KEY_SECRET"
		if raw == "" {
			if path := strings.TrimSpace(os.Getenv("ARTEX_KEY_SECRET_FILE")); path != "" {
				b, err := os.ReadFile(path)
				if err != nil {
					masterErr = fmt.Errorf("读取 %s 失败: %w", path, err)
					return
				}
				raw = strings.TrimSpace(string(b))
				source = path
			}
		}
		if raw == "" {
			masterErr = errors.New("未配置 ARTEX_KEY_SECRET 或 ARTEX_KEY_SECRET_FILE")
			return
		}
		if k, err := base64.StdEncoding.DecodeString(raw); err == nil && len(k) == 32 {
			masterKey = k
			log.Printf("[secretbox] API key 加密已启用（主密钥来自 %s）", source)
			return
		}
		if k, err := base64.RawStdEncoding.DecodeString(raw); err == nil && len(k) == 32 {
			masterKey = k
			log.Printf("[secretbox] API key 加密已启用（主密钥来自 %s）", source)
			return
		}
		if k, err := hex.DecodeString(raw); err == nil && len(k) == 32 {
			masterKey = k
			log.Printf("[secretbox] API key 加密已启用（主密钥来自 %s）", source)
			return
		}
		if len(raw) == 32 {
			masterKey = []byte(raw)
			log.Printf("[secretbox] API key 加密已启用（主密钥来自 %s，原文 32 字节）", source)
			return
		}
		masterErr = fmt.Errorf("主密钥必须是 32 字节（base64/hex/原文），当前解析后 %d 字节", len(raw))
	})
	return masterKey, masterErr
}

// APIKeyEncryptionAvailable reports whether at-rest encryption is active.
func APIKeyEncryptionAvailable() bool {
	k, err := loadMasterKey()
	return err == nil && len(k) == 32
}

// BlobEncryptionAvailable 报告主密钥是否可用于二进制静态加密（流量归档等）。
func BlobEncryptionAvailable() bool { return APIKeyEncryptionAvailable() }

// EncryptBlob 用主密钥对任意二进制数据做 AES-256-GCM，返回 nonce||ciphertext。
// 用于流量归档/证据包这类「落盘即敏感」的数据：抓到的报文里经常带 Cookie、
// Authorization、密码等明文凭据，仅靠 0600 权限挡不住同机其他进程或误打包。
// 没有主密钥时返回错误（调用方决定是否退化为明文并显式记录）。
func EncryptBlob(plain []byte) ([]byte, error) {
	key, err := loadMasterKey()
	if err != nil || len(key) != 32 {
		if err == nil {
			err = errors.New("blob encryption is unavailable")
		}
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plain, nil), nil
}

// DecryptBlob 解开 EncryptBlob 的输出（nonce||ciphertext）。
func DecryptBlob(blob []byte) ([]byte, error) {
	key, err := loadMasterKey()
	if err != nil || len(key) != 32 {
		if err == nil {
			err = errors.New("blob encryption is unavailable")
		}
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(blob) < gcm.NonceSize() {
		return nil, errors.New("blob 太短，无法解密")
	}
	nonce, ct := blob[:gcm.NonceSize()], blob[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ct, nil)
}

// EncryptAPIKey wraps a plaintext key. Empty input and unstructured input pass
// through unchanged; already-encrypted values are returned as-is so the function
// is idempotent.
func EncryptAPIKey(plain string) (string, error) {
	if plain == "" || strings.HasPrefix(plain, apiKeyEncPrefix) {
		return plain, nil
	}
	key, err := loadMasterKey()
	if err != nil || len(key) != 32 {
		if err == nil {
			err = errors.New("API key encryption is unavailable")
		}
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nil, nonce, []byte(plain), nil)
	blob := append(nonce, sealed...)
	return apiKeyEncPrefix + base64.StdEncoding.EncodeToString(blob), nil
}

// DecryptAPIKey unwraps a stored value. Plaintext (unstamped) input, an unset
// master key, or a value encrypted with a different key are all returned
// unchanged rather than erroring: a provider built with a wrong key fails
// authentication upstream, which is far easier to diagnose than a boot failure.
func DecryptAPIKey(stored string) string {
	if stored == "" || !strings.HasPrefix(stored, apiKeyEncPrefix) {
		return stored
	}
	key, err := loadMasterKey()
	if err != nil || len(key) != 32 {
		return stored
	}
	blob, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(stored, apiKeyEncPrefix))
	if err != nil {
		return stored
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return stored
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil || len(blob) < gcm.NonceSize() {
		return stored
	}
	nonce, ct := blob[:gcm.NonceSize()], blob[gcm.NonceSize():]
	out, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return stored
	}
	return string(out)
}

// DecryptAPIKeyStrict 与 DecryptAPIKey 的区别：密文（enc:v1: 前缀）无法解密时返回
// 明确 error，而不是把密文原样返回。明文（无前缀）与空值仍按原样返回（nil error），
// 保证旧数据兼容。错误信息只描述失败原因，绝不包含密钥或明文。
func DecryptAPIKeyStrict(stored string) (string, error) {
	if stored == "" || !strings.HasPrefix(stored, apiKeyEncPrefix) {
		return stored, nil
	}
	key, err := loadMasterKey()
	if err != nil || len(key) != 32 {
		if err == nil {
			err = errors.New("API key encryption is unavailable")
		}
		return "", fmt.Errorf("主密钥不可用: %w", err)
	}
	blob, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(stored, apiKeyEncPrefix))
	if err != nil {
		return "", fmt.Errorf("密文 base64 解析失败: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("AES 初始化失败: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM 初始化失败: %w", err)
	}
	if len(blob) < gcm.NonceSize() {
		return "", errors.New("密文长度不足")
	}
	nonce, ct := blob[:gcm.NonceSize()], blob[gcm.NonceSize():]
	out, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("GCM 认证失败（密文被篡改或密钥不匹配）: %w", err)
	}
	return string(out), nil
}

// APIKeyEncryptionState 返回 API key 静态加密的当前状态，供设置页/日志展示。
// 只返回 enabled/source/error，绝不返回密钥或明文。
func APIKeyEncryptionState() map[string]any {
	k, err := loadMasterKey()
	enabled := err == nil && len(k) == 32
	source := ""
	if strings.TrimSpace(os.Getenv("ARTEX_KEY_SECRET")) != "" {
		source = "ARTEX_KEY_SECRET"
	} else if p := strings.TrimSpace(os.Getenv("ARTEX_KEY_SECRET_FILE")); p != "" {
		source = p
	}
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return map[string]any{
		"enabled": enabled,
		"source":  source,
		"error":   errStr,
	}
}

// IsEncryptedAPIKey reports whether a stored value carries the encryption stamp.
func IsEncryptedAPIKey(stored string) bool {
	return strings.HasPrefix(stored, apiKeyEncPrefix)
}

// MigrateAPIKeys encrypts any plaintext api_key left over from before this
// feature existed. Idempotent; safe to run on every boot. When no master key is
// configured it reports how many rows are still plaintext and does nothing.
func (d *DB) MigrateAPIKeys() error {
	if !APIKeyEncryptionAvailable() {
		_, err := loadMasterKey()
		log.Printf("[secretbox] 未启用 API key 加密（%v）；新 API key 将拒绝保存，请配置 ARTEX_KEY_SECRET_FILE", err)
		return nil
	}
	rows, err := d.Query(`SELECT id, COALESCE(api_key,'') FROM llm_profiles`)
	if err != nil {
		return err
	}
	defer rows.Close()
	type pending struct {
		id  int64
		key string
	}
	var todo []pending
	for rows.Next() {
		var id int64
		var k string
		if err := rows.Scan(&id, &k); err != nil {
			return err
		}
		if k != "" && !IsEncryptedAPIKey(k) {
			todo = append(todo, pending{id, k})
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(todo) == 0 {
		return nil
	}
	n := 0
	for _, p := range todo {
		enc, err := EncryptAPIKey(p.key)
		if err != nil || !IsEncryptedAPIKey(enc) {
			log.Printf("[secretbox] profile %d 加密失败，保持明文: %v", p.id, err)
			continue
		}
		if _, err := d.Exec(`UPDATE llm_profiles SET api_key=$1 WHERE id=$2`, enc, p.id); err != nil {
			log.Printf("[secretbox] profile %d 写回失败: %v", p.id, err)
			continue
		}
		n++
	}
	log.Printf("[secretbox] 已加密 %d/%d 个 profile 的 api_key（其余未能加密，仍为明文）", n, len(todo))
	return nil
}

// ArchiveAEAD 返回用于「流式分块加密」的 AES-256-GCM（任务归档 tar.zst 用）。
// 主密钥不可用时返回 error，调用方据此决定：导出退化为明文（并记录告警），
// 导入遇到加密归档则必须报错——绝不把密文当明文解析。
func ArchiveAEAD() (cipher.AEAD, error) {
	key, err := loadMasterKey()
	if err != nil || len(key) != 32 {
		if err == nil {
			err = errors.New("archive encryption is unavailable")
		}
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}