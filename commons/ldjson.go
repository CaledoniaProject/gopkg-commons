package commons

import (
	"encoding/json"
	"time"
)

type LDJSON struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	Author           string `json:"author"`
	DateModified     string `json:"dateModified"`
	DatePublished    string `json:"datePublished"`
	Headline         string `json:"headline"`
	Image            string `json:"image"`
	MainEntityOfPage string `json:"mainEntityOfPage"`
	Name             string `json:"name"`
	Publisher        struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		Email   string `json:"email"`
		Logo    struct {
			Context string `json:"@context"`
			Type    string `json:"@type"`
			Name    string `json:"name"`
			URL     string `json:"url"`
		} `json:"logo"`
		Name      string `json:"name"`
		Telephone string `json:"telephone"`
	} `json:"publisher"`
	URL string `json:"url"`

	DatePublished2 NullableTime
}

func NewLDJSON(input []byte) (*LDJSON, error) {
	var (
		ldjson = &LDJSON{}
	)

	if err := json.Unmarshal(input, ldjson); err != nil {
		return nil, err
	}

	if ldjson.DatePublished != "" {
		if tmp, err := time.Parse("2006-01-02", ldjson.DatePublished); err != nil {
			return nil, err
		} else {
			ldjson.DatePublished2 = NullableTime(tmp)
		}
	}

	return ldjson, nil
}
