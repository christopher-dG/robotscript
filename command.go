package robotscript

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kylelemons/go-gypsy/yaml"
)

type Command interface {
	Execute() // Execute executes the command.
	// New_____Command(options yaml.Map) (*_____Command, error)
}

// NewCommand creates a new command from a name and YAML map.
func NewCommand(name string, options yaml.Map) (Command, error) {
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
		return nil, errors.New(fmt.Sprintf("unrecognized command: %v", name))
	}
}

// Mouse commands

// MouseMoveCommand moves the mouse.
type MouseMoveCommand struct {
	X, Y     int // Pixels.
	Relative bool
}

// NewMouseMoveCommand creates a new MouseMoveCommand from a YAML map.
func NewMouseMoveCommand(options yaml.Map) (*MouseMoveCommand, error) {
	c := &MouseMoveCommand{}

	for key, val := range options {
		switch canonicalize(key) {

		case "x":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}
			x, err := strconv.Atoi(scalar.String())
			if err != nil {
				return nil, err
			}
			c.X = x

		case "y":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(scalar.String())
			if err != nil {
				return nil, err
			}
			c.Y = y

		case "relative":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}
			relative, err := strconv.ParseBool(scalar.String())
			if err != nil {
				return nil, err
			}
			c.Relative = relative

		default:
			return nil, errors.New("unrecognized option to mouse command")
		}
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
func NewMouseClickCommand(options yaml.Map) (*MouseClickCommand, error) {
	c := &MouseClickCommand{}

	for key, val := range options {
		switch canonicalize(key) {

		case "button":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}

			button := canonicalize(scalar.String())
			if button == "middle" || button == "centre" {
				button = "center"
			}
			if button != "left" && button != "center" && button != "right" {
				return nil, errors.New("invalid button for mouse click")
			}
			c.Button = button

		default:
			return nil, errors.New("unrecognized option to click command")
		}
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
func NewKeyPressCommand(options yaml.Map) (*KeyPressCommand, error) {
	c := &KeyPressCommand{}

	for key, val := range options {
		switch canonicalize(key) {

		case "key":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}
			c.Key = canonicalize(scalar.String())

		case "mods":
			var mods []string

			modNodes, err := toList(val)
			if err != nil {
				modScalar, err := toScalar(val)
				if err != nil {
					return nil, err
				}
				mods = append(mods, modScalar.String())
			} else {
				for _, node := range modNodes {
					scalar, err := toScalar(node)
					if err != nil {
						return nil, err
					}
					mods = append(mods, canonicalize(scalar.String()))
				}
			}
			c.Mods = mods

		default:
			return nil, errors.New("unrecognized option to keypress command")
		}
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
func NewTypeCommand(options yaml.Map) (*TypeCommand, error) {
	c := &TypeCommand{}

	for key, val := range options {
		switch canonicalize(key) {
		case "text":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}
			c.Text = scalar.String()

		default:
			return nil, errors.New("unrecognized option to type command")
		}
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
func NewSleepCommand(options yaml.Map) (*SleepCommand, error) {
	c := &SleepCommand{}

	for key, val := range options {
		switch canonicalize(key) {

		case "seconds":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}
			seconds, err := strconv.Atoi(scalar.String())
			if err != nil {
				return nil, err
			}
			c.Seconds = seconds

		default:
			return nil, errors.New("unrecognized option to sleep command")
		}
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
func NewExecCommand(options yaml.Map) (*ExecCommand, error) {
	c := &ExecCommand{}

	for key, val := range options {
		switch canonicalize(key) {

		case "program":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}
			c.Program = scalar.String()

		case "args":
			scalar, err := toScalar(val)
			if err != nil {
				return nil, err
			}
			c.Args = strings.Split(scalar.String(), " ")

		default:
			return nil, errors.New("unrecognized option to exec command")
		}
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
