package common

func FilterSlice[T any](slice []T, predicate func(t T) bool) []T {
	filteredArr := make([]T, 0, len(slice))
	for _, s := range slice {
		if predicate(s) {
			filteredArr = append(filteredArr, s)
		}
	}

	return filteredArr
}
