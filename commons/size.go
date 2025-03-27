package commons

import "fmt"

func FormatBytes(bytes float64) string {
	var (
		sizes = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
		i     int
	)

	for bytes >= 1024 && i < len(sizes)-1 {
		bytes /= 1024
		i++
	}

	return fmt.Sprintf("%.2f %s", bytes, sizes[i])
}
