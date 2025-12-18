package commands

import (
	"bytes"
	"fmt"
	"io"

	"ryot/internal/index"
	"ryot/internal/objects"
)

func WriteTree(args []string, stdout io.Writer) error {
	_ = args
	idx, err := index.Read(indexPath)
	if err != nil {
		return err
	}
	if len(idx.Entries) == 0 {
		return fmt.Errorf("no index entries")
	}
	idx.Sort()

	store := objects.Store{ObjectsDir: objectsDir}

	var body bytes.Buffer
	for _, e := range idx.Entries {
		line := fmt.Sprintf("%o %s", e.Mode, e.Path)
		body.WriteString(line)
		body.WriteByte(0)
		body.Write(e.Sha1[:])
	}

	treeID, err := store.WriteObject("tree", body.Bytes())
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, treeID)
	return nil
}
