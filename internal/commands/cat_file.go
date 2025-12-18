package commands

import (
	"fmt"
	"io"
	"os"

	"ryot/internal/objects"
)

func CatFile(args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return UsageError{Msg: "cat-file: cat-file <sha1>"}
	}
	id := args[0]
	store := objects.Store{ObjectsDir: objectsDir}
	obj, err := store.ReadObject(id)
	if err != nil {
		return err
	}

	tf, err := os.CreateTemp(".", "temp_git_file_*")
	if err != nil {
		return err
	}
	defer tf.Close()

	typeStr := obj.Type
	if _, err := tf.Write(obj.Data); err != nil {
		typeStr = "bad"
	}
	fmt.Fprintf(stdout, "%s: %s\n", tf.Name(), typeStr)
	return nil
}
