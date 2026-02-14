package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindByTitle(t *testing.T) {
	dir := setupFindTestDir(t)

	var buf bytes.Buffer
	err := Find(&buf, dir, "Golang")
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Golang Tips") {
		t.Errorf("expected to find 'Golang Tips', got: %s", output)
	}
}

func TestFindByTag(t *testing.T) {
	dir := setupFindTestDir(t)

	var buf bytes.Buffer
	err := Find(&buf, dir, "python")
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Python Basics") {
		t.Errorf("expected to find 'Python Basics', got: %s", output)
	}
}

func TestFindByContent(t *testing.T) {
	dir := setupFindTestDir(t)

	var buf bytes.Buffer
	err := Find(&buf, dir, "concurrency")
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Golang Tips") {
		t.Errorf("expected content match for 'Golang Tips', got: %s", output)
	}
}

func TestFindNoResults(t *testing.T) {
	dir := setupFindTestDir(t)

	var buf bytes.Buffer
	err := Find(&buf, dir, "nonexistent-topic-xyz")
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No notes found") {
		t.Errorf("expected 'No notes found', got: %s", output)
	}
}

func TestFindEmptyQuery(t *testing.T) {
	dir := setupFindTestDir(t)

	var buf bytes.Buffer
	err := Find(&buf, dir, "")
	if err == nil {
		t.Error("Find() should error on empty query")
	}
}

func TestFindRanking(t *testing.T) {
	dir := setupFindTestDir(t)

	var buf bytes.Buffer
	err := Find(&buf, dir, "go")
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}

	output := buf.String()
	// "Golang Tips" has "go" in tags, should appear
	if !strings.Contains(output, "Golang Tips") {
		t.Errorf("expected 'Golang Tips' in results, got: %s", output)
	}
}

func setupFindTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	for _, folder := range []string{"Inbox", "Projects", "Areas", "Resources", "Archive"} {
		if err := os.MkdirAll(filepath.Join(dir, folder), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	note1 := `---
title: "Golang Tips"
date: 2026-02-13
tags: [go, programming]
status: active
aliases: [go-tips]
---

# Golang Tips

Some tips about Go concurrency and channels.
`
	note2 := `---
title: "Python Basics"
date: 2026-02-12
tags: [python, programming]
status: draft
aliases: []
---

# Python Basics

Learning Python fundamentals.
`
	note3 := `---
title: "Cooking Recipes"
date: 2026-02-11
tags: [cooking, personal]
status: active
aliases: []
---

# Cooking Recipes

My favorite pasta recipe.
`

	for _, w := range []struct {
		path string
		data []byte
	}{
		{filepath.Join(dir, "Resources", "golang-tips.md"), []byte(note1)},
		{filepath.Join(dir, "Inbox", "2026-02-12-python-basics.md"), []byte(note2)},
		{filepath.Join(dir, "Areas", "cooking-recipes.md"), []byte(note3)},
	} {
		if err := os.WriteFile(w.path, w.data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	return dir
}
