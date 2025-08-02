package atomrss

import (
	"encoding/xml"

	"github.com/CaledoniaProject/gopkg-commons/commons"
	"github.com/pkg/errors"
)

type AtomParser struct {
	OnURLFound func(*AtomEntry)
}

func (p *AtomParser) Load(feedURL string) error {
	if _, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         feedURL,
		Timeout:     120,
		ReadBody:    true,
		MaxBodyRead: 10 * 1024 * 1024, // 10MB
	}); err != nil {
		return errors.Wrapf(err, "load atom feed: %s", feedURL)
	} else {
		return p.parseAtom(body)
	}
}

func (p *AtomParser) parseAtom(data []byte) error {
	var feed AtomFeed

	if err := xml.Unmarshal(data, &feed); err != nil {
		return errors.Wrapf(err, "parse atom xml")
	}

	for _, entry := range feed.Entries {
		entry.Cleanup()
		if p.OnURLFound != nil {
			p.OnURLFound(&entry)
		}
	}

	return nil
}
