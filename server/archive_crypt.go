package server

import (
	"bufio"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// 任务归档（tar.zst）的静态加密：分块 AES-256-GCM。
//
// 为什么不是"整个文件一次 Seal"：归档可能有 GB 级（含流量/转录/工作区文件），
// 一次性读入内存会 OOM。这里按 64KiB 分块加密，每块一个 nonce；nonce = 文件头里
// 的 12 字节随机 base 与块计数器异或，保证同一主密钥下每次导出都不重复。
//
// 文件格式（先写头，再写块）：
//
//	"ARTXTAR1" | baseNonce[12] | [uint32 len][ciphertext+tag] × N
//
// 兼容：没有魔数的文件按旧版明文 tar.zst 处理（老归档仍可导入）。
const (
	taskArchiveMagic     = "ARTXTAR1"
	archiveChunkSize     = 64 * 1024
	archiveNonceSize     = 12
	taskArchiveHeaderLen = len(taskArchiveMagic) + archiveNonceSize
)

// archiveEncWriter 把写入的明文分块加密后落盘。
type archiveEncWriter struct {
	w    io.Writer
	aead cipher.AEAD
	base []byte
	ctr  uint64
	buf  []byte
}

// newArchiveEncWriter 写入文件头（魔数 + 随机 base nonce）。
func newArchiveEncWriter(w io.Writer, aead cipher.AEAD) (*archiveEncWriter, error) {
	base := make([]byte, archiveNonceSize)
	if _, err := io.ReadFull(rand.Reader, base); err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(taskArchiveMagic)); err != nil {
		return nil, err
	}
	if _, err := w.Write(base); err != nil {
		return nil, err
	}
	return &archiveEncWriter{w: w, aead: aead, base: base, buf: make([]byte, 0, archiveChunkSize)}, nil
}

func (e *archiveEncWriter) Write(p []byte) (int, error) {
	written := 0
	for len(p) > 0 {
		space := archiveChunkSize - len(e.buf)
		if space > len(p) {
			space = len(p)
		}
		e.buf = append(e.buf, p[:space]...)
		p = p[space:]
		written += space
		if len(e.buf) == archiveChunkSize {
			if err := e.flush(); err != nil {
				return written, err
			}
		}
	}
	return written, nil
}

// Close 把最后一块（不足 64KiB 的尾块）写出。
func (e *archiveEncWriter) Close() error {
	if len(e.buf) == 0 {
		return nil
	}
	return e.flush()
}

func (e *archiveEncWriter) flush() error {
	nonce := make([]byte, archiveNonceSize)
	copy(nonce, e.base)
	// 计数器异或进 nonce 末尾 8 字节（不会溢出到 base 的其余字节）。
	binary.BigEndian.PutUint64(nonce[archiveNonceSize-8:], binary.BigEndian.Uint64(nonce[archiveNonceSize-8:])^e.ctr)
	e.ctr++
	sealed := e.aead.Seal(nil, nonce, e.buf, nil)
	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(len(sealed)))
	if _, err := e.w.Write(hdr[:]); err != nil {
		return err
	}
	if _, err := e.w.Write(sealed); err != nil {
		return err
	}
	e.buf = e.buf[:0]
	return nil
}

// archiveDecReader 读取并解密分块密文。
type archiveDecReader struct {
	r    *bufio.Reader
	aead cipher.AEAD
	base []byte
	ctr  uint64
	buf  []byte
	eof  bool
}

func (d *archiveDecReader) Read(p []byte) (int, error) {
	for len(d.buf) == 0 {
		if d.eof {
			return 0, io.EOF
		}
		var hdr [4]byte
		if _, err := io.ReadFull(d.r, hdr[:]); err != nil {
			if errors.Is(err, io.EOF) {
				d.eof = true
				return 0, io.EOF
			}
			return 0, err
		}
		n := binary.BigEndian.Uint32(hdr[:])
		if n == 0 || n > archiveChunkSize+uint32(d.aead.Overhead()) {
			return 0, fmt.Errorf("归档分块长度非法: %d", n)
		}
		sealed := make([]byte, n)
		if _, err := io.ReadFull(d.r, sealed); err != nil {
			return 0, err
		}
		nonce := make([]byte, archiveNonceSize)
		copy(nonce, d.base)
		binary.BigEndian.PutUint64(nonce[archiveNonceSize-8:], binary.BigEndian.Uint64(nonce[archiveNonceSize-8:])^d.ctr)
		d.ctr++
		plain, err := d.aead.Open(nil, nonce, sealed, nil)
		if err != nil {
			return 0, errors.New("归档解密失败（主密钥不匹配或文件被篡改）")
		}
		d.buf = plain
	}
	n := copy(p, d.buf)
	d.buf = d.buf[n:]
	return n, nil
}

// newArchiveDecReader 检测文件头：有魔数则解密读取，否则按旧版明文原样读取。
func newArchiveDecReader(r io.Reader, aead cipher.AEAD) (io.Reader, error) {
	br := bufio.NewReaderSize(r, taskArchiveHeaderLen+archiveChunkSize)
	head, err := br.Peek(taskArchiveHeaderLen)
	if err != nil || string(head[:len(taskArchiveMagic)]) != taskArchiveMagic {
		return br, nil // 旧版明文归档（或文件过大读不到头时交给上层报错）
	}
	if aead == nil {
		return nil, errors.New("任务归档已加密，但当前未配置主密钥（ARTEX_KEY_SECRET/_FILE），无法导入")
	}
	if _, err := br.Discard(len(taskArchiveMagic)); err != nil {
		return nil, err
	}
	base := make([]byte, archiveNonceSize)
	if _, err := io.ReadFull(br, base); err != nil {
		return nil, err
	}
	return &archiveDecReader{r: br, aead: aead, base: base}, nil
}
