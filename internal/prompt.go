package internal

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Prompter handles interactive input. It reads from Reader and writes
// prompts to Writer, making it testable by injecting buffers.
type Prompter struct {
	Reader io.Reader
	Writer io.Writer
	scanner *bufio.Scanner
}

// NewPrompter creates a Prompter from the given reader and writer.
func NewPrompter(r io.Reader, w io.Writer) *Prompter {
	return &Prompter{
		Reader: r,
		Writer: w,
		scanner: bufio.NewScanner(r),
	}
}

// Ask prints a prompt and reads a single line of input.
func (p *Prompter) Ask(prompt string) string {
	_, _ = fmt.Fprint(p.Writer, prompt)
	if p.scanner.Scan() {
		return strings.TrimSpace(p.scanner.Text())
	}
	return ""
}

// AskRequired prints a prompt and re-asks until a non-empty answer is given.
func (p *Prompter) AskRequired(prompt string) string {
	for {
		answer := p.Ask(prompt)
		if answer != "" {
			return answer
		}
		_, _ = fmt.Fprintln(p.Writer, "This field is required.")
	}
}

// AskMenu displays a numbered menu and returns the selected option (0-indexed).
// The defaultIdx is used when the user presses Enter without input.
func (p *Prompter) AskMenu(prompt string, options []string, defaultIdx int) int {
	var b strings.Builder
	b.WriteString(prompt)
	for i, opt := range options {
		fmt.Fprintf(&b, " %d) %s ", i+1, opt)
	}
	fmt.Fprintf(&b, "[%d]: ", defaultIdx+1)

	answer := p.Ask(b.String())
	if answer == "" {
		return defaultIdx
	}

	// Parse the number
	var choice int
	_, err := fmt.Sscanf(answer, "%d", &choice)
	if err != nil || choice < 1 || choice > len(options) {
		_, _ = fmt.Fprintf(p.Writer, "Invalid choice, using default: %s\n", options[defaultIdx])
		return defaultIdx
	}
	return choice - 1
}

// AskYesNo prints a y/n prompt and returns true for yes.
func (p *Prompter) AskYesNo(prompt string) bool {
	answer := p.Ask(prompt + " (y/n): ")
	answer = strings.ToLower(answer)
	return answer == "y" || answer == "yes"
}
