package commons

import (
	"net/url"
)

func MustParseURL(input string) *url.URL {
	urlObj, err := url.Parse(input)
	if err != nil {
		panic(err)
	}

	return urlObj
}
