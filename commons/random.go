package commons

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7; rv:119.0) Gecko/20100101 Firefox/119.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edg/119.0.0.0 Safari/537.36",
	}
	charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandomString(length int) string {
	tmp := make([]byte, length)
	for i := range tmp {
		tmp[i] = charset[rand.Intn(len(charset))]
	}

	return string(tmp)
}

func RandomIP() string {
	var (
		ip1 = rand.Intn(200)
		ip2 = rand.Intn(200)
		ip3 = rand.Intn(200)
		ip4 = rand.Intn(200)
	)

	return fmt.Sprintf("%d.%d.%d.%d", ip1, ip2, ip3, ip4)
}

func RandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func RandomStringSlice(input []string) string {
	var (
		rnd   = rand.New(rand.NewSource(time.Now().UnixNano()))
		index = rnd.Intn(len(input))
	)

	return input[index]
}

func RandomStringSliceWithIndex(input []string) (string, int) {
	var (
		rnd   = rand.New(rand.NewSource(time.Now().UnixNano()))
		index = rnd.Intn(len(input))
	)

	return input[index], index
}
