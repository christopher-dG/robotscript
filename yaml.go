package robotscript

import (
	"fmt"
	"strings"

	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/pkg/errors"
)

// toScalar converts a Node to a Scalar.
func toScalar(node yaml.Node) (yaml.Scalar, error) {
	scalar, ok := node.(yaml.Scalar)
	if !ok {
		return "", errors.Errorf("%v is not a scalar", node)
	}
	return scalar, nil
}

// toList converts a Node to a List.
func toList(node yaml.Node) (yaml.List, error) {
	list, ok := node.(yaml.List)
	if !ok {
		return nil, errors.Errorf("%v is not a list", node)
	}
	return list, nil
}

// toMap converts a Node to a Map.
func toMap(node yaml.Node) (yaml.Map, error) {
	m, ok := node.(yaml.Map)
	if !ok {
		return nil, errors.Errorf("%v is not a map", node)
	}
	return m, nil
}

// Determine if l is a scalar. Note that this is not for yaml types.
func isScalar(s interface{}) bool {
	return fmt.Sprintf("%T", s) == "string"
}

// Determine if l is a list. Note that this is not for yaml types.
func isList(l interface{}) bool {
	return strings.HasPrefix(fmt.Sprintf("%T", l), "[]")
}

// Determine if m is a map. Note that this is not for yaml types.
func isMap(m interface{}) bool {
	return strings.HasPrefix(fmt.Sprintf("%T", m), "map[")
}

// unYAML converts the values in a yaml node to their literal values.
func unYAML(node yaml.Node) interface{} {
	if scalar, err := toScalar(node); err == nil {
		return scalar.String()
	} else if list, err := toList(node); err == nil {
		var literalList []interface{}
		for _, listNode := range list {
			literalList = append(literalList, unYAML(listNode))
		}
		return literalList
	} else if m, err := toMap(node); err == nil {
		literalMap := make(map[string]interface{})
		for key, val := range m {
			literalMap[key] = unYAML(val)
		}
		return literalMap
	} else {
		return nil // This should never happen.
	}
}

// getSingleKey gets a map's only key.
// If the map has more than one key, an error is returned.
func getSingleKey(m map[string]interface{}) (string, error) {
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
