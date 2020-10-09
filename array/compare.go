package array

// StringElementsMatch compares two arrays of strings irrespective of order.
func StringElementsMatch(one, two []string) bool {
	if len(one) != len(two) {
		return false
	}
	diff := make(map[string]bool)
	for _, dim := range one {
		diff[dim] = true
	}
	for _, dim := range two {
		if !diff[dim] {
			return false
		}
	}
	return true
}
