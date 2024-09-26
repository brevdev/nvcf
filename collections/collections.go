package collections

func Filter[T any](list []T, f func(T) bool) []T {
	result := []T{}
	for _, item := range list {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}

func ListContains[K comparable](list []K, item K) bool {
	return ListHas(list, func(l K) bool { return l == item })
}

func ListHas[K any](list []K, has func(l K) bool) bool {
	k := Find(list, has)
	return k != nil
}

func Find[T any](list []T, f func(T) bool) *T {
	for _, item := range list {
		if f(item) {
			return &item
		}
	}
	return nil
}
