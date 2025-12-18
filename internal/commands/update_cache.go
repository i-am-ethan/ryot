package commands

import (
	"fmt"
	"io"
	"os"

	"ryot/internal/index"
	"ryot/internal/objects"
)

func UpdateCache(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return UsageError{Msg: "update-cache: update-cache <file> [<file> ...]"}
	}
	store := objects.Store{ObjectsDir: objectsDir}

	idx, err := index.Read(indexPath)
	if err != nil {
		return err
	}

	for _, p := range args {
		if !verifyPath(p) {
			fmt.Fprintf(stdout, "Ignoring path %s\n", p)
			continue
		}
		fi, err := os.Stat(p)
		if err != nil {
			return err
		}
		if !fi.Mode().IsRegular() {
			return fmt.Errorf("only regular files supported: %s", p)
		}
		content, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		blobID, err := store.WriteObject("blob", content)
		if err != nil {
			return err
		}
		shaBytes, err := hexToSha1Bytes(blobID)
		if err != nil {
			return err
		}
		idx.Upsert(index.Entry{
			Path:     p,
			Mode:     fileModeToGitMode(fi),
			MtimeSec: fi.ModTime().Unix(),
			Size:     fi.Size(),
			Sha1:     shaBytes,
		})
		fmt.Fprintf(stdout, "%s %s\n", p, blobID)
	}

	idx.Sort()
	return index.WriteAtomic(indexPath, idx)
}
