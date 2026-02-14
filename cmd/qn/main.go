package main

import (
	"fmt"
	"os"
	"strings"

	"go-spass/quick-note/internal"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	// Handle help before checking MDNOTES_DIR
	if len(args) > 0 {
		switch args[0] {
		case "help", "--help", "-h":
			printUsage()
			return nil
		}
	}

	baseDir, err := internal.NotesDir()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		p := internal.NewPrompter(os.Stdin, os.Stderr)
		return internal.Create(p, baseDir)
	}

	switch args[0] {
	case "list":
		showAll := false
		for _, arg := range args[1:] {
			if arg == "--all" {
				showAll = true
			}
		}
		return internal.List(os.Stdout, baseDir, showAll)
	case "find":
		if len(args) < 2 {
			return fmt.Errorf("usage: qn find <topic>")
		}
		query := strings.Join(args[1:], " ")
		return internal.Find(os.Stdout, baseDir, query)
	default:
		return fmt.Errorf("unknown command: %s\nRun 'qn help' for usage", args[0])
	}
}

func printUsage() {
	fmt.Println(`qn - Quick Note CLI

Usage:
  qn              Create a new note interactively
  qn list         List 10 most recent notes
  qn list --all   List all notes
  qn find <topic> Search notes by topic
  qn help         Show this help message

Environment:
  MDNOTES_DIR     Path to the notes directory (required)
  EDITOR          Editor to open notes in (optional)`)
}
