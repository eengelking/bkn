package config

import (
	"os"

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
