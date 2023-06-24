package utils

import (
	"bytes"
	"io"
	"reflect"
)

func StreamToBytes(stream io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func UniqueStringSlice(strings []string) []string {
	// []string to map
	str_map := make(map[string]bool) // bool type is just a dummy
	for _, str := range strings {
		str_map[str] = true
	}
	// map to []string
	var result []string
	for key := range str_map {
		result = append(result, key)
	}
	return result
}

type myType interface {
	~string | ~int
}

// Compare two slice without ordering
func CompareSlicesAsSets[T comparable](x, y []T) bool {
	x_map := make(map[T]bool)
	y_map := make(map[T]bool)

	for _, x_elem := range x {
		x_map[x_elem] = true
	}
	for _, y_elem := range y {
		y_map[y_elem] = true
	}

	return reflect.DeepEqual(x_map, y_map)
}
