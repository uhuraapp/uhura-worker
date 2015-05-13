package parser

import (
	"crypto/md5"
	"encoding/hex"

	"bitbucket.org/dukex/uhura-api/helpers"
	rss "github.com/jteeuwen/go-pkg-rss"
)

type Channel struct {
	Id            string     `json:"id"`
	Title         string     `json:"title"`
	Subtitle      string     `json:"subtitle"`
	Description   string     `json:"description"`
	Summary       string     `json:"summary"`
	Language      string     `json:"language"`
	Copyright     string     `json:"copyright"`
	PubDate       string     `json:"pub_date"`
	Category      string     `json:"category"`
	Author        string     `json:"author"`
	Image         string     `json:"image_url"`
	Links         []string   `json:"links"`
	Episodes      []*Episode `json:"episodes"`
	UhuraId       string     `json:"uhura_id"`
	LastBuildDate string     `json:"last_build_date"`

	requestedURL string
	URL          string

	Feed *rss.Channel `json:"-"`
	iTunes
}

func (c *Channel) HasNewURL() bool {
	return c.NewURL() != ""
}

func (c *Channel) NewURL() string {
	return c.value(c, "new-feed-url")
}

func (c *Channel) Build() {
	c.Title = c.Feed.Title
	c.Description = c.Feed.Description
	c.Copyright = c.Feed.Copyright
	c.Language = c.Feed.Language
	c.PubDate = c.Feed.PubDate
	c.Subtitle = c.value(c, "subtitle")
	c.Summary = c.value(c, "summary")
	c.Author = c.value(c, "author")
	c.Category = c.attr(c, "category", "text")
	c.Image = c.FixImage()
	c.Links = c.GetLinks()
	c.Id = c.GenerateID()
	c.LastBuildDate = c.LastBuildDate

	log.Debug("%s", c.Feed.Links)
	log.Debug("channel build finished: %s", c.Title)
}

func (c *Channel) FixImage() string {
	if image := c.attr(c, "image", "href"); image != "" {
		return image
	} else {
		return c.Feed.Image.Url
	}
}

func (c *Channel) GetExtensions(ext string) map[string][]rss.Extension {
	if c.Feed != nil && len(c.Feed.Extensions) > 0 {
		return c.Feed.Extensions[ext]
	}
	return nil
}

func (c *Channel) GetLinks() []string {
	links := make([]string, 0)
	links = helpers.AppendIfMissing(links, c.requestedURL)
	links = helpers.AppendIfMissing(links, c.URL)

	if len(c.Feed.Links) > 0 {
		links = helpers.AppendIfMissing(links, c.Feed.Links[0].Href)
	}

	return links
}

func (c *Channel) GenerateID() string {
	h := md5.New()
	h.Write([]byte(c.Title))
	return hex.EncodeToString(h.Sum(nil))
}
