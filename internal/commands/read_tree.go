package commands

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"ryot/internal/objects"
)

func ReadTree(args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return UsageError{Msg: "read-tree: read-tree <key>"}
	}
	id := args[0]
	store := objects.Store{ObjectsDir: objectsDir}
	obj, err := store.ReadObject(id)
	if err != nil {
		return err
	}
	if obj.Type != "tree" {
		return fmt.Errorf("expected a 'tree' node, got %q", obj.Type)
	}

	buf := obj.Data
	for len(buf) > 0 {
		nul := bytes.IndexByte(buf, 0)
		if nul < 0 {
			return fmt.Errorf("corrupt tree: missing NUL")
		}
		entry := buf[:nul]
		rest := buf[nul+1:]
		if len(rest) < 20 {
			return fmt.Errorf("corrupt tree: missing sha1")
		}
		sha := rest[:20]
		buf = rest[20:]

		sp := bytes.IndexByte(entry, ' ')
		if sp < 0 {
			return fmt.Errorf("corrupt tree entry")
		}
		modeStr := string(entry[:sp])
		path := string(entry[sp+1:])
		mode, err := strconv.ParseUint(modeStr, 8, 32)
		if err != nil {
			return fmt.Errorf("corrupt tree mode: %w", err)
		}
		fmt.Fprintf(stdout, "%o %s (%x)\n", mode, path, sha)
	}
	return nil
}
