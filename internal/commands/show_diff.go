package commands

import (
	"fmt"
	"io"
	"os"

	"ryot/internal/diff"
	"ryot/internal/index"
	"ryot/internal/objects"
)

func ShowDiff(args []string, stdout io.Writer) error {
	_ = args
	idx, err := index.Read(indexPath)
	if err != nil {
		return err
	}
	store := objects.Store{ObjectsDir: objectsDir}

	for _, e := range idx.Entries {
		cur, err := os.ReadFile(e.Path)
		if err != nil {
			fmt.Fprintf(stdout, "%s: %v\n", e.Path, err)
			continue
		}
		blobID := sha1BytesToHex(e.Sha1)
		obj, err := store.ReadObject(blobID)
		if err != nil {
			fmt.Fprintf(stdout, "%s: unable to read blob %s: %v\n", e.Path, blobID, err)
			continue
		}
		if obj.Type != "blob" {
			fmt.Fprintf(stdout, "%s: expected blob, got %s\n", e.Path, obj.Type)
			continue
		}
		if bytesEqual(obj.Data, cur) {
			fmt.Fprintf(stdout, "%s: ok\n", e.Path)
			continue
		}

		out := diff.UnifiedDiff("a/"+e.Path, "b/"+e.Path, string(obj.Data), string(cur))
		fmt.Fprint(stdout, out)
	}
	return nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
