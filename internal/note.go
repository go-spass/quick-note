package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Frontmatter holds the YAML frontmatter fields for a note.
type Frontmatter struct {
	Title   string
	Date    string
	Tags    []string
	Status  string
	Aliases []string
}

// Note represents a parsed note file.
type Note struct {
	Frontmatter Frontmatter
	Body        string
	FilePath    string
	Folder      string
	ModTime     time.Time
}

// ParseFrontmatter parses YAML frontmatter from a note file.
// It expects the frontmatter to be delimited by "---" lines.
func ParseFrontmatter(path string) (Frontmatter, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Frontmatter{}, "", err
	}
	return ParseFrontmatterFromBytes(data)
}

// ParseFrontmatterFromBytes parses YAML frontmatter from raw bytes.
func ParseFrontmatterFromBytes(data []byte) (Frontmatter, string, error) {
	content := string(data)
	lines := strings.Split(content, "\n")

	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" {
		return Frontmatter{}, content, nil
	}

	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			endIdx = i
			break
		}
	}
	if endIdx == -1 {
		return Frontmatter{}, content, nil
	}

	fm := Frontmatter{}
	for _, line := range lines[1:endIdx] {
		key, value, ok := parseYAMLLine(line)
		if !ok {
			continue
		}
		switch key {
		case "title":
			fm.Title = unquote(value)
		case "date":
			fm.Date = value
		case "tags":
			fm.Tags = parseYAMLList(value)
		case "status":
			fm.Status = value
		case "aliases":
			fm.Aliases = parseYAMLList(value)
		}
	}

	body := ""
	if endIdx+1 < len(lines) {
		body = strings.Join(lines[endIdx+1:], "\n")
	}

	return fm, body, nil
}

// FormatFrontmatter renders a Frontmatter to YAML string.
func FormatFrontmatter(fm Frontmatter) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("title: %q\n", fm.Title))
	b.WriteString(fmt.Sprintf("date: %s\n", fm.Date))
	b.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(fm.Tags, ", ")))
	b.WriteString(fmt.Sprintf("status: %s\n", fm.Status))
	b.WriteString(fmt.Sprintf("aliases: [%s]\n", strings.Join(fm.Aliases, ", ")))
	b.WriteString("---\n")
	return b.String()
}

// ScanNotes reads all notes from the standard PARA folders under baseDir.
func ScanNotes(baseDir string) ([]Note, error) {
	folders := []string{"Inbox", "Projects", "Areas", "Resources", "Archive"}
	var notes []Note

	for _, folder := range folders {
		dir := filepath.Join(baseDir, folder)
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("reading %s: %w", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			fm, body, err := ParseFrontmatter(path)
			if err != nil {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			notes = append(notes, Note{
				Frontmatter: fm,
				Body:        body,
				FilePath:    path,
				Folder:      folder,
				ModTime:     info.ModTime(),
			})
		}
	}

	return notes, nil
}

// parseYAMLLine splits a simple "key: value" YAML line.
func parseYAMLLine(line string) (string, string, bool) {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])
	return key, value, true
}

// parseYAMLList parses a YAML inline list like "[a, b, c]".
func parseYAMLList(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// unquote removes surrounding double quotes from a string.
func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// ReadTemplate reads a template file from the notes directory.
func ReadTemplate(baseDir, name string) (string, error) {
	path := filepath.Join(baseDir, "_templates", name)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading template %s: %w", name, err)
	}
	return string(data), nil
}

// NotesDir returns the notes directory from the MDNOTES_DIR environment variable.
func NotesDir() (string, error) {
	dir := os.Getenv("MDNOTES_DIR")
	if dir == "" {
		return "", fmt.Errorf("MDNOTES_DIR environment variable is not set")
	}
	info, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("MDNOTES_DIR path %q: %w", dir, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("MDNOTES_DIR path %q is not a directory", dir)
	}
	return dir, nil
}

// FileScanner is a helper for reading lines from a file.
func FileScanner(path string) (*bufio.Scanner, *os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return bufio.NewScanner(f), f, nil
}
