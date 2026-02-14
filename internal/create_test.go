package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestNotesDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	for _, folder := range []string{"Inbox", "Projects", "Areas", "Resources", "Archive"} {
		if err := os.MkdirAll(filepath.Join(dir, folder), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// Create templates
	tmplDir := filepath.Join(dir, "_templates")
	if err := os.MkdirAll(tmplDir, 0o755); err != nil {
		t.Fatal(err)
	}

	basicTmpl := `---
title: "{{title}}"
date: {{date}}
tags: []
status: draft
aliases: []
---

# {{title}}

## Notes

## References
`
	projectTmpl := `---
title: "{{title}}"
date: {{date}}
tags: [project]
status: active
aliases: []
---

# {{title}}

## Goal

## Tasks

- [ ]

## Notes

## Related

## Log

- {{date}} — Project created.
`
	if err := os.WriteFile(filepath.Join(tmplDir, "basic.md"), []byte(basicTmpl), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmplDir, "project.md"), []byte(projectTmpl), 0o644); err != nil {
		t.Fatal(err)
	}

	return dir
}

func TestCreateInbox(t *testing.T) {
	dir := setupTestNotesDir(t)

	// Simulate: title="Test Note", folder=1 (Inbox), tags="go, testing", body="A test note", no editor
	input := "Test Note\n1\ngo, testing\nA test note\nn\n"
	r := strings.NewReader(input)
	w := &bytes.Buffer{}
	p := NewPrompter(r, w)

	// Unset EDITOR to avoid open prompt issues
	t.Setenv("EDITOR", "")

	err := Create(p, dir)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	output := w.String()
	if !strings.Contains(output, "Created:") {
		t.Errorf("expected 'Created:' in output, got: %s", output)
	}

	// Find the created file
	entries, _ := os.ReadDir(filepath.Join(dir, "Inbox"))
	if len(entries) != 1 {
		t.Fatalf("expected 1 file in Inbox, got %d", len(entries))
	}

	filename := entries[0].Name()
	if !strings.HasSuffix(filename, "-test-note.md") {
		t.Errorf("filename %q should end with -test-note.md", filename)
	}

	// Check content
	data, _ := os.ReadFile(filepath.Join(dir, "Inbox", filename))
	content := string(data)
	if !strings.Contains(content, `title: "Test Note"`) {
		t.Errorf("expected title in frontmatter, got:\n%s", content)
	}
	if !strings.Contains(content, "go, testing") {
		t.Errorf("expected tags in frontmatter, got:\n%s", content)
	}
}

func TestCreateProject(t *testing.T) {
	dir := setupTestNotesDir(t)

	// folder=2 (Projects), no tags input, no body
	input := "My Project\n2\n\n\nn\n"
	r := strings.NewReader(input)
	w := &bytes.Buffer{}
	p := NewPrompter(r, w)
	t.Setenv("EDITOR", "")

	err := Create(p, dir)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	entries, _ := os.ReadDir(filepath.Join(dir, "Projects"))
	if len(entries) != 1 {
		t.Fatalf("expected 1 file in Projects, got %d", len(entries))
	}

	data, _ := os.ReadFile(filepath.Join(dir, "Projects", entries[0].Name()))
	content := string(data)

	// Should have project tag auto-added
	if !strings.Contains(content, "project") {
		t.Errorf("expected 'project' tag, got:\n%s", content)
	}
	// Should use project template sections
	if !strings.Contains(content, "## Goal") {
		t.Errorf("expected '## Goal' section, got:\n%s", content)
	}
}

func TestCreateResource(t *testing.T) {
	dir := setupTestNotesDir(t)

	// folder=4 (Resources), with URLs
	input := "API Reference\n4\napi, reference\nUseful API docs\nhttps://example.com/api\nn\n"
	r := strings.NewReader(input)
	w := &bytes.Buffer{}
	p := NewPrompter(r, w)
	t.Setenv("EDITOR", "")

	err := Create(p, dir)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	entries, _ := os.ReadDir(filepath.Join(dir, "Resources"))
	if len(entries) != 1 {
		t.Fatalf("expected 1 file in Resources, got %d", len(entries))
	}

	filename := entries[0].Name()
	// Resources use slug-only filenames
	if strings.Contains(filename, "2026") {
		t.Errorf("Resources filename should not have date prefix: %s", filename)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "Resources", filename))
	content := string(data)
	if !strings.Contains(content, "https://example.com/api") {
		t.Errorf("expected URL in content, got:\n%s", content)
	}
}

func TestBuildNoteContentFallback(t *testing.T) {
	// Test with a non-existent template directory
	dir := t.TempDir()
	content, err := buildNoteContent(dir, "basic.md", "Fallback Test", "2026-02-13", []string{"test"}, "Some body", nil)
	if err != nil {
		t.Fatalf("buildNoteContent() error: %v", err)
	}
	if !strings.Contains(content, `title: "Fallback Test"`) {
		t.Errorf("expected title in fallback content, got:\n%s", content)
	}
	if !strings.Contains(content, "# Fallback Test") {
		t.Errorf("expected heading in fallback content, got:\n%s", content)
	}
}
