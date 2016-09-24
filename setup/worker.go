package setup

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"gopkg.in/redis.v3"

	"github.com/jrallison/go-workers"
	"github.com/stvp/rollbar"
	"github.com/uhuraapp/uhura-api/database"
	"github.com/uhuraapp/uhura-api/models"
)

func Worker(_redisURL string) {
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
		"pool":     "15",
		"process":  "1",
	})

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	client.FlushDb()
	client.Close()

	workers.Process("sync", sync, 7)

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
