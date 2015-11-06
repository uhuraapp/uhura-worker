package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"gopkg.in/redis.v3"

	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/models"
	duplicates "bitbucket.org/dukex/uhura-worker/duplicates"
	syncRunner "bitbucket.org/dukex/uhura-worker/sync"

	"github.com/dukex/too"
	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
	"github.com/stvp/rollbar"
)

var p gorm.DB

func main() {
	redisURL, err := url.Parse(os.Getenv("REDIS_URL"))

	if err != nil {
		panic("REDIS_URL error, " + err.Error())
	}

	rollbar.Token = os.Getenv("ROLLBAR_KEY")
	rollbar.Environment = os.Getenv("ROLLBAR_ENV")

	password, _ := redisURL.User.Password()

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	client.FlushDb()
	client.Close()

	p = database.NewPostgresql()

	workers.Configure(map[string]string{
		"server":   redisURL.Host,
		"password": password,
		"database": "0",
		"pool":     "20",
		"process":  "1",
	})

	// heroku support 20 connections
	workers.Process("sync", sync, 10)
	workers.Process("sync-low", syncLow, 3)
	workers.Process("duplicate-episodes", duplicateEpisodes, 1)
	workers.Process("orphan-channels", orphanChannel, 1)
	workers.Process("delete-episode", deleteEpisode, 2)
	workers.Process("recommendations", recommendations(p), 2)

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	go workers.StatsServer(port)

	go func() {
		workers.Enqueue("orphan-channels", "orphanChannel", nil)

		var c []int64
		p.Table(models.Channel{}.TableName()).Pluck("id", &c)
		for _, id := range c {
			workers.Enqueue("sync-low", "sync", id)
		}
	}()

	workers.Enqueue("recommendations", "recommendations", nil)
	workers.Run()
}

func reporter(message *workers.Msg) {
	if r := recover(); r != nil {
		err, _ := r.(error)
		rollbar.ErrorWithStackSkip(rollbar.ERR, err, 5, &rollbar.Field{Name: "message", Data: message.ToJson()})
		rollbar.Wait()
		panic(r)
	}
}

func recommendations(db gorm.DB) func(*workers.Msg) {
	return func(message *workers.Msg) {
		defer reporter(message)
		var users []int64

		p.Table(models.User{}.TableName()).Pluck("id", &users)

		te, err := too.New(os.Getenv("RREDIS_URL"), "channels")
		checkError(err)

		for i := 0; i < len(users); i++ {
			var subscriptions []models.Subscription
			p.Table(models.Subscription{}.TableName()).Where("user_id = ?", users[i]).Find(&subscriptions)
			userID := too.User(strconv.Itoa(int(users[i])))

			for j := 0; j < len(subscriptions); j++ {
				channelID := too.Item(strconv.Itoa(int(subscriptions[j].ChannelId)))
				te.Likes.Add(userID, channelID)
			}
		}

		workers.EnqueueAt("recommendations", "recommendations", time.Now().Add(5*time.Hour), nil)
	}
}

func syncLow(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	s := syncRunner.NewSync(id)
	s.Sync(p)

	workers.EnqueueAt("sync-low", "sync", time.Now().Add(5*time.Minute), id)
	workers.Enqueue("duplicate-episodes", "duplicateEpisodes", nil)
	workers.Enqueue("orphan-channel", "orphanChannel", nil)
}

func sync(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	s := syncRunner.NewSync(id)
	s.Sync(p)

	nextRunAt, err := s.GetNextRun()
	checkError(err)

	workers.EnqueueAt("sync", "sync", nextRunAt, id)
	workers.Enqueue("duplicate-episodes", "duplicateEpisodes", nil)
	workers.Enqueue("orphan-channel", "orphanChannel", nil)
}

func duplicateEpisodes(message *workers.Msg) {
	defer reporter(message)

	episodes := duplicates.Episodes(p)
	for _, id := range episodes {
		workers.Enqueue("delete-episode", "deleteEpisode", id)
	}
}

func deleteEpisode(message *workers.Msg) {
	defer reporter(message)

	id, err := message.Args().Int64()
	checkError(err)

	p.Table(models.Episode{}.TableName()).Where("id = ?", id).Delete(models.Episode{})
}

func orphanChannel(message *workers.Msg) {
	defer reporter(message)

	var channels []models.Channel
	p.Table(models.Channel{}.TableName()).Find(&channels)

	for _, channel := range channels {
		var users []models.Subscription
		p.Table(models.Subscription{}.TableName()).Where("channel_id = ?", channel.Id).Find(&users)
		if len(users) < 1 {
			var episodes []models.Episode
			p.Table(models.Episode{}.TableName()).Where("channel_id = ?", channel.Id).Find(&episodes)
			for _, e := range episodes {
				workers.Enqueue("delete-episode", "deleteEpisode", e.Id)
			}
			p.Table(models.Channel{}.TableName()).Where("id = ?", channel.Id).Delete(models.Channel{})
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
