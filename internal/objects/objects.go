package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Store struct {
	ObjectsDir string // e.g. .dircache/objects
}

func (s Store) ObjectPath(hexID string) (dir string, file string, err error) {
	if len(hexID) != 40 {
		return "", "", fmt.Errorf("invalid sha1 length: %d", len(hexID))
	}
	_, err = hex.DecodeString(hexID)
	if err != nil {
		return "", "", fmt.Errorf("invalid sha1 hex: %w", err)
	}
	dir = filepath.Join(s.ObjectsDir, hexID[:2])
	file = filepath.Join(dir, hexID[2:])
	return dir, file, nil
}

func BuildObject(typ string, data []byte) []byte {
	hdr := fmt.Sprintf("%s %d\x00", typ, len(data))
	b := make([]byte, 0, len(hdr)+len(data))
	b = append(b, []byte(hdr)...)
	b = append(b, data...)
	return b
}

func Compress(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := zw.Write(in); err != nil {
		_ = zw.Close()
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Inflate(in []byte) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	out, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func HashCompressed(compressed []byte) (hexID string) {
	sum := sha1.Sum(compressed)
	return hex.EncodeToString(sum[:])
}

func (s Store) WriteObject(typ string, data []byte) (hexID string, err error) {
	raw := BuildObject(typ, data)
	compressed, err := Compress(raw)
	if err != nil {
		return "", err
	}
	hexID = HashCompressed(compressed)
	dir, file, err := s.ObjectPath(hexID)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	// Create-if-not-exists semantics.
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return hexID, nil
		}
		return "", err
	}
	defer f.Close()
	if _, err := f.Write(compressed); err != nil {
		return "", err
	}
	return hexID, nil
}

type Object struct {
	Type string
	Size int
	Data []byte
}

func ParseInflated(inflated []byte) (Object, error) {
	nul := bytes.IndexByte(inflated, 0)
	if nul < 0 {
		return Object{}, fmt.Errorf("invalid object: missing NUL")
	}
	hdr := string(inflated[:nul])
	sp := strings.IndexByte(hdr, ' ')
	if sp < 0 {
		return Object{}, fmt.Errorf("invalid object header: %q", hdr)
	}
	typ := hdr[:sp]
	szStr := hdr[sp+1:]
	sz, err := strconv.Atoi(szStr)
	if err != nil {
		return Object{}, fmt.Errorf("invalid size: %w", err)
	}
	data := inflated[nul+1:]
	if sz != len(data) {
		return Object{}, fmt.Errorf("size mismatch: header=%d actual=%d", sz, len(data))
	}
	return Object{Type: typ, Size: sz, Data: data}, nil
}

func (s Store) ReadObject(hexID string) (Object, error) {
	_, file, err := s.ObjectPath(hexID)
	if err != nil {
		return Object{}, err
	}
	compressed, err := os.ReadFile(file)
	if err != nil {
		return Object{}, err
	}
	inflated, err := Inflate(compressed)
	if err != nil {
		return Object{}, err
	}
	return ParseInflated(inflated)
}

func (s Store) Exists(hexID string) bool {
	_, file, err := s.ObjectPath(hexID)
	if err != nil {
		return false
	}
	_, err = os.Stat(file)
	return err == nil
}
