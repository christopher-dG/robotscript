package robotscript

import (
	"fmt"
	"testing"
)

func TestWrongOptType(t *testing.T) {
	err := wrongOptType("a", "b", "c", "d")
	expected := "a: expected b value to 'c' option (value = d)"
	if fmt.Sprintf("%v", err) != expected {
		t.Error("expected err == expected")
	}
}

func TestWrongListEntryType(t *testing.T) {
	err := wrongListEntryType("a", "b", "c", "d")
	expected := "a: expected b entries in 'c' option (value = d)"
	if fmt.Sprintf("%v", err) != expected {
		t.Error("expected err == expected")
	}
}

func TestUnrecognizedOpt(t *testing.T) {
	err := unrecognizedOpt("a", "b")
	expected := "a: unrecognized option 'b'"
	if fmt.Sprintf("%v", err) != expected {
		t.Error("expected err == expected")
	}
}

func TestMissingOpt(t *testing.T) {
	err := missingOpt("a", "b")
	expected := "a: missing required option 'b'"
	if fmt.Sprintf("%v", err) != expected {
		t.Error("expected err == expected")
	}
}

func TestInvalidOpt(t *testing.T) {
	err := invalidOpt("a", "b", "c")
	expected := "a: invalid value for option 'b' (value = c)"
	if fmt.Sprintf("%v", err) != expected {
		t.Error("expected err == expected")
	}
}

func TestUnrecognizedCmd(t *testing.T) {
	err := unrecognizedCmd("a")
	expected := "unrecognized command 'a'"
	if fmt.Sprintf("%v", err) != expected {
		t.Error("expected err == expected")
	}
}
