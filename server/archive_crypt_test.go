package server

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pgdb "github.com/Aurox/artex/db"
)

// wireArchiveTestKey 在测试里配置一个 32 字节主密钥（重置 db 的 sync.Once 缓存）。
func wireArchiveTestKey(t *testing.T) {
	t.Helper()
	pgdb.ResetMasterKeyCacheForTest()
	t.Cleanup(pgdb.ResetMasterKeyCacheForTest)
	t.Setenv("ARTEX_KEY_SECRET", base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x31}, 32)))
}

// TestTaskArchiveEncryptionRoundTrip 验证任务归档在配置主密钥后：
//   - 磁盘上的文件不含明文标记、含 ARTXTAR1 魔数；
//   - 导入仍能还原文件内容；
//   - 明文旧归档仍然兼容。
func TestTaskArchiveEncryptionRoundTrip(t *testing.T) {
	wireArchiveTestKey(t)

	dataDir := t.TempDir()
	taskID := "77"
	explorationID := int64(88)
	taskFile := filepath.Join(dataDir, "tasks", taskID, "uploads", "secret.txt")
	if err := os.MkdirAll(filepath.Dir(taskFile), 0o700); err != nil {
		t.Fatal(err)
	}
	const marker = "PLAINTEXT-MARKER-SHOULD-NOT-APPEAR"
	if err := os.WriteFile(taskFile, []byte(marker), 0o600); err != nil {
		t.Fatal(err)
	}

	stage, err := stageTaskArchiveFiles(dataDir, 2, taskID, explorationID)
	if err != nil {
		t.Fatal(err)
	}
	snapshot := &pgdb.TaskArchiveSnapshot{
		FormatVersion: pgdb.TaskArchiveFormatVersion, TaskID: 77, ExplorationID: explorationID,
	}
	archivePath := taskArchivePath(dataDir, 2, taskID)
	if _, _, _, err := writeTaskArchivePackage(archivePath, stage.payload, snapshot); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(raw, []byte(taskArchiveMagic)) {
		t.Fatalf("归档没有加密头，前 8 字节=%q", raw[:8])
	}
	if bytes.Contains(raw, []byte(marker)) {
		t.Fatal("加密归档中出现了明文")
	}

	// 导入：解密 + 解压回原文件。
	outDir := filepath.Join(t.TempDir(), "restore")
	if err := extractTaskArchivePackage(archivePath, "", outDir); err != nil {
		t.Fatal(err)
	}
	restored, err := os.ReadFile(filepath.Join(outDir, "files", "tasks", taskID, "uploads", "secret.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(restored) != marker {
		t.Fatalf("还原内容不一致: %q", restored)
	}

	// 旧版明文归档兼容：无主密钥导出 → 明文，仍可被同样的导入路径读取。
	pgdb.ResetMasterKeyCacheForTest()
	t.Setenv("ARTEX_KEY_SECRET", "")
	t.Setenv("ARTEX_KEY_SECRET_FILE", "")
	plainDir := t.TempDir()
	plainTask := "78"
	plainFile := filepath.Join(plainDir, "tasks", plainTask, "uploads", "plain.txt")
	if err := os.MkdirAll(filepath.Dir(plainFile), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(plainFile, []byte("legacy-plain"), 0o600); err != nil {
		t.Fatal(err)
	}
	plainStage, err := stageTaskArchiveFiles(plainDir, 3, plainTask, explorationID)
	if err != nil {
		t.Fatal(err)
	}
	plainSnapshot := &pgdb.TaskArchiveSnapshot{FormatVersion: pgdb.TaskArchiveFormatVersion, TaskID: 78}
	plainArchive := taskArchivePath(plainDir, 3, plainTask)
	if _, _, _, err := writeTaskArchivePackage(plainArchive, plainStage.payload, plainSnapshot); err != nil {
		t.Fatal(err)
	}
	rawPlain, _ := os.ReadFile(plainArchive)
	if bytes.HasPrefix(rawPlain, []byte(taskArchiveMagic)) {
		t.Fatal("无主密钥时不应写出加密头")
	}
	outPlain := filepath.Join(t.TempDir(), "restore-plain")
	if err := extractTaskArchivePackage(plainArchive, "", outPlain); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(outPlain, "files", "tasks", plainTask, "uploads", "plain.txt"))
	if err != nil || string(got) != "legacy-plain" {
		t.Fatalf("旧版明文归档导入失败: %v %q", err, got)
	}
}

// TestArchiveDecReaderRejectsWrongKey 主密钥不匹配时必须报错，而不是吐出乱数据。
func TestArchiveDecReaderRejectsWrongKey(t *testing.T) {
	wireArchiveTestKey(t)
	var buf bytes.Buffer
	aead, err := pgdb.ArchiveAEAD()
	if err != nil {
		t.Fatal(err)
	}
	enc, err := newArchiveEncWriter(&buf, aead)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := enc.Write([]byte(strings.Repeat("x", 1000))); err != nil {
		t.Fatal(err)
	}
	if err := enc.Close(); err != nil {
		t.Fatal(err)
	}

	// 换一把错误的主密钥。
	pgdb.ResetMasterKeyCacheForTest()
	t.Setenv("ARTEX_KEY_SECRET", base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x99}, 32)))
	wrong, err := pgdb.ArchiveAEAD()
	if err != nil {
		t.Fatal(err)
	}
	dec, err := newArchiveDecReader(bytes.NewReader(buf.Bytes()), wrong)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(dec); err == nil {
		t.Fatal("错误主密钥解密应当失败")
	}

	// 没有主密钥（aead=nil）导入加密归档必须明确报错。
	if _, err := newArchiveDecReader(bytes.NewReader(buf.Bytes()), nil); err == nil {
		t.Fatal("无主密钥导入加密归档应报错")
	}
	_ = cipher.AEAD(nil)
}
