package commons

import (
	"net/url"

	"github.com/pkg/errors"
)

func MustParseURL(input string) *url.URL {
	urlObj, err := url.Parse(input)
	if err != nil {
		panic(err)
	}

	return urlObj
}

func ResolveURL(currentURL string, href string) (string, error) {
	current, err := url.Parse(currentURL)
	if err != nil {
		return "", errors.Wrapf(err, "parse currentURL: %s", currentURL)
	}

	ref, err := url.Parse(href)
	if err != nil {
		return "", errors.Wrapf(err, "parse href: %s", href)
	}

	finalURL := current.ResolveReference(ref)
	return finalURL.String(), nil
}
