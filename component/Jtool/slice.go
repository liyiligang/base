package Jtool

// ReverseSlice 翻转切片
func ReverseSlice[T any](slice []T) []T {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// FilterSlice 过滤切片
func FilterSlice[T any](collection []T, predicate func(item T, index int) bool) []T {
	result := make([]T, 0, len(collection))
	for i, item := range collection {
		if predicate(item, i) {
			result = append(result, item)
		}
	}
	return result
}
