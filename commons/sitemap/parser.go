package sitemap

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"

	"github.com/CaledoniaProject/gopkg-commons/commons"
	"github.com/pkg/errors"
)

type SitemapParser struct {
	OnURLFound func(url *SitemapURL)
}

func (p *SitemapParser) LoadURL(sitemapURL string) error {
	if _, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         sitemapURL,
		ReadBody:    true,
		MaxBodyRead: 10 * 1024 * 1024, // 10MB
	}); err != nil {
		return errors.Wrapf(err, "load sitemap")
	} else {
		return p.parseXML(bytes.NewReader(body))
	}
}

func (p *SitemapParser) LoadRobots(robotsURL string) error {
	_, body, err := commons.HttpRequest(&commons.RequestOptions{
		URL:         robotsURL,
		ReadBody:    true,
		MaxBodyRead: 1024 * 1024, // 1MB
	})
	if err != nil {
		return errors.Wrapf(err, "load robots")
	}

	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(strings.ToLower(line), "sitemap:") {
			continue
		}

		sitemapURL := strings.TrimSpace(line[len("sitemap:"):])
		if sitemapURL != "" {
			if err := p.LoadURL(sitemapURL); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *SitemapParser) parseXML(r io.Reader) error {
	decoder := xml.NewDecoder(r)
	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch se := tok.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "url":
				url := &SitemapURL{}
				if err := decoder.DecodeElement(&url, &se); err == nil && p.OnURLFound != nil {
					p.OnURLFound(url)
				}
			case "sitemap":
				var loc SitemapLoc
				if err := decoder.DecodeElement(&loc, &se); err == nil {
					_ = p.LoadURL(loc.Loc)
				}
			}
		}
	}
}
