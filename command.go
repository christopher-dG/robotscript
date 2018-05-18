package robotscript

import (
	"log"
	"math"
	"os/exec"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Command interface {
	Execute() // Execute executes the command.
}

// NewCommand creates a new command from an options map.
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
		return nil, errors.Errorf("unrecognized command '%v'", name)
	}
}

// Mouse commands

// MouseMoveCommand moves the mouse.
type MouseMoveCommand struct {
	X, Y     int // Pixels.
	Relative bool
}

// NewMouseMoveCommand creates a new MouseMoveCommand from an options map.
func NewMouseMoveCommand(options map[string]interface{}) (*MouseMoveCommand, error) {
	c := &MouseMoveCommand{}
	if err := checkRequiredOpts(options, "x", "y"); err != nil {
		return nil, errors.Wrap(err, "mouse")
	}
	if err := mapstructure.Decode(options, c); err != nil {
		return nil, errors.Wrap(err, "mouse")
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

	x, y = int(math.Max(float64(x), 0)), int(math.Max(float64(y), 0))

	robotgo.MoveMouse(x, y)
	log.Printf("moved mouse to (%v, %v)", x, y)
}

// MouseClickCommand clicks the mouse.
type MouseClickCommand struct {
	Button string // "left", "center", or "right".
}

// NewMouseClickCommand creates a new MouseClickCommand from an options map.
func NewMouseClickCommand(options map[string]interface{}) (*MouseClickCommand, error) {
	c := &MouseClickCommand{}
	if err := checkRequiredOpts(options, "button"); err != nil {
		return nil, errors.Wrap(err, "click")
	}
	if err := mapstructure.Decode(options, c); err != nil {
		return nil, errors.Wrap(err, "click")
	}
	if c.Button == "middle" || c.Button == "centre" {
		c.Button = "center"
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
	Key  string // https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
	Mods []string
}

// NewKeyPressCommand creates a new KeyPressCommand from an options map.
func NewKeyPressCommand(options map[string]interface{}) (*KeyPressCommand, error) {
	c := &KeyPressCommand{}
	if err := checkRequiredOpts(options, "key"); err != nil {
		return nil, errors.Wrap(err, "keypress")
	}
	if err := mapstructure.Decode(options, c); err != nil {
		return nil, errors.Wrap(err, "keypress")
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

// NewTypeCommand creates a new TypeCommand from an options map.
func NewTypeCommand(options map[string]interface{}) (*TypeCommand, error) {
	c := &TypeCommand{}
	if err := checkRequiredOpts(options, "text"); err != nil {
		return nil, errors.Wrap(err, "type")
	}
	if err := mapstructure.Decode(options, c); err != nil {
		return nil, errors.Wrap(err, "type")
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

// NewSleepCommand creates a new SleepCommand from an options map.
func NewSleepCommand(options map[string]interface{}) (*SleepCommand, error) {
	c := &SleepCommand{}
	if err := checkRequiredOpts(options, "seconds"); err != nil {
		return nil, errors.Wrap(err, "sleep")
	}
	if err := mapstructure.Decode(options, c); err != nil {
		return nil, errors.Wrap(err, "sleep")
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
	Args    string
}

// NewExecCommand creates a new ExecCommand from an options map.
func NewExecCommand(options map[string]interface{}) (*ExecCommand, error) {
	c := &ExecCommand{}
	if err := checkRequiredOpts(options, "program"); err != nil {
		return nil, errors.Wrap(err, "exec")
	}
	if err := mapstructure.Decode(options, c); err != nil {
		return nil, errors.Wrap(err, "exec")
	}
	return c, nil
}

// Execute executes the command.
func (c *ExecCommand) Execute() {
	cmd := exec.Command(c.Program, strings.Split(c.Args, " ")...)
	if err := cmd.Start(); err != nil {
		log.Printf("error executing command %v: %v", cmd, err)
	} else {
		log.Printf("executed command: %v %v", c.Program, c.Args)
	}
}
