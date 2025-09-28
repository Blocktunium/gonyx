package gormkit

// arrayContains - tells whether arr contains x.
func arrayContains[T comparable](arr *[]T, x T) bool {
	for _, n := range *arr {
		if x == n {
			return true
		}
	}
	return false
}
