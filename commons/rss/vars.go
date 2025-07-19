package rss

import "strings"

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func (i *RSSItem) Cleanup() {
	i.Title = strings.TrimSpace(i.Title)
	i.Link = strings.TrimSpace(i.Link)
	i.Description = strings.TrimSpace(i.Description)
	i.PubDate = strings.TrimSpace(i.PubDate)
	i.GUID = strings.TrimSpace(i.GUID)
}

type RSSChannel struct {
	Items []*RSSItem `xml:"item"`
}

type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}
