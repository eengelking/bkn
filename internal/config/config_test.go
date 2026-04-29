package config

import (
	"os"
	"path/filepath"
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
