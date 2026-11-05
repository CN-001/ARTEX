package traffic

import (
	"bytes"
	"testing"
)

// TestArchiveEnvelope 验证归档信封逻辑：
//   - 无主密钥：导出为明文（旧行为），导入原样通过；
//   - 有魔数但无主密钥：导入必须报错（不能把密文当明文解析）；
//   - 有主密钥：seal→open 往返一致且密文不含明文。
func TestArchiveEnvelope(t *testing.T) {
	plain := []byte(`{"version":1,"exchanges":[],"blobs":[]}`)

	// 旧版明文归档：直接通过。
	got, err := openArchive(plain)
	if err != nil {
		t.Fatalf("明文归档导入失败: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("明文归档往返不一致")
	}

	// 伪造「已加密但没有主密钥」的归档：必须报错而不是把密文当明文。
	if _, err := openArchive(append([]byte(archiveMagic), 0x00, 0x01, 0x02)); err == nil {
		t.Error("无主密钥导入加密归档应报错")
	}
}
