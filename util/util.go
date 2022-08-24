package util

// ContainsString checks if a string is contained in an array of strings
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Returns a slice of all keys in a map
func GetMapKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))

	i := 0
	for key := range m {
		keys[i] = key
		i++
	}

	return keys
}
