package config

import (
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Name        string
	Description string
	Command     string
	OS          []string `yaml:"os"`
}

func (c Command) AllowedOnHost() bool {
	if len(c.OS) == 0 {
		return true
	}
	for _, o := range c.OS {
		if o == runtime.GOOS {
			return true
		}
	}
	return false
}

func FilterForHost(commands []Command) []Command {
	out := make([]Command, 0, len(commands))
	for _, c := range commands {
		if c.AllowedOnHost() {
			out = append(out, c)
		}
	}
	return out
}

type Include struct {
	Commands []Command `yaml:"commands"`
	Include  []string  `yaml:"include"`
}

func ParseYAML(yamlFile string) ([]Command, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var includes Include
	if err := yaml.Unmarshal(data, &includes); err != nil {
		return nil, err
	}

	commands := includes.Commands

	for _, includeFile := range includes.Include {
		includeCommands, err := ParseYAML(includeFile)
		if err != nil {
			return nil, err
		}
		commands = append(commands, includeCommands...)
	}

	return commands, nil
}
