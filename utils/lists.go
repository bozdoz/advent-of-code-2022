package utils

// python's all() function, kinda
func All[T any](arr []T, cb func(x T, i int) bool) bool {
	for i, v := range arr {
		if !cb(v, i) {
			return false
		}
	}
	return true
}
