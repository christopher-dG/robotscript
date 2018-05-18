package robotscript

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Script is a sequence of commands to execute.
type Script struct {
	Commands []Command
}

// NewScript parses the YAML file at filename into a new Script.
func NewScript(filename string) (*Script, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "reading script file failed")
	}
	return newScript(contents)
}

// newScript parses the contents of a YAML file into a new Script.
func newScript(contents []byte) (*Script, error) {
	m := make(map[string][]map[string]map[string]interface{})
	err := yaml.Unmarshal(contents, &m)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling YAML failed")
	}

	commands, ok := m["commands"]
	if !ok {
		return nil, errors.New("script: missing 'commands' section")
	}

	script := &Script{}
	for _, cmd := range commands {
		key, err := getSingleKey(cmd)
		if err != nil {
			return nil, err
		}

		command, err := NewCommand(key, cmd[key])
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
