package array

func MapArray[T, V any](arr []T, mapFn func(T) V) []V {
	result := make([]V, len(arr))
	for i, t := range arr {
		result[i] = mapFn(t)
	}
	return result
}

func CountElements[T comparable](arr []T) map[T]int {
	countMap := make(map[T]int)
	for _, element := range arr {
		if _, hasKey := countMap[element]; !hasKey {
			countMap[element] = 0
		}

		countMap[element] = countMap[element] + 1
	}

	return countMap
}
