package robotscript

import (
	"errors"
	"fmt"

	"github.com/kylelemons/go-gypsy/yaml"
)

// toList converts a Node to a Map.
func toMap(node yaml.Node) (yaml.Map, error) {
	m, ok := node.(yaml.Map)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%v is not a map", node))
	}
	return m, nil
}

// toList converts a Node to a List.
func toList(node yaml.Node) (yaml.List, error) {
	list, ok := node.(yaml.List)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%v is not a list", node))
	}
	return list, nil
}

// toScalar converts a Node to a Scalar.
func toScalar(node yaml.Node) (yaml.Scalar, error) {
	scalar, ok := node.(yaml.Scalar)
	if !ok {
		return "", errors.New(fmt.Sprintf("%v is not a scalar", node))
	}
	return scalar, nil
}

// getSingleKey gets a map's only key.
// If the map has more than one key, an error is returned.
func getSingleKey(m yaml.Map) (string, error) {
	var i int
	var key string
	for k := range m {
		if i > 0 {
			return "", errors.New("map has multiple keys")
		}
		i++
		key = k
	}
	return key, nil
}
