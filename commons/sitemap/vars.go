package sitemap

import (
	"encoding/xml"
	"time"

	"github.com/CaledoniaProject/gopkg-commons/commons"
)

const (
	urlSetStart = `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1">`
	urlSetEnd   = `</urlset>`
	indexStart  = `<?xml version="1.0" encoding="UTF-8"?><sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	indexEnd    = `</sitemapindex>`
)

type ChangeFreq string

const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

type SitemapImage struct {
	Loc string `xml:"image:loc"`
}

type SitemapURL struct {
	XMLName    xml.Name             `xml:"url"`
	Loc        string               `xml:"loc"`
	LastMod    commons.NullableTime `xml:"lastmod"`
	ChangeFreq ChangeFreq           `xml:"changefreq"`
	Priority   float64              `xml:"priority"`
	Images     []*SitemapImage      `xml:"image:image"`
}

type SitemapLoc struct {
	XMLName xml.Name  `xml:"sitemap"`
	Loc     string    `xml:"loc"`
	LastMod time.Time `xml:"lastmod"`
}
