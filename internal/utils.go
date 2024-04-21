package internal

// contains checks if a string is contained in a slice of strings.
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// Convert a slice to a map for quicker lookups
func sliceToMap(slice []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range slice {
		result[item] = struct{}{} // Use an empty struct to minimize memory usage
	}
	return result
}
