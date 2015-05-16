package main

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/dukex/uhura-api/channels"
	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/helpers"
	"bitbucket.org/dukex/uhura-api/models"
	"bitbucket.org/dukex/uhura-api/parser"

	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
	"github.com/stvp/rollbar"
)

const imageHost string = "http://arcane-forest-5063.herokuapp.com"
const weekHours time.Duration = 24 * 7

var p gorm.DB

type rollbarMidleware struct{}

func (r *rollbarMidleware) Call(queue string, message *workers.Msg, next func() bool) (acknowledge bool) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			rollbar.ErrorWithStackSkip(rollbar.ERR, err, 5, &rollbar.Field{Name: "message", Data: message.ToJson()})
			rollbar.Wait()
			panic(r)
		}
	}()
	acknowledge = next()
	return
}

func main() {
	redis, err := url.Parse(os.Getenv("REDIS_URL"))

	if err != nil {
		panic("REDIS_URL error, " + err.Error())
	}

	rollbar.Token = os.Getenv("ROLLBAR_KEY")
	rollbar.Environment = os.Getenv("ROLLBAR_ENV")

	password, _ := redis.User.Password()

	workers.Configure(map[string]string{
		"server":   redis.Host,
		"password": password,
		"database": "0",
		"pool":     "20",
		"process":  "1",
	})

	workers.Middleware.Append(&rollbarMidleware{})

	workers.Process("sync", sync, 18)
	workers.Process("sync-low", sync, 2)

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	go workers.StatsServer(port)

	p = database.NewPostgresql()
	workers.Run()
}

func sync(message *workers.Msg) {
	var channel models.Channel

	id, err := message.Args().Int64()
	checkError(err)

	// 	// [x] channel := FindChannel(xml)
	err = p.Table(models.Channel{}.TableName()).First(&channel, id).Error
	checkError(err)

	channelURL, err := url.Parse(channel.Url)
	checkError(err)

	// 	// [x] xml := ParserXML(body)
	channelsFeed, errors := parser.URL(channelURL)
	if len(errors) > 0 {
		panic(errors[0])
	}

	channelFeed := channelsFeed[0]

	//  // [x] UpdateChannel(channel, xml)
	updateChannel(&channel, channelFeed)

	// 	// [x] CacheImage(channel)
	cacheImage(&channel)

	p.Save(&channel)

	channels.CreateLinks(r(channelFeed.Links), channel.Id, p)

	//  // [x] episodes := FindOrCreateEpisodes(channel, xml)
	saveEpisodes(p, channelFeed.Episodes, channel)

	// 	// [x] GetDelayBetweenEpisodes(episodes)
	// 	// [x] SetNewRun(channel)
	scheduleNextRun(channelFeed, id)
}

func updateChannel(model *models.Channel, feed *parser.Channel) *models.Channel {
	return channels.TranslateFromFeed(model, feed)
}

func cacheImage(model *models.Channel) {
	currentImageURL := model.ImageUrl
	if strings.Contains(currentImageURL, imageHost) {
		// Check image is OK
		return
	}

	resp, err := http.Get(imageHost + "/resolve?url=" + currentImageURL)
	if err != nil {
		return
	}

	newImageURL := resp.Request.URL.String()

	if resp.StatusCode == 200 && strings.Contains(newImageURL, imageHost+"/cache") {
		model.ImageUrl = newImageURL
	}
}

func saveEpisodes(p gorm.DB, episodes []*parser.Episode, channel models.Channel) {
	for _, episode := range episodes {
		saveEpisode(p, episode, channel)
	}
}

func saveEpisode(p gorm.DB, episodeFeed *parser.Episode, channel models.Channel) {
	var episode models.Episode

	description := episodeFeed.Summary
	if episodeFeed.Summary == "" {
		description = episodeFeed.Description
	}

	publishedAt, err := episodeFeed.Feed.ParsedPubDate()
	if err != nil {
		publishedAt, err = getPubDate(episodeFeed)
		checkError(err)
	}

	audioData := &channels.EpisodeAudioData{
		ContentLength: episodeFeed.Enclosures[0].Length,
		ContentType:   episodeFeed.Enclosures[0].Type,
	}

	if audioData.ContentLength == 0 || audioData.ContentType == "" {
		audioData, err = channels.GetEpisodeAudioData(episodeFeed.Source)
		checkError(err)
	}

	err = p.Table(models.Episode{}.TableName()).
		Where("key = ?", episodeFeed.GetKey()).
		Assign(models.Episode{
		Key:           episodeFeed.GetKey(),
		Uri:           helpers.MakeUri(episodeFeed.Title),
		Title:         episodeFeed.Title,
		SourceUrl:     episodeFeed.Source,
		Description:   description,
		ChannelId:     channel.Id,
		PublishedAt:   publishedAt,
		ContentType:   audioData.ContentType,
		ContentLength: audioData.ContentLength,
	}).
		FirstOrCreate(&episode).Error
	checkError(err)
}

func scheduleNextRun(channelFeed *parser.Channel, id int64) {
	if len(channelFeed.Episodes) > 1 {
		t0, errT1 := channelFeed.Episodes[0].Feed.ParsedPubDate()
		if errT1 != nil {
			scheduleNextWeek(id)
			return
		}

		t1, errT2 := channelFeed.Episodes[1].Feed.ParsedPubDate()
		if errT2 != nil {
			scheduleNextWeek(id)
			return
		}

		nextRunAt := t0.Add(t0.Sub(t1))

		now := time.Now()
		if !nextRunAt.After(now) {
			scheduleNextWeek(id)
			return
		}

		scheduleAt(nextRunAt, id)
	} else {
		scheduleNextWeek(id)
		return
	}
}

func scheduleAt(at time.Time, id int64) {
	workers.EnqueueAt("sync", "sync", at, id)
}

func scheduleNextWeek(id int64) {
	now := time.Now()
	at := now.Add(time.Hour * weekHours)
	scheduleAt(at, id)
}

const episodePubDateFormat = "Mon, _2 Jan 2006 15:04:05 -0700"

func getPubDate(e *parser.Episode) (time.Time, error) {
	pubDate := strings.Replace(e.PubDate, "GMT", "-0100", -1)
	pubDate = strings.Replace(pubDate, "PST", "-0800", -1)
	pubDate = strings.Replace(pubDate, "PDT", "-0700", -1)

	return time.Parse(episodePubDateFormat, pubDate)
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
