package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/eengelking/bkn/internal/config"
	"github.com/eengelking/bkn/internal/runner"
	"github.com/eengelking/bkn/internal/ui"
)

func main() {
	var configPath string
	var showHelp bool

	flag.StringVar(&configPath, "c", "bkn.yaml", "path to YAML file")
	flag.StringVar(&configPath, "config", "bkn.yaml", "path to YAML file")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.BoolVar(&showHelp, "help", false, "show help")

	flag.Usage = func() {
		commands, _ := config.ParseYAML(configPath)
		ui.PrintUsage(os.Stdout, commands)
	}

	flag.Parse()

	commands, err := config.ParseYAML(configPath)

	args := flag.Args()
	if showHelp || len(args) == 0 {
		ui.PrintUsage(os.Stdout, commands)
		return
	}

	if err != nil {
		fmt.Printf("Failed to parse YAML file: %s\n", err)
		os.Exit(1)
	}

	option := args[0]
	cmdArgs := args[1:]

	for _, cmd := range commands {
		if cmd.Name == option {
			runner.Execute(os.Stdout, os.Stderr, cmd, cmdArgs)
			return
		}
	}

	fmt.Printf("Invalid option: %s\n", option)
}
