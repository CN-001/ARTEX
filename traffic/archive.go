package traffic

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Aurox/artex/db"
)

// archiveMagic 标记「主密钥加密的归档文件」。
//
// 为什么归档需要加密：抓到的流量里经常带 Cookie / Authorization / 表单密码，
// 0600 只挡同机其他用户，挡不住误打包、备份外发、镜像层残留。这里用与 LLM key
// 相同的主密钥做 AES-256-GCM；没有主密钥时退化为明文（旧行为），并在导出日志里
// 明确记录，避免"以为加密了其实没有"。
const archiveMagic = "ARTXENC1"

// sealArchive 有主密钥则加密（前置魔数），否则原样返回。
func sealArchive(data []byte) ([]byte, error) {
	if !db.BlobEncryptionAvailable() {
		return data, nil
	}
	sealed, err := db.EncryptBlob(data)
	if err != nil {
		return nil, err
	}
	return append([]byte(archiveMagic), sealed...), nil
}

// openArchive 解开 sealArchive 的输出；没有魔数视为旧版明文归档，原样返回。
func openArchive(raw []byte) ([]byte, error) {
	if !bytes.HasPrefix(raw, []byte(archiveMagic)) {
		return raw, nil
	}
	if !db.BlobEncryptionAvailable() {
		return nil, errors.New("归档已加密，但当前未配置主密钥（ARTEX_KEY_SECRET/_FILE），无法导入")
	}
	return db.DecryptBlob(raw[len(archiveMagic):])
}

// ArchiveSnapshot is the portable traffic subset embedded in one task archive.
// Large content-addressed bodies are copied alongside this manifest rather than
// base64-encoded, preserving deduplication and allowing the tar writer to stream.
type ArchiveSnapshot struct {
	Version   int               `json:"version"`
	Exchanges []ArchiveExchange `json:"exchanges"`
	Blobs     []string          `json:"blobs"`
}

type ArchiveExchange struct {
	ID          string `json:"id"`
	TS          int64  `json:"ts"`
	Host        string `json:"host"`
	Method      string `json:"method"`
	URLTemplate string `json:"url_template"`
	URL         string `json:"url"`
	Status      int    `json:"status"`
	ContentType string `json:"content_type"`
	ReqLen      int    `json:"req_len"`
	RespLen     int    `json:"resp_len"`
	ReqHead     string `json:"req_head"`
	ReqBody     []byte `json:"req_body,omitempty"`
	ReqBlob     string `json:"req_blob,omitempty"`
	RespHead    string `json:"resp_head"`
	RespBody    []byte `json:"resp_body,omitempty"`
	RespBlob    string `json:"resp_blob,omitempty"`
}

// ExportHosts writes an exact-host snapshot to dir. It holds the traffic writer
// lock while reading SQLite and blobs, so each body and its index row come from
// one consistent point in time.
func (t *Traffic) ExportHosts(hosts []string, dir string) (int64, error) {
	if t == nil || len(hosts) == 0 {
		return 0, nil
	}
	if err := os.MkdirAll(filepath.Join(dir, "blobs"), 0o700); err != nil {
		return 0, err
	}
	t.wmu.Lock()
	defer t.wmu.Unlock()
	unique := uniqueArchiveHosts(hosts)
	placeholders := make([]string, len(unique))
	args := make([]any, len(unique))
	for i, host := range unique {
		placeholders[i] = "?"
		args[i] = host
	}
	rows, err := t.db.Query(`SELECT id,ts,host,method,url_template,url,status,content_type,req_len,resp_len,path
FROM exchanges WHERE host IN (`+strings.Join(placeholders, ",")+`) ORDER BY ts,id`, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	snapshot := ArchiveSnapshot{Version: 1}
	blobs := map[string]struct{}{}
	for rows.Next() {
		var item ArchiveExchange
		var legacyPath string
		if err := rows.Scan(&item.ID, &item.TS, &item.Host, &item.Method, &item.URLTemplate,
			&item.URL, &item.Status, &item.ContentType, &item.ReqLen, &item.RespLen, &legacyPath); err != nil {
			return 0, err
		}
		var reqBlob, respBlob sql.NullString
		err := t.db.QueryRow(`SELECT req_head,req_body,req_blob,resp_head,resp_body,resp_blob
FROM exchange_bodies WHERE id=?`, item.ID).Scan(&item.ReqHead, &item.ReqBody, &reqBlob, &item.RespHead, &item.RespBody, &respBlob)
		if errors.Is(err, sql.ErrNoRows) && strings.TrimSpace(legacyPath) != "" {
			req, readErr := os.ReadFile(filepath.Join(t.dir, legacyPath, "request.http"))
			if readErr != nil && !os.IsNotExist(readErr) {
				return 0, readErr
			}
			resp, readErr := os.ReadFile(filepath.Join(t.dir, legacyPath, "response.http"))
			if readErr != nil && !os.IsNotExist(readErr) {
				return 0, readErr
			}
			item.ReqHead, item.RespHead = string(req), string(resp)
		} else if err != nil {
			return 0, err
		}
		item.ReqBlob, item.RespBlob = reqBlob.String, respBlob.String
		for _, hash := range []string{item.ReqBlob, item.RespBlob} {
			if hash != "" {
				blobs[hash] = struct{}{}
			}
		}
		snapshot.Exchanges = append(snapshot.Exchanges, item)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for hash := range blobs {
		source, err := t.blobPath(hash)
		if err != nil {
			return 0, err
		}
		data, err := os.ReadFile(source)
		if err != nil {
			return 0, err
		}
		sealed, err := sealArchive(data)
		if err != nil {
			return 0, err
		}
		if err := os.WriteFile(filepath.Join(dir, "blobs", hash+".bin"), sealed, 0o600); err != nil {
			return 0, err
		}
		snapshot.Blobs = append(snapshot.Blobs, hash)
	}
	sort.Strings(snapshot.Blobs)
	raw, err := json.Marshal(snapshot)
	if err != nil {
		return 0, err
	}
	if !db.BlobEncryptionAvailable() {
		log.Printf("[traffic] 警告: 未配置主密钥，归档以明文写出（配置 ARTEX_KEY_SECRET 后自动加密）")
	}
	sealedJSON, err := sealArchive(raw)
	if err != nil {
		return 0, err
	}
	if err := os.WriteFile(filepath.Join(dir, "traffic.json"), sealedJSON, 0o600); err != nil {
		return 0, err
	}
	return int64(len(snapshot.Exchanges)), nil
}

// ImportArchive imports only missing exchange IDs. Current hot rows always win,
// and repeated restore attempts are safe after a partial external failure.
func (t *Traffic) ImportArchive(dir string) (int64, error) {
	if t == nil {
		return 0, nil
	}
	raw, err := os.ReadFile(filepath.Join(dir, "traffic.json"))
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	plain, err := openArchive(raw)
	if err != nil {
		return 0, err
	}
	var snapshot ArchiveSnapshot
	if err := json.Unmarshal(plain, &snapshot); err != nil {
		return 0, err
	}
	if snapshot.Version != 1 {
		return 0, fmt.Errorf("unsupported traffic archive version %d", snapshot.Version)
	}
	t.wmu.Lock()
	defer t.wmu.Unlock()
	for _, hash := range snapshot.Blobs {
		if !blobHashRe.MatchString(hash) {
			return 0, fmt.Errorf("invalid archived traffic blob %q", hash)
		}
		rawBlob, err := os.ReadFile(filepath.Join(dir, "blobs", hash+".bin"))
		if err != nil {
			return 0, err
		}
		data, err := openArchive(rawBlob) // 校验和针对明文计算
		if err != nil {
			return 0, err
		}
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != hash {
			return 0, fmt.Errorf("traffic blob checksum mismatch: %s", hash)
		}
		destination := filepath.Join(t.dir, "_blobs", "sha256", hash[:2], hash+".bin")
		if _, err := os.Stat(destination); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
				return 0, err
			}
			if err := os.WriteFile(destination, data, 0o644); err != nil {
				return 0, err
			}
		}
	}
	tx, err := t.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck
	var imported int64
	for _, item := range snapshot.Exchanges {
		res, err := tx.Exec(`INSERT OR IGNORE INTO exchanges(id,ts,host,method,url_template,url,status,content_type,req_len,resp_len,path)
VALUES(?,?,?,?,?,?,?,?,?,?,'')`, item.ID, item.TS, item.Host, item.Method, item.URLTemplate,
			item.URL, item.Status, item.ContentType, item.ReqLen, item.RespLen)
		if err != nil {
			return imported, err
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO exchange_bodies(id,req_head,req_body,req_blob,resp_head,resp_body,resp_blob)
VALUES(?,?,?,?,?,?,?)`, item.ID, item.ReqHead, item.ReqBody, nullIfEmpty(item.ReqBlob),
			item.RespHead, item.RespBody, nullIfEmpty(item.RespBlob)); err != nil {
			return imported, err
		}
		for _, hash := range []string{item.ReqBlob, item.RespBlob} {
			if hash != "" {
				if _, err := tx.Exec(`INSERT OR IGNORE INTO blob_refs(hash,exchange_id) VALUES(?,?)`, hash, item.ID); err != nil {
					return imported, err
				}
			}
		}
		if t.fts {
			var rowID int64
			if err := tx.QueryRow(`SELECT rowid FROM exchanges WHERE id=?`, item.ID).Scan(&rowID); err != nil {
				return imported, err
			}
			content := strings.Join([]string{item.URL, item.ReqHead, string(item.ReqBody), item.RespHead, string(item.RespBody)}, "\n")
			if _, err := tx.Exec(`INSERT INTO ex_fts(rowid,content) VALUES(?,?)`, rowID, content); err != nil {
				return imported, err
			}
		}
		imported++
	}
	if err := tx.Commit(); err != nil {
		return imported, err
	}
	return imported, nil
}

func uniqueArchiveHosts(hosts []string) []string {
	seen := make(map[string]struct{}, len(hosts))
	out := make([]string, 0, len(hosts))
	for _, host := range hosts {
		host = strings.ToLower(strings.TrimSpace(host))
		if host == "" {
			continue
		}
		if _, exists := seen[host]; exists {
			continue
		}
		seen[host] = struct{}{}
		out = append(out, host)
	}
	sort.Strings(out)
	return out
}
