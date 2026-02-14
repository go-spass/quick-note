package internal

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrompterAsk(t *testing.T) {
	input := "hello world\n"
	r := strings.NewReader(input)
	w := &bytes.Buffer{}
	p := NewPrompter(r, w)

	got := p.Ask("Enter: ")
	if got != "hello world" {
		t.Errorf("Ask() = %q, want %q", got, "hello world")
	}
	if !strings.Contains(w.String(), "Enter: ") {
		t.Errorf("prompt not written to output")
	}
}

func TestPrompterAskRequired(t *testing.T) {
	// First line empty, second line has content
	input := "\nactual answer\n"
	r := strings.NewReader(input)
	w := &bytes.Buffer{}
	p := NewPrompter(r, w)

	got := p.AskRequired("Name: ")
	if got != "actual answer" {
		t.Errorf("AskRequired() = %q, want %q", got, "actual answer")
	}
	if !strings.Contains(w.String(), "This field is required.") {
		t.Errorf("expected required message, got: %s", w.String())
	}
}

func TestPrompterAskMenu(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		options    []string
		defaultIdx int
		want       int
	}{
		{"select option 2", "2\n", []string{"A", "B", "C"}, 0, 1},
		{"default on empty", "\n", []string{"A", "B", "C"}, 0, 0},
		{"invalid falls to default", "99\n", []string{"A", "B"}, 1, 1},
		{"text falls to default", "abc\n", []string{"A", "B"}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			w := &bytes.Buffer{}
			p := NewPrompter(r, w)

			got := p.AskMenu("Pick:", tt.options, tt.defaultIdx)
			if got != tt.want {
				t.Errorf("AskMenu() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPrompterAskYesNo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"yes", "y\n", true},
		{"YES", "YES\n", true},
		{"no", "n\n", false},
		{"empty", "\n", false},
		{"random", "maybe\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			w := &bytes.Buffer{}
			p := NewPrompter(r, w)

			got := p.AskYesNo("Continue?")
			if got != tt.want {
				t.Errorf("AskYesNo() = %v, want %v", got, tt.want)
			}
		})
	}
}
