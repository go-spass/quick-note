package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Folders defines the PARA folder names and their indices.
var Folders = []string{"Inbox", "Projects", "Areas", "Resources"}

// Create runs the interactive note creation flow.
func Create(p *Prompter, baseDir string) error {
	title := p.AskRequired("Title: ")
	folderIdx := p.AskMenu("Folder:", Folders, 0)
	folder := Folders[folderIdx]

	tagsInput := p.Ask("Tags (comma-separated): ")
	tags := NormalizeTags(tagsInput)

	body := p.Ask("Body: ")

	var urls []string
	if folder == "Resources" {
		urlInput := p.Ask("URLs (comma-separated): ")
		if urlInput != "" {
			for _, u := range strings.Split(urlInput, ",") {
				u = strings.TrimSpace(u)
				if u != "" {
					urls = append(urls, u)
				}
			}
		}
	}

	// Select template
	templateName := "basic.md"
	if folder == "Projects" {
		templateName = "project.md"
		// Ensure "project" tag is present
		hasProject := false
		for _, t := range tags {
			if t == "project" {
				hasProject = true
				break
			}
		}
		if !hasProject {
			tags = append(tags, "project")
		}
	}

	// Generate filename
	slug := Slugify(title)
	today := time.Now().Format("2006-01-02")
	var filename string
	if folder == "Inbox" || folder == "Projects" {
		filename = fmt.Sprintf("%s-%s.md", today, slug)
	} else {
		filename = slug + ".md"
	}

	destDir := filepath.Join(baseDir, folder)
	destPath := filepath.Join(destDir, filename)

	// Check for duplicate
	if _, err := os.Stat(destPath); err == nil {
		_, _ = fmt.Fprintf(p.Writer, "Warning: file already exists: %s\n", destPath)
	}

	// Read and fill template
	content, err := buildNoteContent(baseDir, templateName, title, today, tags, body, urls)
	if err != nil {
		return err
	}

	// Ensure folder exists
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", destDir, err)
	}

	if err := os.WriteFile(destPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing note: %w", err)
	}

	_, _ = fmt.Fprintf(p.Writer, "Created: %s\n", destPath)

	// Offer to open in editor
	editor := os.Getenv("EDITOR")
	if editor != "" && p.AskYesNo("Open in editor?") {
		cmd := exec.Command(editor, destPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			_, _ = fmt.Fprintf(p.Writer, "Error opening editor: %v\n", err)
		}
	}

	return nil
}

func buildNoteContent(baseDir, templateName, title, date string, tags []string, body string, urls []string) (string, error) {
	tmpl, err := ReadTemplate(baseDir, templateName)
	if err != nil {
		// Fall back to generating content without a template
		return buildFallbackContent(title, date, tags, body, urls, templateName == "project.md"), nil
	}

	// Replace template placeholders
	content := strings.ReplaceAll(tmpl, "{{title}}", title)
	content = strings.ReplaceAll(content, "{{date}}", date)

	// Replace the frontmatter tags and status
	fm := Frontmatter{
		Title:   title,
		Date:    date,
		Tags:    tags,
		Status:  "draft",
		Aliases: nil,
	}
	if templateName == "project.md" {
		fm.Status = "active"
	}

	// Build fresh frontmatter and replace the template's
	freshFM := FormatFrontmatter(fm)

	// Find and replace the existing frontmatter block
	lines := strings.Split(content, "\n")
	inFM := false
	fmEnd := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if !inFM {
				inFM = true
			} else {
				fmEnd = i
				break
			}
		}
	}

	if fmEnd > 0 {
		remaining := strings.Join(lines[fmEnd+1:], "\n")
		content = freshFM + remaining
	}

	// Add body if provided
	if body != "" {
		content = addBodyToContent(content, body)
	}

	// Add URLs for Resources
	if len(urls) > 0 {
		content = addURLsToContent(content, urls)
	}

	return content, nil
}

func buildFallbackContent(title, date string, tags []string, body string, urls []string, isProject bool) string {
	fm := Frontmatter{
		Title:   title,
		Date:    date,
		Tags:    tags,
		Status:  "draft",
		Aliases: nil,
	}
	if isProject {
		fm.Status = "active"
	}

	var b strings.Builder
	b.WriteString(FormatFrontmatter(fm))
	b.WriteString("\n# " + title + "\n")

	if isProject {
		b.WriteString("\n## Goal\n\n")
		if body != "" {
			b.WriteString(body + "\n")
		}
		b.WriteString("\n## Tasks\n\n- [ ]\n")
		b.WriteString("\n## Notes\n\n")
		b.WriteString("\n## Related\n\n")
		b.WriteString("\n## Log\n\n")
		b.WriteString(fmt.Sprintf("- %s — Project created.\n", date))
	} else {
		b.WriteString("\n## Notes\n\n")
		if body != "" {
			b.WriteString(body + "\n")
		}
		b.WriteString("\n## References\n")
		if len(urls) > 0 {
			b.WriteString("\n")
			for _, u := range urls {
				b.WriteString("- " + u + "\n")
			}
		}
	}

	return b.String()
}

func addBodyToContent(content, body string) string {
	// Insert body after the first "## Notes" heading
	marker := "## Notes"
	idx := strings.Index(content, marker)
	if idx < 0 {
		return content + "\n" + body + "\n"
	}
	insertAt := idx + len(marker)
	// Skip any newlines after the heading
	for insertAt < len(content) && content[insertAt] == '\n' {
		insertAt++
	}
	return content[:insertAt] + body + "\n\n" + content[insertAt:]
}

func addURLsToContent(content string, urls []string) string {
	marker := "## References"
	idx := strings.Index(content, marker)
	if idx < 0 {
		// Append a References section
		content += "\n## References\n\n"
		for _, u := range urls {
			content += "- " + u + "\n"
		}
		return content
	}
	insertAt := idx + len(marker)
	for insertAt < len(content) && content[insertAt] == '\n' {
		insertAt++
	}
	var urlBlock strings.Builder
	for _, u := range urls {
		urlBlock.WriteString("- " + u + "\n")
	}
	urlBlock.WriteString("\n")
	return content[:insertAt] + urlBlock.String() + content[insertAt:]
}
