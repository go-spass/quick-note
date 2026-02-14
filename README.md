# Quick Note CLI

A lightweight CLI (`qn`) for quickly capturing notes into a plain-file second brain.

## Install

```bash
go install go-spass/quick-note/cmd/qn@latest
```

Or build from source:

```bash
make build
```

## Configuration

Set the `MDNOTES_DIR` environment variable to your notes directory:

```bash
export MDNOTES_DIR=~/mdnotes
```

## Usage

### Create a Note

Run `qn` with no arguments to create a note interactively:

```
$ qn
Title: API Design Patterns
Folder: 1) Inbox  2) Projects  3) Areas  4) Resources  [1]: 4
Tags (comma-separated): golang, api, patterns
Body: Notes on common API design patterns in Go
URLs (comma-separated): https://example.com/api-patterns
Open in editor? (y/n): y
Created: ~/mdnotes/Resources/api-design-patterns.md
```

Prompts:

1. **Title** — required
2. **Folder** — numbered menu, defaults to Inbox
3. **Tags** — comma-separated, optional
4. **Body** — single-line description, optional
5. **URLs** — only for Resources; added to `## References` section
6. **Open in editor?** — opens in `$EDITOR` if yes

### List Recent Notes

```bash
qn list          # Show 10 most recent notes
qn list --all    # Show all notes
```

Output format: `2026-02-13  Title  (Folder)  [tag1, tag2]`

### Find Notes

```bash
qn find <topic>
```

Searches title, tags, and aliases first (ranked higher), then body content.

## How It Works

- **Templates** are auto-selected by folder: Projects use `_templates/project.md`, everything else uses `_templates/basic.md`
- **Filenames** are date-prefixed in Inbox and Projects (`2026-02-13-topic.md`), slug-only in Areas and Resources (`topic.md`)
- **Tags** are normalized to lowercase, hyphen-separated
- **Duplicate filenames** produce a warning but don't block creation

## Project Layout

```
cmd/
  qn/
    main.go           # Entry point, arg parsing
internal/
  create.go           # Note creation logic
  list.go             # List subcommand
  find.go             # Find/search subcommand
  prompt.go           # Interactive prompt helpers
  note.go             # Note/frontmatter types and parsing
  slug.go             # Filename slugification
```

## Development

```bash
make build     # Build binary to ./bin/qn
make test      # Run tests
make install   # Install to $GOPATH/bin
make clean     # Remove build artifacts
```

## Technical Details

- **Language:** Go 1.25
- **Dependencies:** Standard library only
- **Testing:** `go test`, table-driven
