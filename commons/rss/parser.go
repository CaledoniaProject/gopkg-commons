package rss

import (
	"encoding/xml"

	"github.com/CaledoniaProject/gopkg-commons/commons"
	"github.com/pkg/errors"
)

type RSSParser struct {
	OnURLFound func(*RSSItem)
}

func (r *RSSParser) Load(feedURL string) error {
	if _, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         feedURL,
		Timeout:     120,
		ReadBody:    true,
		MaxBodyRead: 10 * 1024 * 1024, // 1MB
	}); err != nil {
		return errors.Wrapf(err, "load rss feed: %s", feedURL)
	} else {
		return r.parseRSS(body)
	}
}

func (r *RSSParser) parseRSS(data []byte) error {
	var (
		feed RSSFeed
	)

	if err := xml.Unmarshal(data, &feed); err != nil {
		return errors.Wrapf(err, "parse rss xml")
	}

	for _, item := range feed.Channel.Items {
		item.Image = feed.Channel.Image
		item.Cleanup()

		if r.OnURLFound != nil {
			r.OnURLFound(item)
		}
	}

	return nil
}
