package opml

import (
	"encoding/xml"
	"strings"
)

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    OPMLHead `xml:"head"`
	Body    OPMLBody `xml:"body"`
}

type OPMLHead struct {
	Title string `xml:"title"`
}

type OPMLBody struct {
	Outlines []Outline `xml:"outline"`
}

type Outline struct {
	Text     string    `xml:"text,attr"`
	Title    string    `xml:"title,attr"`
	Type     string    `xml:"type,attr,omitempty"`
	XMLURL   string    `xml:"xmlUrl,attr,omitempty"`
	HTMLURL  string    `xml:"htmlUrl,attr,omitempty"`
	Children []Outline `xml:"outline"`
}

func (o *Outline) Cleanup() {
	o.Text = strings.TrimSpace(o.Text)
	o.Title = strings.TrimSpace(o.Title)
	o.Type = strings.TrimSpace(o.Type)
	o.XMLURL = strings.TrimSpace(o.XMLURL)
	o.HTMLURL = strings.TrimSpace(o.HTMLURL)
	for i := range o.Children {
		o.Children[i].Cleanup()
	}
}
