package commands

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func mustWriter(w io.Writer) io.Writer { return w }

func hexToSha1Bytes(hexID string) ([20]byte, error) {
	var out [20]byte
	if len(hexID) != 40 {
		return out, fmt.Errorf("invalid sha1 length: %d", len(hexID))
	}
	b, err := hex.DecodeString(hexID)
	if err != nil {
		return out, err
	}
	copy(out[:], b)
	return out, nil
}

func sha1BytesToHex(b [20]byte) string {
	return hex.EncodeToString(b[:])
}

// Early git's verify_path is stricter; we keep a simple version.
func verifyPath(p string) bool {
	if p == "" {
		return false
	}
	if strings.Contains(p, "//") {
		return false
	}
	if strings.HasSuffix(p, "/") {
		return false
	}
	clean := filepath.Clean(p)
	if clean != p {
		// Avoid ambiguous paths
		return false
	}
	parts := strings.Split(p, string(os.PathSeparator))
	for _, part := range parts {
		if part == "." || part == ".." {
			return false
		}
		if strings.HasPrefix(part, ".") {
			return false
		}
	}
	return true
}

func fileModeToGitMode(fi os.FileInfo) uint32 {
	// We only support regular files.
	perm := fi.Mode().Perm()
	exec := (perm & 0o111) != 0
	if exec {
		return 0o100755
	}
	return 0o100644
}

func sha1OfBytes(b []byte) [20]byte {
	return sha1.Sum(b)
}
