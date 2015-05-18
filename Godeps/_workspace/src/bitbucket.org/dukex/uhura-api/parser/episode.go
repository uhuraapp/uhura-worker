package parser

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"

	rss "github.com/jteeuwen/go-pkg-rss"
)

type Episode struct {
	Uid         string           `json:"uid"`
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Subtitle    string           `json:"subtitle"`
	Summary     string           `json:"summary"`
	Description string           `json:"description"`
	Enclosures  []*rss.Enclosure `json:"enclosures"`
	PubDate     string           `json:"pub_date"`
	Duration    string           `json:"duration"`
	Source      string           `json:"source"`

	iTunes
	Feed *rss.Item `json:"-"`
}

func (e *Episode) Build() bool {
	if len(e.Feed.Enclosures) < 1 {
		return false
	}

	e.Uid = e.GetKey()
	e.Title = e.Feed.Title
	e.Enclosures = e.Feed.Enclosures
	e.PubDate = e.Feed.PubDate
	e.Description = e.Feed.Description
	e.Subtitle = e.value(e, "subtitle")
	e.Summary = e.value(e, "summary")
	e.Duration = e.value(e, "duration")
	e.Source = strings.TrimSpace(e.Enclosures[0].Url)

	if e.Feed.Guid != nil {
		e.ID = *e.Feed.Guid
	}

	if e.Source == "" {
		return false
	}

	return true
}

func (e *Episode) GetExtensions(ext string) map[string][]rss.Extension {
	if e.Feed != nil && len(e.Feed.Extensions) > 0 {
		return e.Feed.Extensions[ext]
	}
	return nil
}

func (e *Episode) GetKey() string {
	h := md5.New()
	io.WriteString(h, e.Source)
	return hex.EncodeToString(h.Sum(nil))
}
