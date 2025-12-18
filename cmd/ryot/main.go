package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"ryot/internal/commands"
)

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr, os.Stdin); err != nil {
		var ue commands.UsageError
		if errors.As(err, &ue) {
			fmt.Fprintln(os.Stderr, ue.Error())
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, usage())
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer, stdin io.Reader) error {
	if len(args) < 2 {
		return commands.UsageError{Msg: "missing subcommand"}
	}

	sub := args[1]
	subArgs := args[2:]

	switch sub {
	case "help", "-h", "--help":
		fmt.Fprintln(stdout, usage())
		return nil
	case "init-db":
		return commands.InitDB(subArgs, stdout)
	case "update-cache":
		return commands.UpdateCache(subArgs, stdout)
	case "write-tree":
		return commands.WriteTree(subArgs, stdout)
	case "commit-tree":
		return commands.CommitTree(subArgs, stdout, stdin)
	case "cat-file":
		return commands.CatFile(subArgs, stdout)
	case "read-tree":
		return commands.ReadTree(subArgs, stdout)
	case "show-diff":
		return commands.ShowDiff(subArgs, stdout)
	default:
		return commands.UsageError{Msg: fmt.Sprintf("unknown subcommand: %s", sub)}
	}
}

func usage() string {
	lines := []string{
		"ryot - minimal git-like content tracker (early-git inspired)",
		"",
		"Usage:",
		"  ryot <subcommand> [args...]",
		"",
		"Subcommands:",
		"  init-db                      Initialize .dircache and objects directory",
		"  update-cache <file>...        Add/update file(s) to index and objects (add)",
		"  write-tree                    Create a tree object from index",
		"  commit-tree <tree> [-p <c>]   Create a commit object from a tree (message via stdin)",
		"  cat-file <sha1>               Inflate object into a temp file and print its type",
		"  read-tree <sha1>              List entries in a tree object",
		"  show-diff                     Diff index (staged) vs working tree",
	}
	return strings.Join(lines, "\n")
}
