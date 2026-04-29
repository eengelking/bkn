package runner

import (
	"bytes"
	"runtime"
	"strings"
	"testing"

	"github.com/eengelking/bkn/internal/config"
)

func skipIfNotUnix(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("runner tests require sh; skipping on windows")
	}
}

func TestExecute_PositionalArgsForwarded(t *testing.T) {
	skipIfNotUnix(t)
	var out, errBuf bytes.Buffer
	Execute(&out, &errBuf, config.Command{Command: `echo "$1-$2"`}, []string{"a", "b"})
	if got := strings.TrimSpace(out.String()); got != "a-b" {
		t.Errorf("stdout = %q, want %q", got, "a-b")
	}
}

func TestExecute_AtExpansionCount(t *testing.T) {
	skipIfNotUnix(t)
	var out, errBuf bytes.Buffer
	Execute(&out, &errBuf, config.Command{Command: `echo $#`}, []string{"x", "y", "z"})
	if got := strings.TrimSpace(out.String()); got != "3" {
		t.Errorf("stdout = %q, want 3", got)
	}
}

func TestExecute_NoArgs(t *testing.T) {
	skipIfNotUnix(t)
	var out, errBuf bytes.Buffer
	Execute(&out, &errBuf, config.Command{Command: "echo hi"}, nil)
	if got := strings.TrimSpace(out.String()); got != "hi" {
		t.Errorf("stdout = %q, want hi", got)
	}
}

func TestExecute_ShellFailureDoesNotPanicOrExit(t *testing.T) {
	skipIfNotUnix(t)
	var out, errBuf bytes.Buffer
	// Documented quirk: a non-zero shell exit is reported on stdout but does
	// not propagate from Execute (caller continues normally).
	Execute(&out, &errBuf, config.Command{Command: "exit 7"}, nil)
	if !strings.Contains(out.String(), "Failed to execute command") {
		t.Errorf("expected failure message on stdout, got %q", out.String())
	}
}
