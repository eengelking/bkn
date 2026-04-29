package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/eengelking/bkn/internal/config"
)

func sampleCommands() []config.Command {
	return []config.Command{
		{Name: "a", Description: "first"},
		{Name: "longest-name", Description: "second"},
		{Name: "mid", Description: "third"},
	}
}

func TestListCommands_ContainsNamesAndDescriptions(t *testing.T) {
	var buf bytes.Buffer
	ListCommands(&buf, sampleCommands())
	out := buf.String()
	for _, want := range []string{"a", "longest-name", "mid", "first", "second", "third", "Available commands"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n--- output ---\n%s", want, out)
		}
	}
}

func TestListCommands_PaddingAlignsColumns(t *testing.T) {
	var buf bytes.Buffer
	ListCommands(&buf, sampleCommands())

	lines := strings.Split(buf.String(), "\n")
	cols := make([]int, 0, 3)
	for _, line := range lines {
		idx := strings.Index(line, "first")
		if idx < 0 {
			idx = strings.Index(line, "second")
		}
		if idx < 0 {
			idx = strings.Index(line, "third")
		}
		if idx >= 0 {
			cols = append(cols, idx)
		}
	}
	if len(cols) != 3 {
		t.Fatalf("expected 3 description lines, got %d", len(cols))
	}
	for i := 1; i < len(cols); i++ {
		if cols[i] != cols[0] {
			t.Errorf("description column not aligned: %v", cols)
		}
	}
}

func TestListCommands_EmptyOmitsHeader(t *testing.T) {
	var buf bytes.Buffer
	ListCommands(&buf, nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty commands, got: %q", buf.String())
	}
}

func TestPrintUsage_EmptyCommandsStillPrintsFlags(t *testing.T) {
	var buf bytes.Buffer
	PrintUsage(&buf, nil)
	out := buf.String()
	for _, want := range []string{"Usage:", "-h, --help", "BKN"} {
		if !strings.Contains(out, want) {
			t.Errorf("usage output missing %q", want)
		}
	}
	if strings.Contains(out, "Available commands") {
		t.Errorf("did not expect 'Available commands' header with no commands; got:\n%s", out)
	}
}

func TestPrintUsage_IncludesFlagsAndCommands(t *testing.T) {
	var buf bytes.Buffer
	PrintUsage(&buf, sampleCommands())
	out := buf.String()
	for _, want := range []string{"Usage:", "-c, --config", "-h, --help", "BKN", "longest-name"} {
		if !strings.Contains(out, want) {
			t.Errorf("usage output missing %q", want)
		}
	}
}
