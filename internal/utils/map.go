package utils

func Map[GIVEN any, RETURN any](collection []GIVEN, fn func(val GIVEN) RETURN) (result []RETURN) {
	result = make([]RETURN, 0, len(collection))
	for _, v := range collection {
		result = append(result, fn(v))
	}
	return result
}
