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
