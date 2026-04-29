package runner

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/eengelking/bkn/internal/config"
)

func Execute(stdout, stderr io.Writer, command config.Command, args []string) {
	shArgs := append([]string{"-c", command.Command, "bkn"}, args...)

	cmd := exec.Command("sh", shArgs...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(stdout, "Failed to execute command: %s\n", err)
	}
}
