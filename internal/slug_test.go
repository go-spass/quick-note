package internal

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "Hello World", "hello-world"},
		{"special chars", "API Design & Patterns!", "api-design-patterns"},
		{"multiple spaces", "too   many   spaces", "too-many-spaces"},
		{"leading trailing", "  leading and trailing  ", "leading-and-trailing"},
		{"numbers", "Go 1.25 Release", "go-1-25-release"},
		{"already slug", "my-cool-note", "my-cool-note"},
		{"mixed case", "MyFirstNote", "myfirstnote"},
		{"unicode", "café latte", "caf-latte"},
		{"empty", "", ""},
		{"only special", "!!!", ""},
		{"hyphens collapse", "a---b---c", "a-b-c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"simple", "go, cli, tooling", []string{"go", "cli", "tooling"}},
		{"mixed case", "Go, CLI", []string{"go", "cli"}},
		{"extra spaces", "  go ,  cli  ", []string{"go", "cli"}},
		{"hyphenated", "machine learning, deep-learning", []string{"machine-learning", "deep-learning"}},
		{"empty", "", nil},
		{"whitespace only", "   ", nil},
		{"single tag", "golang", []string{"golang"}},
		{"special chars removed", "c++, c#", []string{"c", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeTags(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("NormalizeTags(%q) = %v, want %v", tt.input, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("NormalizeTags(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}
