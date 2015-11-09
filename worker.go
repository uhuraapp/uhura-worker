package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/models"

	"gopkg.in/redis.v3"

	"github.com/jrallison/go-workers"
	"github.com/stvp/rollbar"
)

func main() {
	redisURL, err := url.Parse(os.Getenv("REDIS_URL"))

	if err != nil {
		panic("REDIS_URL error, " + err.Error())
	}

	rollbar.Token = os.Getenv("ROLLBAR_KEY")
	rollbar.Environment = os.Getenv("ROLLBAR_ENV")

	var password string
	if redisURL.User != nil {
		password, _ = redisURL.User.Password()
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	client.FlushDb()
	client.Close()

	workers.Configure(map[string]string{
		"server":   redisURL.Host,
		"password": password,
		"database": "0",
		"pool":     "15",
		"process":  "1",
	})

	workers.Process("duplicate-episodes", duplicateEpisodes, 1)
	workers.Process("delete-episode", deleteEpisode, 2)
	workers.Process("sync-low", syncLow, 8)
	workers.Process("sync", sync, 7)
	// workers.Process("orphan-channel", orphanChannel(p), 2)
	// workers.Process("recommendations", recommendations(p), 1)

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	go workers.StatsServer(port)

	go func() {
		workers.Enqueue("orphan-channel", "orphanChannel", 0)

		var c []int64

		p := database.NewPostgresql()
		p.Table(models.Channel{}.TableName()).Pluck("id", &c)
		for _, id := range c {
			workers.Enqueue("sync", "sync", id)
		}

		p.Close()
	}()

	workers.Enqueue("duplicate-episodes", "duplicateEpisodes", nil)
	// workers.Enqueue("recommendations", "recommendations", nil)
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

// func recommendations(p gorm.DB) func(*workers.Msg) {
// 	return func(message *workers.Msg) {
// 		defer reporter(message)
// 		var users []int64
//
// 		p.Table(models.User{}.TableName()).Pluck("id", &users)
//
// 		te, err := too.New(os.Getenv("RREDIS_URL"), "channels")
// 		checkError(err)
//
// 		for i := 0; i < len(users); i++ {
// 			var subscriptions []models.Subscription
// 			p.Table(models.Subscription{}.TableName()).Where("user_id = ?", users[i]).Find(&subscriptions)
// 			userID := too.User(strconv.Itoa(int(users[i])))
// 			log.Println("too", userID)
//
// 			for j := 0; j < len(subscriptions); j++ {
// 				channelID := too.Item(strconv.Itoa(int(subscriptions[j].ChannelId)))
// 				log.Println("too", channelID)
// 				log.Println("too", te.Likes.Add(userID, channelID))
// 			}
// 			log.Println("too", "--------------------------")
// 		}
//
// 		workers.EnqueueAt("recommendations", "recommendations", time.Now().Add(12*time.Hour), nil)
// 	}
// }

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
