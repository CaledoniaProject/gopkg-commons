package commons

import "fmt"

func FormatBytes(bytes float64) string {
	var (
		sizes = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
		i     int
		neg   = bytes < 0
	)

	if neg {
		bytes = -bytes
	}

	for bytes >= 1024 && i < len(sizes)-1 {
		bytes /= 1024
		i++
	}

	sign := ""
	if neg {
		sign = "-"
	}

	return fmt.Sprintf("%s%.2f %s", sign, bytes, sizes[i])
}
