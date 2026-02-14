package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func setupListDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, folder := range []string{"Inbox", "Projects", "Areas", "Resources", "Archive"} {
		if err := os.MkdirAll(filepath.Join(dir, folder), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestListEmpty(t *testing.T) {
	dir := setupListDir(t)

	var buf bytes.Buffer
	err := List(&buf, dir, false)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if buf.String() != "No notes found.\n" {
		t.Errorf("expected 'No notes found.', got: %s", buf.String())
	}
}

func TestListNotes(t *testing.T) {
	dir := setupListDir(t)

	for i := range 3 {
		note := fmt.Sprintf(`---
title: "Note %c"
date: 2026-02-1%d
tags: [test]
status: draft
aliases: []
---

Content.
`, rune('A'+i), i)
		path := filepath.Join(dir, "Inbox", fmt.Sprintf("2026-02-1%d-note-%c.md", i, rune('a'+i)))
		if err := os.WriteFile(path, []byte(note), 0o644); err != nil {
			t.Fatal(err)
		}
		modTime := time.Now().Add(time.Duration(-i) * time.Hour)
		if err := os.Chtimes(path, modTime, modTime); err != nil {
			t.Fatal(err)
		}
	}

	var buf bytes.Buffer
	err := List(&buf, dir, false)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	output := buf.String()
	lines := nonEmptyLines(output)
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d: %s", len(lines), output)
	}
}

func TestListLimit(t *testing.T) {
	dir := setupListDir(t)

	for i := range 15 {
		note := fmt.Sprintf(`---
title: "Note %d"
date: 2026-02-13
tags: []
status: draft
aliases: []
---
`, i)
		path := filepath.Join(dir, "Inbox", fmt.Sprintf("2026-02-13-note-%02d.md", i))
		if err := os.WriteFile(path, []byte(note), 0o644); err != nil {
			t.Fatal(err)
		}
		modTime := time.Now().Add(time.Duration(-i) * time.Minute)
		if err := os.Chtimes(path, modTime, modTime); err != nil {
			t.Fatal(err)
		}
	}

	// Without --all, should show 10 + footer
	var buf bytes.Buffer
	err := List(&buf, dir, false)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if !strings.Contains(buf.String(), "Showing 10 of 15") {
		t.Errorf("expected limit message, got: %s", buf.String())
	}

	// With --all, should show all 15
	buf.Reset()
	err = List(&buf, dir, true)
	if err != nil {
		t.Fatalf("List(--all) error: %v", err)
	}
	if strings.Contains(buf.String(), "Showing") {
		t.Errorf("should not have limit message with --all, got: %s", buf.String())
	}
}

func nonEmptyLines(s string) []string {
	var result []string
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}
