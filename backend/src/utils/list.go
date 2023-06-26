package utils

func Intersects[T string | int](listA []T, listB []T) bool {
	for _, elemA := range listA {
		for _, elemB := range listB {
			if elemA == elemB {
				return true
			}
		}
	}

	return false
}
