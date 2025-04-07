package commons

func IntSliceAppendDedup(target []int, value int) []int {
	for _, v := range target {
		if v == value {
			return target
		}
	}

	return append(target, value)
}
