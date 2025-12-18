package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func InitDB(args []string, stdout io.Writer) error {
	// args unused (kept for symmetry)
	_ = args

	if err := os.MkdirAll(dircacheDir, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(objectsDir, 0o700); err != nil {
		return err
	}
	for i := 0; i < 256; i++ {
		p := filepath.Join(objectsDir, fmt.Sprintf("%02x", i))
		if err := os.MkdirAll(p, 0o700); err != nil {
			return err
		}
	}
	fmt.Fprintln(stdout, "initialized", dircacheDir)
	return nil
}
