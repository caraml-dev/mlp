package util

// JoinStringMaps joins multiple string maps into one
// If there are duplicate keys, the last one wins
func JoinStringMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
