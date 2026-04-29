package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestParseYAML_ValidSingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bkn.yaml")
	writeFile(t, path, `
commands:
  - name: hello
    description: say hi
    command: echo hi
  - name: bye
    description: say bye
    command: echo bye
`)
	cmds, err := ParseYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmds) != 2 {
		t.Fatalf("want 2 commands, got %d", len(cmds))
	}
	if cmds[0].Name != "hello" || cmds[0].Description != "say hi" || cmds[0].Command != "echo hi" {
		t.Errorf("cmd[0] = %+v", cmds[0])
	}
	if cmds[1].Name != "bye" {
		t.Errorf("cmd[1].Name = %q, want bye", cmds[1].Name)
	}
}

func TestParseYAML_MissingFile(t *testing.T) {
	_, err := ParseYAML(filepath.Join(t.TempDir(), "nope.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseYAML_Malformed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	writeFile(t, path, "commands: [this is : not valid")
	if _, err := ParseYAML(path); err == nil {
		t.Fatal("expected error for malformed YAML, got nil")
	}
}

func TestParseYAML_RecursiveInclude(t *testing.T) {
	dir := t.TempDir()
	child := filepath.Join(dir, "child.yaml")
	writeFile(t, child, `
commands:
  - name: from-child
    description: child cmd
    command: echo child
`)
	parent := filepath.Join(dir, "parent.yaml")
	writeFile(t, parent, `
commands:
  - name: from-parent
    description: parent cmd
    command: echo parent
include:
  - `+child+`
`)
	cmds, err := ParseYAML(parent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmds) != 2 {
		t.Fatalf("want 2 commands, got %d", len(cmds))
	}
	if cmds[0].Name != "from-parent" || cmds[1].Name != "from-child" {
		t.Errorf("merge order wrong: %q, %q", cmds[0].Name, cmds[1].Name)
	}
}

func TestParseYAML_IncludeMissing(t *testing.T) {
	dir := t.TempDir()
	parent := filepath.Join(dir, "parent.yaml")
	writeFile(t, parent, `
include:
  - `+filepath.Join(dir, "missing.yaml")+`
`)
	if _, err := ParseYAML(parent); err == nil {
		t.Fatal("expected error from missing include, got nil")
	}
}

func TestParseYAML_OSField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bkn.yaml")
	writeFile(t, path, `
commands:
  - name: cross
    description: dual-target
    command: echo cross
    os: [linux, darwin]
  - name: anywhere
    description: no restriction
    command: echo anywhere
`)
	cmds, err := ParseYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := cmds[0].OS; len(got) != 2 || got[0] != "linux" || got[1] != "darwin" {
		t.Errorf("cmds[0].OS = %v, want [linux darwin]", got)
	}
	if len(cmds[1].OS) != 0 {
		t.Errorf("cmds[1].OS = %v, want empty", cmds[1].OS)
	}
}

func TestCommand_AllowedOnHost(t *testing.T) {
	host := runtime.GOOS
	cases := []struct {
		name string
		os   []string
		want bool
	}{
		{"omitted", nil, true},
		{"empty", []string{}, true},
		{"matches host", []string{host}, true},
		{"only foreign", []string{"plan9"}, false},
		{"includes host", []string{"plan9", host}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := Command{OS: tc.os}
			if got := c.AllowedOnHost(); got != tc.want {
				t.Errorf("AllowedOnHost() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFilterForHost(t *testing.T) {
	host := runtime.GOOS
	in := []Command{
		{Name: "a"},
		{Name: "b", OS: []string{host}},
		{Name: "c", OS: []string{"plan9"}},
		{Name: "d", OS: []string{"plan9", host}},
	}
	out := FilterForHost(in)
	if len(out) != 3 {
		t.Fatalf("want 3 commands, got %d (%+v)", len(out), out)
	}
	for _, c := range out {
		if c.Name == "c" {
			t.Errorf("filtered slice should not contain plan9-only command")
		}
	}
}

func TestParseYAML_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.yaml")
	writeFile(t, path, "")
	cmds, err := ParseYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cmds) != 0 {
		t.Errorf("want 0 commands, got %d", len(cmds))
	}
}
