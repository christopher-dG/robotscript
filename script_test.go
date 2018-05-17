package robotscript

import (
	"testing"
)

func TestNewScript(t *testing.T) {
	script, err := NewScript("test.yml")
	if err != nil {
		t.Errorf("parsing script failed: %v", err)
	}
	if len(script.Commands) != 6 {
		t.Errorf("expected script.Commands == 6, got %v", len(script.Commands))
	}
}
