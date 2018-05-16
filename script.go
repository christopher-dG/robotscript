package robotscript

import (
	"log"

	"github.com/kylelemons/go-gypsy/yaml"
)

// Script is a sequence of commands to execute.
type Script struct {
	Commands []Command
}

// NewScript parses a YAML file into a new Script.
func NewScript(filename string) (*Script, error) {
	script := &Script{}

	file, err := yaml.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cmdsNode, err := yaml.Child(file.Root, "commands")
	if err != nil {
		return nil, err
	}

	cmdMaps, err := toList(cmdsNode)
	if err != nil {
		log.Fatal(err)
	}

	for _, cmdMap := range cmdMaps {
		cmdMap, err := toMap(cmdMap)
		if err != nil {
			return nil, err
		}

		cmdKey, err := getSingleKey(cmdMap)
		if err != nil {
			return nil, err
		}

		cmdOptions, err := toMap(cmdMap[cmdKey])
		if err != nil {
			return nil, err
		}

		command, err := NewCommand(cmdKey, cmdOptions)
		if err != nil {
			return nil, err
		}

		script.Commands = append(script.Commands, command)
	}

	return script, nil
}

// Execute executes each script command.
func (s *Script) Execute() {
	for _, cmd := range s.Commands {
		cmd.Execute()
	}
}
