package rss

import (
	"encoding/xml"
	"strings"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Content     string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func (i *RSSItem) Cleanup() {
	i.Title = strings.TrimSpace(i.Title)
	i.Link = strings.TrimSpace(i.Link)
	i.Description = strings.TrimSpace(i.Description)
	i.Content = strings.TrimSpace(i.Content)
	i.PubDate = strings.TrimSpace(i.PubDate)
	i.GUID = strings.TrimSpace(i.GUID)
}

type RSSImage struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type RSSChannel struct {
	Title          string               `xml:"title"`
	Link           string               `xml:"link"`
	Description    string               `xml:"description"`
	ManagingEditor string               `xml:"managingEditor"`
	PubDate        commons.NullableTime `xml:"pubDate"`
	LastBuildDate  commons.NullableTime `xml:"lastBuildDate"`
	Image          *RSSImage            `xml:"image"`
	Items          []*RSSItem           `xml:"item"`
}

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}
