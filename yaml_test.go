package robotscript

import (
	"reflect"
	"testing"

	"github.com/kylelemons/go-gypsy/yaml"
)

func TestIsScalar(t *testing.T) {
	if !isScalar("foo") {
		t.Error("expected 'foo' to be a scalar")
	}
}

func TestIsList(t *testing.T) {
	if !isList([]int{1, 2, 3}) {
		t.Error("expected []int{1, 2, 3} to be a list")
	}
}

func TestIsMap(t *testing.T) {
	if !isMap(make(map[string]int)) {
		t.Error("expected map[string]int to be a map")
	}
}
func TestUnYAML(t *testing.T) {
	file := yaml.Config("a: scalar\nb:\n  - list\n  - value2")
	m := unYAML(file.Root).(map[string]interface{})

	if len(m) != 2 {
		t.Errorf("expected len(m) == 2, got %v", len(m))
	} else if m["a"] != "scalar" {
		t.Errorf("expected m['a'] == 'scalar', got %v", m["a"])
	} else if reflect.DeepEqual(m["b"], []string{"list", "value2"}) {
		t.Errorf("expected m['b'] == []string{'list', 'value2'}, got %v", m["b"])
	}
}

func TestGetSingleKey(t *testing.T) {
	m := make(map[string]interface{})
	m["foo"] = "bar"
	if key, err := getSingleKey(m); err != nil {
		t.Errorf("expected err == nil, got %v", err)
	} else if key != "foo" {
		t.Errorf("expected key == 'foo', got '%v'", key)
	}

	m["bar"] = "foo"
	if key, err := getSingleKey(m); err == nil {
		t.Error("expected err != nil, got nil")
	} else if key != "" {
		t.Errorf("expected key == '', got '%v'", key)
	}
}

func TestCanonicalize(t *testing.T) {
	if canon := canonicalize("foo"); canon != "foo" {
		t.Errorf("expected canonicalize('foo') == 'foo', got '%v'", canon)
	}
	if canon := canonicalize(" foo"); canon != "foo" {
		t.Errorf("expected canonicalize(' foo') == 'foo', got '%v'", canon)
	}
	if canon := canonicalize("foo   "); canon != "foo" {
		t.Errorf("expected canonicalize('foo   ') == 'foo', got '%v'", canon)
	}
	if canon := canonicalize("fOo"); canon != "foo" {
		t.Errorf("expected canonicalize('fOo') == 'foo', got '%v'", canon)
	}
	if canon := canonicalize("  FOo   "); canon != "foo" {
		t.Errorf("expected canonicalize('  FOo   ') == 'foo', got '%v'", canon)
	}

}
