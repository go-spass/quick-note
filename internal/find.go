package internal

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type searchResult struct {
	Note  Note
	Score int
}

// Find searches for notes matching the given query and displays results.
func Find(w io.Writer, baseDir, query string) error {
	if query == "" {
		return fmt.Errorf("search query is required")
	}

	notes, err := ScanNotes(baseDir)
	if err != nil {
		return err
	}

	q := strings.ToLower(query)
	var results []searchResult

	for _, note := range notes {
		score := scoreNote(note, q)
		if score > 0 {
			results = append(results, searchResult{Note: note, Score: score})
		}
	}

	if len(results) == 0 {
		_, _ = fmt.Fprintf(w, "No notes found matching %q.\n", query)
		return nil
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	for _, r := range results {
		title := r.Note.Frontmatter.Title
		if title == "" {
			title = "(untitled)"
		}
		tags := ""
		if len(r.Note.Frontmatter.Tags) > 0 {
			tags = "  [" + strings.Join(r.Note.Frontmatter.Tags, ", ") + "]"
		}

		_, _ = fmt.Fprintf(w, "%s  (%s)%s\n", title, r.Note.Folder, tags)

		excerpt := findExcerpt(r.Note.Body, q)
		if excerpt != "" {
			_, _ = fmt.Fprintf(w, "  %s\n", excerpt)
		}
	}

	return nil
}

func scoreNote(note Note, query string) int {
	score := 0

	// Title match (highest priority)
	if strings.Contains(strings.ToLower(note.Frontmatter.Title), query) {
		score += 10
	}

	// Tag match
	for _, tag := range note.Frontmatter.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			score += 8
			break
		}
	}

	// Alias match
	for _, alias := range note.Frontmatter.Aliases {
		if strings.Contains(strings.ToLower(alias), query) {
			score += 8
			break
		}
	}

	// Body match
	if strings.Contains(strings.ToLower(note.Body), query) {
		score += 3
	}

	return score
}

func findExcerpt(body, query string) string {
	lower := strings.ToLower(body)
	idx := strings.Index(lower, query)
	if idx < 0 {
		return ""
	}

	// Extract a window around the match
	start := idx - 40
	if start < 0 {
		start = 0
	}
	end := idx + len(query) + 40
	if end > len(body) {
		end = len(body)
	}

	excerpt := body[start:end]
	// Clean up newlines
	excerpt = strings.ReplaceAll(excerpt, "\n", " ")
	excerpt = strings.TrimSpace(excerpt)

	prefix := ""
	suffix := ""
	if start > 0 {
		prefix = "..."
	}
	if end < len(body) {
		suffix = "..."
	}

	return prefix + excerpt + suffix
}
