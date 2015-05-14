package channels

import (
	"net/url"
	"os"

	"bitbucket.org/dukex/uhura-api/models"
	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
)

func Ping(database gorm.DB, id int64) {
	var channel models.Channel

	ok := database.Table(models.Channel{}.TableName()).
		First(&channel, id).Error != gorm.RecordNotFound

	if !ok {
		return
	}

	workers.Enqueue("sync-low", "sync", id)
}

func init() {
	redis, err := url.Parse(os.Getenv("REDIS_URL"))

	if err != nil {
		panic("REDIS_URL error, " + err.Error())
	}

	password, _ := redis.User.Password()

	workers.Configure(map[string]string{
		"server":   redis.Host,
		"password": password,
		"database": "0",
		"pool":     "1",
		"process":  "1",
	})
}
