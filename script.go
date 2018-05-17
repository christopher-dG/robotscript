package robotscript

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "reading script file failed")
	}

	cmdsNode, err := yaml.Child(file.Root, "commands")
	if err != nil {
		return nil, errors.Wrap(err, "commands section not found")
	}

	cmdMaps := unYAML(cmdsNode)
	if !isList(cmdMaps) {
		return nil, wrongOptType("script", "list", "commands", cmdMaps)
	}

	for _, cmdMap := range cmdMaps.([]interface{}) {
		if !isMap(cmdMap) {
			return nil, wrongListEntryType("script", "map", "commands", cmdMap)
		}
		cmdMap := cmdMap.(map[string]interface{})

		cmdKey, err := getSingleKey(cmdMap)
		if err != nil {
			return nil, err
		}

		command, err := NewCommand(cmdKey, cmdMap[cmdKey].(map[string]interface{}))
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
