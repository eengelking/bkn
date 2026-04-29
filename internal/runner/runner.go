package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/eengelking/bkn/internal/config"
)

func Execute(stdout, stderr io.Writer, command config.Command, args []string) {
	var err error
	if runtime.GOOS == "windows" {
		err = runWindows(stdout, stderr, command.Command, args)
	} else {
		err = runPosix(stdout, stderr, command.Command, args)
	}
	if err != nil {
		fmt.Fprintf(stdout, "Failed to execute command: %s\n", err)
	}
}

func runPosix(stdout, stderr io.Writer, body string, args []string) error {
	shArgs := append([]string{"-c", body, "bkn"}, args...)
	cmd := exec.Command("sh", shArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

// runWindows writes the command body to a temp .bat file and invokes it via
// cmd.exe so that args bind to %1, %2, %* the way `sh -c` binds $1, $2, $@ on
// posix. cmd /C with an inline string does NOT perform %1 substitution — it
// only happens inside batch files — which is why we go through a temp file.
func runWindows(stdout, stderr io.Writer, body string, args []string) error {
	f, err := os.CreateTemp("", "bkn-*.bat")
	if err != nil {
		return err
	}
	path := f.Name()
	defer os.Remove(path)

	if _, err := f.WriteString("@echo off\r\n" + body + "\r\n"); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	cmdArgs := append([]string{"/C", filepath.Clean(path)}, args...)
	cmd := exec.Command("cmd.exe", cmdArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
