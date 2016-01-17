package parser

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"bitbucket.org/dukex/uhura-api/helpers"
	rss "github.com/jteeuwen/go-pkg-rss"
)

// Channel parsed
type Channel struct {
	ID            string      `json:"id"`
	Title         string      `json:"title"`
	Subtitle      string      `json:"subtitle"`
	Description   string      `json:"description"`
	Summary       string      `json:"summary"`
	Language      string      `json:"language"`
	Copyright     string      `json:"copyright"`
	PubDate       string      `json:"pub_date"`
	Categories    []*Category `json:"categories"`
	Author        string      `json:"author"`
	Image         string      `json:"image_url"`
	Links         []string    `json:"links"`
	Episodes      []*Episode  `json:"episodes"`
	UhuraID       string      `json:"uhura_id"`
	LastBuildDate string      `json:"last_build_date"`

	requestedURL string
	URL          string

	Feed *rss.Channel `json:"-"`
	iTunes
}

// HasNewURL check if channel has a new URL
func (c *Channel) HasNewURL() bool {
	return c.NewURL() != ""
}

// NewURL get feed new URL
func (c *Channel) NewURL() string {
	return c.value(c, "new-feed-url")
}

// Build channel from parsed rss
func (c *Channel) Build() {
	c.Title = c.Feed.Title
	c.Description = c.Feed.Description
	c.Copyright = c.Feed.Copyright
	c.Language = NormalizeLanguage(c.Feed.Language)
	c.PubDate = c.Feed.PubDate
	c.Subtitle = c.value(c, "subtitle")
	c.Summary = c.value(c, "summary")
	c.Author = c.value(c, "author")
	c.Image = c.FixImage()
	c.Links = c.GetLinks()
	c.ID = c.GenerateID()
	c.LastBuildDate = c.Feed.LastBuildDate

	categories := c.attrs(c, "category", "text")
	for _, category := range categories {
		c.Categories = append(c.Categories, &Category{
			Name: category,
		})
	}

	log.Debug("%s", c.Feed.Links)
	log.Debug("channel build finished: %s", c.Title)
}

var Languages = map[string]string{
	"br": "portuguese",
	"da": "danish",
	"de": "german",
	"en": "english",
	"es": "spanish",
	"fi": "finnish",
	"fr": "french",
	"hu": "hungarian",
	"it": "italian",
	"jp": "japan",
	"nl": "dutch",
	"no": "norwegian",
	"pt": "portuguese",
	"ro": "romanian",
	"ru": "russian",
	"sv": "swedish",
	"tr": "turkish",
}

func NormalizeLanguage(language string) string {
	// Normalize word: remove spaces and downcase
	language = strings.ToLower(strings.TrimSpace(language))

	// remove region from language, so 'en-gb' -> en
	_language := strings.Split(language, "-")
	language = _language[0]

	if language == "" {
		// Default is english
		return "english"
	}

	if newLanguage, ok := Languages[language]; ok {
		return newLanguage
	}

	return language
}

// FixImage get correct image
func (c *Channel) FixImage() string {
	if image := c.attr(c, "image", "href"); image != "" {
		return image
	}

	return c.Feed.Image.Url
}

// GetExtensions from rss
func (c *Channel) GetExtensions(ext string) map[string][]rss.Extension {
	if c.Feed != nil && len(c.Feed.Extensions) > 0 {
		return c.Feed.Extensions[ext]
	}
	return nil
}

// GetLinks from urls requested and feed data
func (c *Channel) GetLinks() []string {
	links := make([]string, 0)
	links = helpers.AppendIfMissing(links, c.requestedURL)
	links = helpers.AppendIfMissing(links, c.URL)

	if len(c.Feed.Links) > 0 {
		links = helpers.AppendIfMissing(links, c.Feed.Links[0].Href)
	}

	return links
}

// GenerateID md5 from title
func (c *Channel) GenerateID() string {
	h := md5.New()
	h.Write([]byte(c.Title))
	return hex.EncodeToString(h.Sum(nil))
}
