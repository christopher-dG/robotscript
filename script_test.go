package robotscript

import (
	"testing"
)

func TestNewScript(t *testing.T) {
	script, err := NewScript("test.yml")
	if err != nil {
		t.Errorf("Parsing script failed: %v", err)
	}
	if len(script.Commands) != 6 {
		t.Errorf("Expected script.Commands == 5, got %d", len(script.Commands))
	}
}
