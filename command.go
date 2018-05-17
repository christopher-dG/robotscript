package robotscript

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
)

type Command interface {
	Execute() // Execute executes the command.
	// New_____Command(options yaml.Map) (*_____Command, error)
}

// NewCommand creates a new command from a name options map.
func NewCommand(name string, options map[string]interface{}) (Command, error) {
	switch canonicalize(name) {
	case "click":
		return NewMouseClickCommand(options)
	case "exec":
		return NewExecCommand(options)
	case "keypress":
		return NewKeyPressCommand(options)
	case "mouse":
		return NewMouseMoveCommand(options)
	case "sleep":
		return NewSleepCommand(options)
	case "type":
		return NewTypeCommand(options)
	default:
		return nil, unrecognizedCmd(name)
	}
}

// Mouse commands

// MouseMoveCommand moves the mouse.
type MouseMoveCommand struct {
	X, Y     int // Pixels.
	Relative bool
}

// NewMouseMoveCommand creates a new MouseMoveCommand from a YAML map.
func NewMouseMoveCommand(options map[string]interface{}) (*MouseMoveCommand, error) {
	c := &MouseMoveCommand{}
	foundX := false
	foundY := false

	for key, val := range options {
		switch canonicalize(key) {

		case "x":
			foundX = true
			if !isScalar(val) {
				return nil, wrongOptType("mouse", "scalar", "x", val)
			}
			x, err := strconv.Atoi(val.(string))
			if err != nil {
				return nil, errors.Wrap(err, "mouse: option 'x'")
			}
			c.X = x

		case "y":
			foundY = true
			if !isScalar(val) {
				return nil, wrongOptType("mouse", "scalar", "y", val)
			}
			y, err := strconv.Atoi(val.(string))
			if err != nil {
				return nil, errors.Wrap(err, "mouse: option 'y'")
			}
			c.Y = y

		case "relative":
			if !isScalar(val) {
				return nil, wrongOptType("mouse", "scalar", "relative", val)
			}
			relative, err := strconv.ParseBool(val.(string))
			if err != nil {
				return nil, errors.Wrap(err, "mouse: option 'relative'")
			}
			c.Relative = relative

		default:
			return nil, unrecognizedOpt("mouse", key)
		}
	}

	if !foundX {
		return nil, missingOpt("mouse", "x")
	} else if !foundY {
		return nil, missingOpt("mouse", "y")
	}

	return c, nil

}

// Execute executes the command.
func (c *MouseMoveCommand) Execute() {
	var x, y int

	if c.Relative {
		curX, curY := robotgo.GetMousePos()
		x, y = curX+c.X, curY+c.Y
	} else {
		x, y = c.X, c.Y
	}

	robotgo.MoveMouse(x, y)
	log.Printf("moved mouse to (%v, %v)", x, y)
}

// MouseClickCommand clicks the mouse.
type MouseClickCommand struct {
	Button string // "left", "center", or "right".
}

// NewMouseClickCommand creates a new MouseClickCommand from a YAML map.
func NewMouseClickCommand(options map[string]interface{}) (*MouseClickCommand, error) {
	c := &MouseClickCommand{}
	foundButton := false

	for key, val := range options {
		switch canonicalize(key) {

		case "button":
			foundButton = true
			if !isScalar(val) {
				return nil, wrongOptType("click", "scalar", "button", val)
			}
			button := canonicalize(val.(string))
			if button == "middle" || button == "centre" {
				button = "center"
			}
			if button != "left" && button != "center" && button != "right" {
				return nil, invalidOpt("click", "button", val)
			}
			c.Button = button

		default:
			return nil, unrecognizedOpt("click", key)
		}
	}

	if !foundButton {
		return nil, missingOpt("click", "button")
	}

	return c, nil
}

// Execute executes the command.
func (c *MouseClickCommand) Execute() {
	robotgo.MouseClick(c.Button)
	log.Printf("clicked %v mouse button", c.Button)
}

// Keyboard commands

// KeyPressCommand presses a key.
type KeyPressCommand struct {
	Key  string // TODO: Check robotgo docs for keycode specifics.
	Mods []string
}

// NewKeyPressCommand creates a new KeyPressCommand from a YAML map.
func NewKeyPressCommand(options map[string]interface{}) (*KeyPressCommand, error) {
	c := &KeyPressCommand{}
	foundKey := false

	for key, val := range options {
		switch canonicalize(key) {

		case "key":
			foundKey = true
			if !isScalar(val) {
				return nil, wrongOptType("keypress", "scalar", "key", val)
			}
			c.Key = canonicalize(val.(string))

		case "mods":
			var mods []string
			if !isList(mods) {
				return nil, wrongOptType("keypress", "scalar list", "mods", val)
			}
			for _, mod := range val.([]interface{}) {
				if !isScalar(mod) {
					return nil, wrongListEntryType("keypress", "scalar", "mods", mod)
				}
				mods = append(mods, canonicalize(mod.(string)))
			}
			c.Mods = mods

		default:
			return nil, unrecognizedOpt("keypress", key)
		}
	}

	if !foundKey {
		return nil, missingOpt("keypress", "key")
	}

	return c, nil
}

// Execute executes the command.
func (c *KeyPressCommand) Execute() {
	// Sometimes the key stays pressed for some reason, so we release it manually.
	if len(c.Mods) > 0 {
		robotgo.KeyTap(c.Key, c.Mods)
		// robotgo.KeyToggle(c.Key, "up", c.Mods...) should work but it doesn't.
		args := []string{c.Key, "up"}
		for i := range c.Mods {
			args = append(args, c.Mods[i])
		}
		robotgo.KeyToggle(args...)
		log.Printf("pressed key %v+%v", strings.Join(c.Mods, "+"), c.Key)
	} else {
		robotgo.KeyTap(c.Key)
		robotgo.KeyToggle(c.Key, "up")
		log.Printf("pressed key %v", c.Key)
	}
}

// TypeCommand types some text.
type TypeCommand struct {
	Text string
}

// NewTypeCommand creates a new TypeCommand from a YAML map.
func NewTypeCommand(options map[string]interface{}) (*TypeCommand, error) {
	c := &TypeCommand{}
	foundText := false

	for key, val := range options {
		switch canonicalize(key) {

		case "text":
			foundText = true
			if !isScalar(val) {
				return nil, wrongOptType("type", "scalar", "text", val)
			}
			c.Text = val.(string)

		default:
			return nil, unrecognizedOpt("type", key)
		}
	}

	if !foundText {
		return nil, missingOpt("type", "text")
	}

	return c, nil
}

// Execute executes the command.
func (c *TypeCommand) Execute() {
	robotgo.TypeString(c.Text)
	log.Printf("typed '%s'", strings.Replace(c.Text, "\n", "\\n", -1))
}

// Misc commands

// SleepCommand does nothing for some amount of time.
type SleepCommand struct {
	Seconds int
}

// NewSleepCommand creates a new SleepCommand from a YAML map.
func NewSleepCommand(options map[string]interface{}) (*SleepCommand, error) {
	c := &SleepCommand{}
	foundSeconds := false

	for key, val := range options {
		switch canonicalize(key) {

		case "seconds":
			foundSeconds = true
			if !isScalar(val) {
				return nil, wrongOptType("sleep", "scalar", "seconds", val)
			}
			seconds, err := strconv.Atoi(val.(string))
			if err != nil {
				return nil, errors.Wrap(err, "sleep: option 'seconds'")
			}
			c.Seconds = seconds

		default:
			return nil, unrecognizedOpt("sleep", key)
		}
	}

	if !foundSeconds {
		return nil, missingOpt("sleep", "seconds")
	}

	return c, nil
}

// Execute executes the command.
func (c *SleepCommand) Execute() {
	time.Sleep(time.Second * time.Duration(c.Seconds))
	log.Printf("slept for %v seconds", c.Seconds)
}

// ExecCommand executes a system command.
type ExecCommand struct {
	Program string
	Args    []string
}

// NewExecCommand creates a new ExecCommand from a YAML map.
func NewExecCommand(options map[string]interface{}) (*ExecCommand, error) {
	c := &ExecCommand{}
	foundProgram := false

	for key, val := range options {
		switch canonicalize(key) {

		case "program":
			foundProgram = true
			if !isScalar(val) {
				return nil, wrongOptType("exec", "scalar", "program", val)
			}
			c.Program = strings.TrimSpace(val.(string))

		case "args":
			if !isScalar(val) {
				return nil, wrongOptType("exec", "scalar", "args", val)
			}
			c.Args = strings.Split(val.(string), " ")

		default:
			return nil, unrecognizedOpt("exec", key)
		}
	}

	if !foundProgram {
		return nil, missingOpt("exec", "program")
	}

	return c, nil
}

// Execute executes the command.
func (c *ExecCommand) Execute() {
	cmd := exec.Command(c.Program, c.Args...)
	if err := cmd.Start(); err != nil {
		log.Printf("error executing command %v: %v", cmd, err)
	} else {
		log.Printf("executed command: %v %v", c.Program, strings.Join(c.Args, " "))
	}
}
