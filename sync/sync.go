package sync

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/dukex/uhura-api/channels"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/parser"
)

const (
	imageHost            string        = "http://arcane-forest-5063.herokuapp.com"
	episodePubDateFormat string        = "Mon, _2 Jan 2006 15:04:05 -0700"
	weekHours            time.Duration = 24 * 7
)

// Sync channel Class
type Sync struct {
	channelID int64
	model     *models.Channel
	feed      *parser.Channel
}

// NewSync Create new instance of Sync
func NewSync(channelID int64) Sync {
	return Sync{
		channelID: channelID,
	}
}

// Sync syncronization
func (s *Sync) Sync(p gorm.DB) {
	s.defineModel(p)
	s.defineFeed()
	s.update()
	s.cacheImage()
	s.save(p)
	s.episodes(p)
	s.createCategory(p)
}

func (s *Sync) defineModel(p gorm.DB) {
	var model models.Channel
	err := p.Table(models.Channel{}.TableName()).First(&model, s.channelID).Error
	checkError(err)
	s.model = &model
}

func (s *Sync) defineFeed() {
	channelURL, err := url.Parse(s.model.Url)
	checkError(err)

	channelsFeed, _errors := parser.URL(channelURL)
	if len(_errors) > 0 {
		panic(_errors[0])
	}

	if len(channelsFeed) == 0 {
		panic(errors.New("no channels - " + channelURL.String()))
	}

	s.feed = channelsFeed[0]
}

func (s *Sync) update() {
	channels.TranslateFromFeed(s.model, s.feed)
}

func (s *Sync) cacheImage() {
	currentImageURL := s.model.ImageUrl

	if strings.Contains(currentImageURL, imageHost) {
		return
	}

	resp, err := http.Get(imageHost + "/resolve?url=" + currentImageURL)
	if err != nil {
		return
	}

	newImageURL := resp.Request.URL.String()

	if resp.StatusCode == 200 && strings.Contains(newImageURL, imageHost+"/cache") {
		s.model.ImageUrl = newImageURL
	}
}

func (s *Sync) save(p gorm.DB) {
	p.Save(s.model)
	channels.CreateLinks(r(s.feed.Links), s.model.Id, p)
}

func (s *Sync) episodes(p gorm.DB) {
	for _, data := range s.feed.Episodes {
		s.saveEpisode(p, s.buildEpisode(data))
	}
}

func (s *Sync) saveEpisode(p gorm.DB, episode models.Episode) {
	err := p.Table(models.Episode{}.TableName()).
		Where("key = ?", episode.Key).
		Assign(episode).
		FirstOrCreate(&episode).Error

	checkError(err)
}

func (s Sync) buildEpisode(data *parser.Episode) models.Episode {
	description := data.Summary
	if description == "" {
		description = data.Description
	}

	publishedAt, err := data.Feed.ParsedPubDate()
	if err != nil {
		publishedAt, err = s.fixPubDate(data)
		checkError(err)
	}

	audioData := &channels.EpisodeAudioData{
		ContentLength: data.Enclosures[0].Length,
		ContentType:   data.Enclosures[0].Type,
	}

	if audioData.ContentLength == 0 || audioData.ContentType == "" {
		audioData, err = channels.GetEpisodeAudioData(data.Source)
		checkError(err)
	}

	return models.Episode{
		Description:   description,
		Key:           data.GetKey(),
		Uri:           helpers.MakeUri(data.Title),
		Title:         data.Title,
		SourceUrl:     data.Source,
		ChannelId:     s.model.Id,
		PublishedAt:   publishedAt,
		ContentType:   audioData.ContentType,
		ContentLength: audioData.ContentLength,
	}
}

func (s Sync) fixPubDate(e *parser.Episode) (time.Time, error) {
	pubDate := strings.Replace(e.PubDate, "GMT", "-0100", -1)
	pubDate = strings.Replace(pubDate, "PST", "-0800", -1)
	pubDate = strings.Replace(pubDate, "PDT", "-0700", -1)

	return time.Parse(episodePubDateFormat, pubDate)
}

// GetNextRun returns the next run to channel
func (s *Sync) GetNextRun() (time.Time, error) {
	now := time.Now()

	if len(s.feed.Episodes) > 1 {
		last, errLast := s.feed.Episodes[0].Feed.ParsedPubDate()
		if errLast != nil {
			return now, errLast
		}

		penultimate, errPenutimate := s.feed.Episodes[1].Feed.ParsedPubDate()
		if errPenutimate != nil {
			return now, errLast
		}

		// The next run is the duration of last less penultimate episode
		nextRunAt := last.Add(last.Sub(penultimate))

		// If next run date was a old date
		if !nextRunAt.After(now) {
			return now.Add(time.Hour * weekHours), nil
		}

		return nextRunAt, nil
	}

	return now.Add(time.Hour * weekHours), nil
}

func (s *Sync) createCategory(p gorm.DB) {
	for _, cat := range s.feed.Categories {
		var category models.Category

		p.Table(models.Category{}.TableName()).
			Where("name = ?", cat.Name).
			FirstOrCreate(&category)

		p.Table(models.Categoriable{}.TableName()).
			Where("channel_id = ? AND category_id = ?", s.model.Id, category.Id).
			FirstOrCreate(&models.Categoriable{})
	}

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func r(s []string) []string {
	m := map[string]bool{}
	t := []string{}
	for _, v := range s {
		if _, seen := m[v]; !seen {
			t = append(t, v)
			m[v] = true
		}
	}
	return t
}
