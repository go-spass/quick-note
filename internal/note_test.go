package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatterFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFM   Frontmatter
		wantBody string
	}{
		{
			name: "complete frontmatter",
			input: `---
title: "My Note"
date: 2026-02-13
tags: [go, cli]
status: draft
aliases: [my-note, mn]
---

# My Note

Some content here.
`,
			wantFM: Frontmatter{
				Title:   "My Note",
				Date:    "2026-02-13",
				Tags:    []string{"go", "cli"},
				Status:  "draft",
				Aliases: []string{"my-note", "mn"},
			},
			wantBody: "\n# My Note\n\nSome content here.\n",
		},
		{
			name: "empty tags and aliases",
			input: `---
title: "Empty"
date: 2026-01-01
tags: []
status: active
aliases: []
---

Body.
`,
			wantFM: Frontmatter{
				Title:   "Empty",
				Date:    "2026-01-01",
				Tags:    nil,
				Status:  "active",
				Aliases: nil,
			},
			wantBody: "\nBody.\n",
		},
		{
			name:     "no frontmatter",
			input:    "Just some text without frontmatter.\n",
			wantFM:   Frontmatter{},
			wantBody: "Just some text without frontmatter.\n",
		},
		{
			name: "frontmatter without closing",
			input: `---
title: "Broken"
tags: [broken]
some content after
`,
			wantFM:   Frontmatter{},
			wantBody: "---\ntitle: \"Broken\"\ntags: [broken]\nsome content after\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, err := ParseFrontmatterFromBytes([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fm.Title != tt.wantFM.Title {
				t.Errorf("Title = %q, want %q", fm.Title, tt.wantFM.Title)
			}
			if fm.Date != tt.wantFM.Date {
				t.Errorf("Date = %q, want %q", fm.Date, tt.wantFM.Date)
			}
			if fm.Status != tt.wantFM.Status {
				t.Errorf("Status = %q, want %q", fm.Status, tt.wantFM.Status)
			}
			if !sliceEqual(fm.Tags, tt.wantFM.Tags) {
				t.Errorf("Tags = %v, want %v", fm.Tags, tt.wantFM.Tags)
			}
			if !sliceEqual(fm.Aliases, tt.wantFM.Aliases) {
				t.Errorf("Aliases = %v, want %v", fm.Aliases, tt.wantFM.Aliases)
			}
			if body != tt.wantBody {
				t.Errorf("Body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}

func TestFormatFrontmatter(t *testing.T) {
	fm := Frontmatter{
		Title:   "Test Note",
		Date:    "2026-02-13",
		Tags:    []string{"go", "testing"},
		Status:  "draft",
		Aliases: nil,
	}

	got := FormatFrontmatter(fm)
	want := "---\ntitle: \"Test Note\"\ndate: 2026-02-13\ntags: [go, testing]\nstatus: draft\naliases: []\n---\n"

	if got != want {
		t.Errorf("FormatFrontmatter() =\n%s\nwant:\n%s", got, want)
	}
}

func TestScanNotes(t *testing.T) {
	dir := t.TempDir()

	// Create folder structure
	for _, folder := range []string{"Inbox", "Projects", "Areas", "Resources", "Archive"} {
		if err := os.MkdirAll(filepath.Join(dir, folder), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// Create test notes
	note1 := `---
title: "Note One"
date: 2026-02-13
tags: [test]
status: draft
aliases: []
---

Content one.
`
	note2 := `---
title: "Note Two"
date: 2026-02-12
tags: [test, go]
status: active
aliases: []
---

Content two.
`
	if err := os.WriteFile(filepath.Join(dir, "Inbox", "2026-02-13-note-one.md"), []byte(note1), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "Resources", "note-two.md"), []byte(note2), 0o644); err != nil {
		t.Fatal(err)
	}

	notes, err := ScanNotes(dir)
	if err != nil {
		t.Fatalf("ScanNotes() error: %v", err)
	}

	if len(notes) != 2 {
		t.Fatalf("ScanNotes() returned %d notes, want 2", len(notes))
	}

	// Check that folders are set correctly
	folderMap := make(map[string]string)
	for _, n := range notes {
		folderMap[n.Frontmatter.Title] = n.Folder
	}
	if folderMap["Note One"] != "Inbox" {
		t.Errorf("Note One folder = %q, want %q", folderMap["Note One"], "Inbox")
	}
	if folderMap["Note Two"] != "Resources" {
		t.Errorf("Note Two folder = %q, want %q", folderMap["Note Two"], "Resources")
	}
}

func TestNotesDir(t *testing.T) {
	// Test with unset env var
	t.Setenv("MDNOTES_DIR", "")
	_, err := NotesDir()
	if err == nil {
		t.Error("NotesDir() should error when MDNOTES_DIR is not set")
	}

	// Test with valid dir
	dir := t.TempDir()
	t.Setenv("MDNOTES_DIR", dir)
	got, err := NotesDir()
	if err != nil {
		t.Errorf("NotesDir() unexpected error: %v", err)
	}
	if got != dir {
		t.Errorf("NotesDir() = %q, want %q", got, dir)
	}

	// Test with non-existent path
	t.Setenv("MDNOTES_DIR", "/nonexistent/path")
	_, err = NotesDir()
	if err == nil {
		t.Error("NotesDir() should error for non-existent path")
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
