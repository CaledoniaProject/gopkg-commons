package opml

import (
	"encoding/xml"
	"io"
)

func ParseOPMLFromReader(r io.Reader) (*OPML, error) {
	var (
		opml = &OPML{}
	)

	if err := xml.NewDecoder(r).Decode(opml); err != nil {
		return nil, err
	}

	for i := range opml.Body.Outlines {
		opml.Body.Outlines[i].Cleanup()
	}

	return opml, nil
}
