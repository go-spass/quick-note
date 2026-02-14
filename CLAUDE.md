# Quick Note CLI — Claude Instructions

## Overview

`qn` is a Go CLI tool for creating, listing, and searching Markdown notes in a second brain system. The notes live in a separate directory specified by `$MDNOTES_DIR` (typically `~/mdnotes`).

## Project Structure

```
cmd/qn/main.go       — Entry point, arg parsing, command dispatch
internal/
  create.go           — Note creation: prompts, template reading, file writing
  list.go             — List subcommand: scan folders, parse frontmatter, display
  find.go             — Find subcommand: search frontmatter + body, rank results
  prompt.go           — Interactive prompt helpers (stdin reading, menus, y/n)
  note.go             — Note and frontmatter types, YAML parsing
  slug.go             — Filename slugification
```

## Build & Test

```bash
make build     # Build to ./bin/qn
make test      # Run all tests
make install   # go install
make clean     # Remove artifacts
```

## Key Design Decisions

- **No third-party dependencies.** Everything uses the Go standard library.
- **No YAML library.** Frontmatter is a fixed schema — parse it line-by-line instead of pulling in `gopkg.in/yaml.v3`.
- **`$MDNOTES_DIR` is required.** The CLI errors immediately if the env var is not set. No fallback, no config file.
- **Templates live in the notes directory** at `$MDNOTES_DIR/_templates/`. The CLI reads them at runtime, not embedded.

## Notes Directory Structure

The CLI expects this layout in `$MDNOTES_DIR`:

```
Inbox/          — Quick captures, date-prefixed filenames
Projects/       — Active projects, date-prefixed filenames
Areas/          — Ongoing responsibilities, slug-only filenames
Resources/      — Reference material, slug-only filenames
Archive/        — Completed items
_templates/
  basic.md      — Default template (used for Inbox, Areas, Resources)
  project.md    — Project template (used for Projects)
```

## Frontmatter Schema

Every note has this YAML frontmatter:

```yaml
---
title: "Note Title"
date: 2026-02-13
tags: [tag1, tag2]
status: draft
aliases: []
---
```

- `status` values: `draft`, `active`, `archived`
- Tags: flat, lowercase, hyphen-separated

## Filename Rules

- **Inbox & Projects:** `YYYY-MM-DD-slugified-title.md`
- **Areas & Resources:** `slugified-title.md`
- Slugification: lowercase, spaces/special chars become hyphens, collapse consecutive hyphens

## Commands

### `qn` (no args) — Create

Interactive prompts: title → folder → tags → body → URLs (Resources only) → open in editor?

### `qn list` — List

Show 10 most recent notes by modification time. `--all` flag shows everything. Format: `date  title  (folder)  [tags]`.

### `qn find <topic>` — Search

Search frontmatter (title, tags, aliases) first, then body. Rank frontmatter matches higher. Display title, folder, and excerpt.

## Conventions

- Use `fmt.Fprintf(os.Stderr, ...)` for prompts, `fmt.Println` for output
- Exit code 0 on success, 1 on errors
- Table-driven tests in `*_test.go` files alongside implementation
