package robotscript

import "testing"

var s = `
commands:
  - mouse:
      x: 200
      y: 200
      relative: false
  - click:
      button: left
  - keypress:
      key: enter
      mods:
        - control
        - shift
  - sleep:
      seconds: 1
  - type:
      text: hello, world
  - exec:
      program: stat
      args: /tmp/fakefile
`

func TestNewScript(t *testing.T) {
	script, err := newScript([]byte(s))
	if err != nil {
		t.Errorf("parsing script failed: %v", err)
	}
	if len(script.Commands) != 6 {
		t.Errorf("expected script.Commands == 6, got %v", len(script.Commands))
	}

	if _, err = NewScript("foo"); err == nil {
		t.Error("expected err != nil, got nil")
	}

	if _, err = newScript([]byte("foo")); err == nil {
		t.Error("expected err != nil, got nil")
	}
}
