package parser

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"regexp"
	"strings"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
)

const (
	episodePubDateFormat                   string = "Mon, _2 Jan 2006 15:04:05 -0700"
	episodePubDateFormatWithoutMiliseconds string = "Mon, _2 Jan 2006 15:04 -0700"
	episodePubDateFormatRFC822Extendend    string = "_2 Mon 2006 15:04:05 -0700"
)

var (
	dateWithoutMiliseconds = regexp.MustCompile(`^\w{3}.{13,14}\d{2}:\d{2}\s`)
	dateRFC822Extedend     = regexp.MustCompile(`^\d{2}.\w{3}.\d{4}.\d{2}:\d{2}:\d{2}.-\d{4}`)
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
	PublishedAt time.Time

	iTunes
	Feed *rss.Item `json:"-"`
}

func (e *Episode) Build() bool {
	if len(e.Feed.Enclosures) < 1 {
		// log.Debug("No enclosures")
		return false
	}

	e.Uid = e.GetKey()
	e.Title = e.Feed.Title
	e.Enclosures = e.Feed.Enclosures
	e.PubDate = e.Feed.PubDate
	e.Subtitle = e.value(e, "subtitle")
	e.Duration = e.value(e, "duration")
	e.Source = strings.TrimSpace(e.Enclosures[0].Url)

	description := e.value(e, "summary")
	if description == "" {
		description = e.Feed.Description
	}
	e.Description = description

	var publishedAt time.Time
	publishedAt, err := e.Feed.ParsedPubDate()
	if err != nil {
		publishedAt, _ = e.FixPubDate()
	}
	e.PublishedAt = publishedAt

	if e.Feed.Guid != nil {
		e.ID = *e.Feed.Guid
	}

	if e.Source == "" {
		// log.Debug("No source")
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

func (e *Episode) FixPubDate() (time.Time, error) {
	pubDate := strings.Replace(e.PubDate, "-5GMT", "-0500", -1)
	pubDate = strings.Replace(e.PubDate, "GMT", "-0100", -1)
	pubDate = strings.Replace(pubDate, "PST", "-0800", -1)
	pubDate = strings.Replace(pubDate, "PDT", "-0700", -1)
	pubDate = strings.Replace(pubDate, "EDT", "-0400", -1)

	if dateWithoutMiliseconds.MatchString(pubDate) {
		return time.Parse(episodePubDateFormatWithoutMiliseconds, pubDate)
	}

	if dateRFC822Extedend.MatchString(pubDate) {
		return time.Parse(episodePubDateFormatRFC822Extendend, pubDate)
	}

	return time.Parse(episodePubDateFormat, pubDate)
}
