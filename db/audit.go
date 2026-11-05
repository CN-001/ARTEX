package db

import (
	"context"
	"database/sql"
	"time"
)

// AuditEntry 是一条审计记录。
//
// 设计取舍：审计日志只记「动作与结果」，不记请求体/响应体——那些里面常有密码、
// API key、Cookie，把敏感数据抄进审计表等于制造第二个泄露面。需要细节时用
// detail 记短摘要（截断、去敏感字段）。
type AuditEntry struct {
	ID     int64     `json:"id"`
	TS     time.Time `json:"ts"`
	Actor  string    `json:"actor"`
	IP     string    `json:"ip"`
	Method string    `json:"method"`
	Path   string    `json:"path"`
	Action string    `json:"action"`
	Target string    `json:"target"`
	Status int       `json:"status"`
	Detail string    `json:"detail"`
}

// InsertAuditCtx 带 context 写入一条审计记录（供异步审计路径用超时控制）。
func (d *DB) InsertAuditCtx(ctx context.Context, e AuditEntry) error {
	_, err := d.ExecContext(ctx, `
INSERT INTO audit_log (actor, ip, method, path, action, target, status, detail)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		e.Actor, e.IP, e.Method, e.Path, e.Action, e.Target, e.Status, e.Detail)
	return err
}

// InsertAudit 写入一条审计记录。调用方负责截断过长字段（见 server 侧 truncate）。
func (d *DB) InsertAudit(e AuditEntry) error {
	_, err := d.Exec(`
INSERT INTO audit_log (actor, ip, method, path, action, target, status, detail)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		e.Actor, e.IP, e.Method, e.Path, e.Action, e.Target, e.Status, e.Detail)
	return err
}

// ListAudit 按时间倒序返回审计记录；beforeID>0 时做游标分页（取更早的记录）。
func (d *DB) ListAudit(limit int, beforeID int64) ([]AuditEntry, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	var (
		rows *sql.Rows
		err  error
	)
	if beforeID > 0 {
		rows, err = d.Query(`
SELECT id, ts, actor, ip, method, path, action, target, status, detail
FROM audit_log WHERE id < $1 ORDER BY id DESC LIMIT $2`, beforeID, limit)
	} else {
		rows, err = d.Query(`
SELECT id, ts, actor, ip, method, path, action, target, status, detail
FROM audit_log ORDER BY id DESC LIMIT $1`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.TS, &e.Actor, &e.IP, &e.Method, &e.Path, &e.Action, &e.Target, &e.Status, &e.Detail); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// PruneAudit 删除早于保留期的审计记录，返回删除行数。默认保留 180 天。
func (d *DB) PruneAudit(retention time.Duration) (int64, error) {
	if retention <= 0 {
		retention = 180 * 24 * time.Hour
	}
	res, err := d.Exec(`DELETE FROM audit_log WHERE ts < now() - $1::interval`,
		retention.String())
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
