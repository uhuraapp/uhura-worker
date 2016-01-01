package setup

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"gopkg.in/redis.v3"

	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/jrallison/go-workers"
	"github.com/stvp/rollbar"
)

func Worker(_redisURL string, runner bool) {
	pool := "1"
	if runner {
		pool = "15"
	}
	redisURL, err := url.Parse(_redisURL)

	if err != nil {
		panic("REDIS_URL error, " + err.Error())
	}

	rollbar.Token = os.Getenv("ROLLBAR_KEY")
	rollbar.Environment = os.Getenv("ROLLBAR_ENV")

	var password string
	if redisURL.User != nil {
		password, _ = redisURL.User.Password()
	}

	workers.Configure(map[string]string{
		"server":   redisURL.Host,
		"password": password,
		"database": "0",
		"pool":     pool,
		"process":  "1",
	})

	if runner {
		client := redis.NewClient(&redis.Options{
			Addr:     redisURL.Host,
			Password: password, // no password set
			DB:       0,        // use default DB
		})

		pong, err := client.Ping().Result()
		fmt.Println(pong, err)

		client.FlushDb()
		client.Close()

		workers.Process("duplicate-episodes", duplicateEpisodes, 1)
		workers.Process("delete-episode", deleteEpisode, 2)
		workers.Process("sync-low", syncLow, 7)
		workers.Process("sync", sync, 7)
		// workers.Process("orphan-channel", orphanChannel(p), 2)
		workers.Process("recommendations", recommendations, 1)

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
		workers.Enqueue("recommendations", "recommendations", nil)
		workers.Run()
	}
}

func reporter(message *workers.Msg) {
	if r := recover(); r != nil {
		err, _ := r.(error)
		rollbar.ErrorWithStackSkip(rollbar.ERR, err, 5, &rollbar.Field{Name: "message", Data: message.ToJson()})
		rollbar.Wait()
		panic(r)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
