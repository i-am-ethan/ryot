package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"ryot/internal/objects"
)

func CommitTree(args []string, stdout io.Writer, stdin io.Reader) error {
	if len(args) < 1 {
		return UsageError{Msg: "commit-tree: commit-tree <tree-sha1> [-p <parent-sha1>]* < message"}
	}
	treeID := args[0]
	parents := []string{}

	for i := 1; i < len(args); i++ {
		if args[i] != "-p" {
			return UsageError{Msg: "commit-tree: commit-tree <tree-sha1> [-p <parent-sha1>]* < message"}
		}
		if i+1 >= len(args) {
			return UsageError{Msg: "commit-tree: missing parent sha1"}
		}
		parents = append(parents, args[i+1])
		i++
	}

	msgBytes, err := io.ReadAll(bufio.NewReader(stdin))
	if err != nil {
		return err
	}
	message := string(msgBytes)
	if message == "" {
		message = "(no message)\n"
	}
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}

	name := os.Getenv("COMMITTER_NAME")
	if name == "" {
		name = os.Getenv("USER")
		if name == "" {
			name = "unknown"
		}
	}
	email := os.Getenv("COMMITTER_EMAIL")
	if email == "" {
		h, _ := os.Hostname()
		if h == "" {
			h = "localhost"
		}
		email = fmt.Sprintf("%s@%s", name, h)
	}
	date := os.Getenv("COMMITTER_DATE")
	if date == "" {
		// git-ish: "<unix> <tz>"
		now := time.Now()
		_, off := now.Zone()
		sign := "+"
		if off < 0 {
			sign = "-"
			off = -off
		}
		hh := off / 3600
		mm := (off % 3600) / 60
		date = fmt.Sprintf("%d %s%02d%02d", now.Unix(), sign, hh, mm)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("tree %s\n", treeID))
	for _, p := range parents {
		sb.WriteString(fmt.Sprintf("parent %s\n", p))
	}
	sb.WriteString(fmt.Sprintf("author %s <%s> %s\n", name, email, date))
	sb.WriteString(fmt.Sprintf("committer %s <%s> %s\n\n", name, email, date))
	sb.WriteString(message)

	store := objects.Store{ObjectsDir: objectsDir}
	commitID, err := store.WriteObject("commit", []byte(sb.String()))
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, commitID)
	return nil
}
