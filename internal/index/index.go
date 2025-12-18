package index

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// これは「雰囲気優先」の自前バイナリindex。
// READMEの言う "current directory cache" の役割（パス→blob/メタ情報）を満たす。
// 本家最初期のDIRC互換ではない。

var (
	magic   = [8]byte{'R', 'Y', 'O', 'T', 'I', 'D', 'X', 1}
	endian  = binary.BigEndian
	version = uint32(1)
)

type Entry struct {
	Path     string
	Mode     uint32
	MtimeSec int64
	Size     int64
	Sha1     [20]byte
}

type Index struct {
	Entries []Entry
}

func Read(path string) (Index, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Index{}, nil
		}
		return Index{}, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	var gotMagic [8]byte
	if _, err := io.ReadFull(r, gotMagic[:]); err != nil {
		return Index{}, err
	}
	if gotMagic != magic {
		return Index{}, fmt.Errorf("bad index magic")
	}
	var ver uint32
	if err := binary.Read(r, endian, &ver); err != nil {
		return Index{}, err
	}
	if ver != version {
		return Index{}, fmt.Errorf("unsupported index version: %d", ver)
	}
	var n uint32
	if err := binary.Read(r, endian, &n); err != nil {
		return Index{}, err
	}

	entries := make([]Entry, 0, n)
	for i := uint32(0); i < n; i++ {
		var pathLen uint16
		if err := binary.Read(r, endian, &pathLen); err != nil {
			return Index{}, err
		}
		var mode uint32
		var mtime int64
		var size int64
		if err := binary.Read(r, endian, &mode); err != nil {
			return Index{}, err
		}
		if err := binary.Read(r, endian, &mtime); err != nil {
			return Index{}, err
		}
		if err := binary.Read(r, endian, &size); err != nil {
			return Index{}, err
		}
		var sha [20]byte
		if _, err := io.ReadFull(r, sha[:]); err != nil {
			return Index{}, err
		}
		p := make([]byte, pathLen)
		if _, err := io.ReadFull(r, p); err != nil {
			return Index{}, err
		}
		entries = append(entries, Entry{Path: string(p), Mode: mode, MtimeSec: mtime, Size: size, Sha1: sha})
	}
	return Index{Entries: entries}, nil
}

func WriteAtomic(path string, idx Index) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	tmp := path + ".lock"
	buf := new(bytes.Buffer)
	buf.Write(magic[:])
	_ = binary.Write(buf, endian, version)
	_ = binary.Write(buf, endian, uint32(len(idx.Entries)))

	for _, e := range idx.Entries {
		p := []byte(e.Path)
		if len(p) > 0xFFFF {
			return fmt.Errorf("path too long: %s", e.Path)
		}
		_ = binary.Write(buf, endian, uint16(len(p)))
		_ = binary.Write(buf, endian, e.Mode)
		_ = binary.Write(buf, endian, e.MtimeSec)
		_ = binary.Write(buf, endian, e.Size)
		buf.Write(e.Sha1[:])
		buf.Write(p)
	}

	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (idx *Index) Upsert(e Entry) {
	for i := range idx.Entries {
		if idx.Entries[i].Path == e.Path {
			idx.Entries[i] = e
			return
		}
	}
	idx.Entries = append(idx.Entries, e)
}

func (idx *Index) Sort() {
	sort.Slice(idx.Entries, func(i, j int) bool {
		return idx.Entries[i].Path < idx.Entries[j].Path
	})
}
