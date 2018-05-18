package robotscript

import "testing"

func TestCheckRequiredOpts(t *testing.T) {
	m := map[string]interface{}{"foo": 1, "bar": 2}
	if err := checkRequiredOpts(m, "foo"); err != nil {
		t.Errorf("expected err == nil, got '%v'", err)
	}
	if err := checkRequiredOpts(m, "foo", "bar"); err != nil {
		t.Errorf("expected err == nil, got '%v'", err)
	}
	if err := checkRequiredOpts(m, "foo", "baz"); err == nil {
		t.Error("expected err != nil, got nil")
	}
}

func TestGetSingleKey(t *testing.T) {
	m := map[string]map[string]interface{}{"foo": map[string]interface{}{"bar": 1}}
	if key, err := getSingleKey(m); err != nil {
		t.Errorf("expected err == nil, got %v", err)
	} else if key != "foo" {
		t.Errorf("expected key == 'foo', got '%v'", key)
	}

	m["bar"] = map[string]interface{}{"x": "y"}
	if _, err := getSingleKey(m); err == nil {
		t.Error("expected err != nil, got nil")
	}

	if _, err := getSingleKey(make(map[string]map[string]interface{})); err == nil {
		t.Error("expected err != nil, got nil")
	}
}

func TestCanonicalize(t *testing.T) {
	var canon string
	if canon = canonicalize("foo"); canon != "foo" {
		t.Errorf("expected canonicalize('foo') == 'foo', got '%v'", canon)
	} else if canon := canonicalize(" foo"); canon != "foo" {
		t.Errorf("expected canonicalize(' foo') == 'foo', got '%v'", canon)
	} else if canon := canonicalize("foo   "); canon != "foo" {
		t.Errorf("expected canonicalize('foo   ') == 'foo', got '%v'", canon)
	} else if canon := canonicalize("fOo"); canon != "foo" {
		t.Errorf("expected canonicalize('fOo') == 'foo', got '%v'", canon)
	} else if canon := canonicalize("  FOo   "); canon != "foo" {
		t.Errorf("expected canonicalize('  FOo   ') == 'foo', got '%v'", canon)
	}
}
