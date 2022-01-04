package util

import (
	"sort"
)

// GetMapKeyIndicesStringString gets the keys of a string map, sorts them (golang maps don't have fixed iteration order) and inverts index with vals
func GetMapKeyIndicesStringString(input map[string]string) []string {
	keySlice := make([]string, len(input))

	i := 0
	for key, _ := range input {
		keySlice[i] = key
		i++
	}
	sort.Strings(keySlice)

	return keySlice
}
