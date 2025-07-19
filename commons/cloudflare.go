package commons

import (
	"encoding/hex"
)

func DecodeCloudflareEmail(encoded string) (string, error) {
	data, err := hex.DecodeString(encoded)
	if err != nil || len(data) < 1 {
		return "", err
	}

	key := data[0]
	for i := 1; i < len(data); i++ {
		data[i] ^= key
	}

	return string(data[1:]), nil
}
