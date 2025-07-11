package sitemap

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type SitemapGenerator struct {
	// configurations
	OutputDirectory string // Where to hold sitemaps
	LinksPerFile    int    // maximum items per file, 5000 for google

	// internal variables
	files         []string
	writer        *gzip.Writer
	numberOfLinks int
}

func (s *SitemapGenerator) Add(sitemapURL *SitemapURL) error {
	// create writer
	if s.writer == nil {
		filename := fmt.Sprintf("file-%d.gz", len(s.files)+1)

		// create directory
		if err := os.MkdirAll(s.OutputDirectory, 0755); err != nil {
			return err
		} else if filp, err := os.OpenFile(s.OutputDirectory+"/"+filename, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return err
		} else {
			s.writer = gzip.NewWriter(filp)
			s.files = append(s.files, filename)
		}

		if _, err := s.writer.Write([]byte(urlSetStart)); err != nil {
			return err
		}
	}

	// write marshalled xml
	if xmlData, err := xml.Marshal(sitemapURL); err != nil {
		return err
	} else if _, err := s.writer.Write(xmlData); err != nil {
		return err
	} else {
		s.numberOfLinks++

		// move to next file
		if s.numberOfLinks == s.LinksPerFile {
			s.Close()
		}
	}

	return nil
}

func (s *SitemapGenerator) Close() error {
	// no active file
	if s.writer == nil {
		return nil
	}

	// write closing header and close
	if _, err := s.writer.Write([]byte(urlSetEnd)); err != nil {
		return err
	} else if err := s.writer.Close(); err != nil {
		return err
	}

	// reset staths
	s.numberOfLinks = 0
	s.writer = nil
	return nil
}

func (s *SitemapGenerator) WriteIndex(baseURL string) (string, error) {
	var (
		indexFile = s.OutputDirectory + "/sitemap-index.xml"
	)
	if err := s.Close(); err != nil {
		return "", err
	}

	filp, err := os.OpenFile(indexFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	} else if _, err := filp.Write([]byte(indexStart)); err != nil {
		return "", err
	}

	for _, file := range s.files {
		if data, err := xml.Marshal(&SitemapLoc{
			Loc:     baseURL + "/" + file,
			LastMod: time.Now(),
		}); err != nil {
			return "", err
		} else if _, err := filp.Write(data); err != nil {
			return "", err
		}
	}

	if _, err := filp.Write([]byte(indexEnd)); err != nil {
		return "", err
	} else if err := filp.Close(); err != nil {
		return "", err
	}

	return indexFile, nil
}
