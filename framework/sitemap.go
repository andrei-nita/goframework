package framework

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"path"
)

// ChangeFreq specifies the change frequency of a sitemap entry. It is just a string.
type ChangeFreq string

// Feel free to use these constants for ChangeFreq (or you can just supply
// a string directly).
const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

// URL entry in sitemap or sitemap index. LastMod is a pointer
// to time.Time because omitempty does not work otherwise. Loc is the
// only mandatory item. ChangeFreq and Priority must be left empty when
// using with a sitemap index.
type URL struct {
	Loc        string     `xml:"loc"`
	LastMod    string     `xml:"lastmod,omitempty"`
	ChangeFreq ChangeFreq `xml:"changefreq,omitempty"`
	Priority   float32    `xml:"priority,omitempty"`
}

// Sitemap represents a complete sitemap which can be marshaled to XML.
// New instances must be created with New() in order to set the xmlns
// attribute correctly. Minify can be set to make the output less human-
// readable.
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`

	URLs []*URL `xml:"url"`

	Minify bool `xml:"-"`
}

// New returns a new Sitemap.
func New() *Sitemap {
	return &Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]*URL, 0),
	}
}

// Add adds a URL to a Sitemap.
func (s *Sitemap) Add(u *URL) {
	s.URLs = append(s.URLs, u)
}

// WriteToFile writes XML encoded sitemap to given io.Writer.
// Implements io.WriterTo.
func (s *Sitemap) WriteToFile(pathFile string) error {
	smf, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	defer smf.Close()

	_, err = smf.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	en := xml.NewEncoder(smf)
	if !s.Minify {
		en.Indent("", "  ")
	}
	err = en.Encode(s)
	smf.Write([]byte{'\n'})
	smf.Sync()
	return err
}

// ReadFrom reads and parses an XML encoded sitemap from io.Reader.
// Implements io.ReaderFrom.
func (s *Sitemap) ReadFrom(r io.Reader) (n int64, err error) {
	de := xml.NewDecoder(r)
	err = de.Decode(s)
	return de.InputOffset(), err
}

var _ io.ReaderFrom = (*Sitemap)(nil)

func BuildSitemap(urls ...*URL) {
	sm := New()

	for _, url := range urls {
		sm.Add(url)
	}

	err := sm.WriteToFile("sitemap.xml")
	if err != nil {
		log.Println(err)
	}
}

func CreateURL(loc string, date string, changeFreq ChangeFreq, priority float32) *URL {
	return &URL{
		Loc:        loc,
		LastMod:    date,
		ChangeFreq: changeFreq,
		Priority:   priority,
	}
}

func Domain(pathDomain string) string {
	return path.Join(Server.Domain, pathDomain)
}
