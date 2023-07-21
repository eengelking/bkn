package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "gopkg.in/yaml.v2"
)

type Command struct {
    Name        string
    Description string
    Command     string
}

type Include struct {
	Commands []Command `yaml:"commands"`
	Include  []string  `yaml:"include"`
}

func parseYaml(yamlFile string) ([]Command, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var includes Include
	err = yaml.Unmarshal(data, &includes)
	if err != nil {
		return nil, err
	}

	commands := includes.Commands

	for _, includeFile := range includes.Include {
		includeCommands, err := parseYaml(includeFile) // Treat included files the same as the main file
		if err != nil {
			return nil, err
		}

		commands = append(commands, includeCommands...)
	}

	return commands, nil
}

func listCommands(commands []Command) {
	teal := "\033[1;36m"
	white := "\033[0;37m"
	reset := "\033[0m"

    fmt.Printf("\nBKN (bacon): A better Make using Golang\n\n")
    fmt.Printf("-----\nAvailable commands:\n-----\n\n")

	longestName := 0
	for _, cmd := range commands {
		if len(cmd.Name) > longestName {
			longestName = len(cmd.Name)
		}
	}

	for _, cmd := range commands {
		padding := longestName - len(cmd.Name) + 5
		spaces := strings.Repeat(" ", padding)

		fmt.Printf("%s%s%s%s%s%s\n", teal, cmd.Name, reset, spaces, white, cmd.Description)
	}
}

func executeCommand(command Command, variable string) {
    cmdStr := strings.Replace(command.Command, "<variable>", variable, -1)

    cmd := exec.Command("sh", "-c", cmdStr)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Failed to execute command: %s\n", err)
    }
}

func main() {
    commands, err := parseYaml("bkn.yaml")
    if err != nil {
        fmt.Printf("Failed to parse YAML file: %s\n", err)
        os.Exit(1)
    }

    if len(os.Args) == 1 || os.Args[1] == "--help" {
        listCommands(commands)
        return
    }

    option := os.Args[1]
    variable := ""

    if len(os.Args) > 2 {
        variable = os.Args[2]
    }

    for _, cmd := range commands {
        if cmd.Name == option {
            executeCommand(cmd, variable)
            return
        }
    }

    fmt.Printf("Invalid option: %s\n", option)
}