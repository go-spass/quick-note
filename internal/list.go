package internal

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// List displays recent notes sorted by modification time.
func List(w io.Writer, baseDir string, showAll bool) error {
	notes, err := ScanNotes(baseDir)
	if err != nil {
		return err
	}

	if len(notes) == 0 {
		_, _ = fmt.Fprintln(w, "No notes found.")
		return nil
	}

	// Sort by modification time, newest first
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].ModTime.After(notes[j].ModTime)
	})

	limit := 10
	if showAll || len(notes) <= limit {
		limit = len(notes)
	}

	for _, note := range notes[:limit] {
		title := note.Frontmatter.Title
		if title == "" {
			title = "(untitled)"
		}
		date := note.Frontmatter.Date
		if date == "" {
			date = note.ModTime.Format("2006-01-02")
		}
		tags := ""
		if len(note.Frontmatter.Tags) > 0 {
			tags = "  [" + strings.Join(note.Frontmatter.Tags, ", ") + "]"
		}
		_, _ = fmt.Fprintf(w, "%s  %s  (%s)%s\n", date, title, note.Folder, tags)
	}

	if !showAll && len(notes) > 10 {
		_, _ = fmt.Fprintf(w, "\nShowing 10 of %d notes. Use --all to show all.\n", len(notes))
	}

	return nil
}
