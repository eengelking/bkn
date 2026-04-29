package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/eengelking/bkn/internal/config"
)

const (
	Teal  = "\033[1;36m"
	White = "\033[0;37m"
	Reset = "\033[0m"
)

func ListCommands(w io.Writer, commands []config.Command) {
	fmt.Fprintf(w, "%s-----\nAvailable commands:\n-----%s\n\n", Teal, Reset)

	longestName := 0
	for _, cmd := range commands {
		if len(cmd.Name) > longestName {
			longestName = len(cmd.Name)
		}
	}

	for _, cmd := range commands {
		padding := longestName - len(cmd.Name) + 5
		spaces := strings.Repeat(" ", padding)
		fmt.Fprintf(w, "%s%s%s%s%s%s%s\n", Teal, cmd.Name, Reset, spaces, White, cmd.Description, Reset)
	}
}

func PrintUsage(w io.Writer, commands []config.Command) {
	fmt.Fprintf(w, "\n%sBKN (bacon): A better Make using Golang%s\n\n", Teal, Reset)
	fmt.Fprintf(w, "%sUsage:%s bkn [flags] [command] [args...]\n\n", Teal, Reset)
	fmt.Fprintf(w, "%sFlags:%s\n", Teal, Reset)
	fmt.Fprintf(w, "  %s-c, --config <path>%s   Path to YAML file (default \"bkn.yaml\")\n", Teal, Reset)
	fmt.Fprintf(w, "  %s-h, --help%s            Show this help and exit\n\n", Teal, Reset)
	ListCommands(w, commands)
}
