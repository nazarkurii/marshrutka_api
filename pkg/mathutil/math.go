package mathutil

func IfLessThan(value, replacement int) int {
	if value > 0 {
		return value
	}
	return replacement
}
