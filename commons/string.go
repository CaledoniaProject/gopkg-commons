package commons

import "strings"

func UpperCaseFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func StringSliceContainsNoCase(haystack []string, needle string) bool {
	for _, row := range haystack {
		if strings.EqualFold(row, needle) {
			return true
		}
	}

	return false
}

func StringSliceContains(haystack []string, needle string) bool {
	for _, row := range haystack {
		if row == needle {
			return true
		}
	}

	return false
}

func TrimSplit(input, sep string) []string {
	var (
		result = []string{}
	)

	for _, tmp := range strings.Split(input, sep) {
		tmp = strings.TrimSpace(tmp)
		if tmp != "" {
			result = append(result, tmp)
		}
	}

	return result
}
