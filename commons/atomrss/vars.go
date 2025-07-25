package atomrss

import (
	"encoding/xml"
	"strings"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

type AtomEntry struct {
	Title     string               `xml:"title"`
	Link      AtomLink             `xml:"link"`
	Summary   string               `xml:"summary"`
	Content   string               `xml:"content"`
	Published commons.NullableTime `xml:"published"`
	Updated   commons.NullableTime `xml:"updated"`
	ID        string               `xml:"id"`
}

func (e *AtomEntry) Cleanup() {
	e.Title = strings.TrimSpace(e.Title)
	e.Summary = strings.TrimSpace(e.Summary)
	e.Content = strings.TrimSpace(e.Content)
	e.ID = strings.TrimSpace(e.ID)
	e.Link.Href = strings.TrimSpace(e.Link.Href)
}

type AtomLink struct {
	Href string `xml:"href,attr"`
}

type AtomFeed struct {
	XMLName xml.Name             `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string               `xml:"title"`
	Link    AtomLink             `xml:"link"`
	Updated commons.NullableTime `xml:"updated"`
	Entries []AtomEntry          `xml:"entry"`
}
