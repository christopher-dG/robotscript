package robotscript

import (
	"strings"

	"github.com/pkg/errors"
)

// checkRequiredOpts checks that all required keys are present in a map.
func checkRequiredOpts(m map[string]interface{}, keys ...string) error {
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return errors.Errorf("missing required key '%v'", key)
		}
	}
	return nil
}

// getSingleKey gets a map's first key, as long as it is its only key.
func getSingleKey(m map[string]map[string]interface{}) (string, error) {
	if len(m) == 0 {
		return "", errors.New("map is empty")
	}
	var i int
	var key string
	for k := range m {
		if i > 0 {
			return "", errors.Errorf("map has multiple keys (value = %v)", m)
		}
		i++
		key = k
	}
	return key, nil
}

// canonicalize converts a string to lowercase and removes outer whitespace.
func canonicalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
