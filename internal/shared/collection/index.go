package collection

func IndexBy[T any, K comparable](items []T, key func(T) K) map[K]T {
	indexed := make(map[K]T, len(items))
	for _, item := range items {
		indexed[key(item)] = item
	}

	return indexed
}
