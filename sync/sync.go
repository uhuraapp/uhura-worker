package sync

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/uhuraapp/uhura-api/channels"
	"github.com/uhuraapp/uhura-api/helpers"
	"github.com/uhuraapp/uhura-api/models"
	"github.com/uhuraapp/uhura-api/parser"
)

const (
	episodePubDateFormat                   string        = "Mon, _2 Jan 2006 15:04:05 -0700"
	episodePubDateFormatWithoutMiliseconds string        = "Mon, _2 Jan 2006 15:04 -0700"
	episodePubDateFormatRFC822Extendend    string        = "_2 Mon 2006 15:04:05 -0700"
	weekHours                              time.Duration = 24 * 7
)

var (
	imageHosts             []string = []string{"https://images1.uhura.io", "https://images2.uhura.io"}
	imageHostRegEx                  = regexp.MustCompile(`.+images\d.uhura.io`)
	dateWithoutMiliseconds          = regexp.MustCompile(`^\w{3}.{13,14}\d{2}:\d{2}\s`)
	dateRFC822Extedend              = regexp.MustCompile(`^\d{2}.\w{3}.\d{4}.\d{2}:\d{2}:\d{2}.-\d{4}`)
)

func Sync(channelID int64, p *gorm.DB) (model models.Channel, feed parser.Channel) {
	model = GetModel(channelID, p)

	if model.Enabled {
		_feed := getFeed(model)
		if _feed == nil {
			model.Enabled = false
		} else {
			feed = *_feed
			model = translate(model, feed)
			model = cacheImage(model)

			createLinks(feed, model, p)
			episodes := getEpisodes(feed, model)
			log.Println(episodes)

			for _, episode := range episodes {
				saveEpisode(episode, p)
			}

			createCategory(feed, model, p)
		}
		p.Save(model)
	}

	return model, feed
}

func GetModel(id int64, p *gorm.DB) (model models.Channel) {
	err := p.Table(models.Channel{}.TableName()).First(&model, id).Error
	checkError(err)

	return model
}

func getFeed(model models.Channel) *parser.Channel {
	channelURL, err := url.Parse(model.Url)
	checkError(err)

	channelsFeed, parserError := parser.URL(channelURL)
	if parserError != nil {
		if parserError.Error() == "Status is 404 Not Found" {
			return nil
		} else {
			panic(parserError.Error())
		}
	}

	if channelsFeed == nil {
		panic(errors.New("Not found channel in " + channelURL.String()))
	}

	return channelsFeed
}

func translate(model models.Channel, feed parser.Channel) models.Channel {
	return channels.TranslateFromFeed(model, &feed)
}

func cacheImage(model models.Channel) models.Channel {
	currentImageURL := model.ImageUrl

	if imageHostRegEx.MatchString(currentImageURL) {
		return model
	}

	imageHost := random(imageHosts)
	resp, err := http.Get(imageHost + "/resolve?url=" + currentImageURL)
	if err != nil {
		return model
	}

	newImageURL := resp.Request.URL.String()

	if resp.StatusCode == 200 && strings.Contains(newImageURL, imageHost+"/cache") {
		model.ImageUrl = newImageURL
	}

	return model
}

func createLinks(feed parser.Channel, model models.Channel, p *gorm.DB) {
	channels.CreateLinks(unique(feed.Links), model.Id, p)
}

func getEpisodes(feed parser.Channel, model models.Channel) []models.Episode {
	episodes := make([]models.Episode, 0)

	for _, data := range feed.Episodes {
		episodes = append(episodes, buildEpisode(data, model))
	}
	return episodes
}

func buildEpisode(data *parser.Episode, model models.Channel) models.Episode {
	description := data.Summary
	if description == "" {
		description = data.Description
	}

	var err error

	// hack to feed without date
	var publishedAt time.Time
	if data.Feed.PubDate != "" {
		publishedAt, err = data.Feed.ParsedPubDate()
		if err != nil {
			publishedAt, err = fixPubDate(data)
			checkError(err)
		}
	}

	audioData := &channels.EpisodeAudioData{
		ContentLength: data.Enclosures[0].Length,
		ContentType:   data.Enclosures[0].Type,
	}

	if audioData.ContentLength == 0 || audioData.ContentType == "" {
		// audioData, err = channels.GetEpisodeAudioData(data.Source)
		// return models.Episode{}, err
	}

	return models.Episode{
		Description:   description,
		Key:           data.GetKey(),
		Uri:           helpers.MakeUri(data.Title),
		Title:         data.Title,
		SourceUrl:     data.Source,
		ChannelId:     model.Id,
		PublishedAt:   publishedAt,
		ContentType:   audioData.ContentType,
		ContentLength: audioData.ContentLength,
	}
}

func fixPubDate(e *parser.Episode) (time.Time, error) {
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

func saveEpisode(episode models.Episode, p *gorm.DB) models.Episode {
	err := p.Table(models.Episode{}.TableName()).
		Where("source_url = ?", episode.SourceUrl).First(&models.Episode{}).Error

	if err == gorm.ErrRecordNotFound {
		err = p.Table(models.Episode{}.TableName()).
			Where("key = ?", episode.Key).
			Assign(episode).
			FirstOrCreate(&episode).Error

		checkError(err)
	}

	return episode
}

// GetNextRun returns the next run to channel
func GetNextRun(feed parser.Channel) (time.Time, error) {
	now := time.Now()

	if len(feed.Episodes) > 1 {
		last, errLast := fixPubDate(feed.Episodes[0])
		if errLast != nil {
			return now, errLast
		}

		penultimate, errPenutimate := fixPubDate(feed.Episodes[1])
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

func createCategory(feed parser.Channel, model models.Channel, p *gorm.DB) {
	for _, data := range feed.Categories {
		var category models.Category

		p.Table(models.Category{}.TableName()).
			Where("name = ?", data.Name).
			FirstOrCreate(&category)

		p.Table(models.Categoriable{}.TableName()).
			Where("channel_id = ? AND category_id = ?", model.Id, category.Id).
			FirstOrCreate(&models.Categoriable{})
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func random(slice []string) string {
	rand.Seed(time.Now().Unix())

	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice[0]
}

func unique(s []string) []string {
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
