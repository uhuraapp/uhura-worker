package setup

import (
	"log"
	"os"
	"strconv"
	"time"

	"bitbucket.org/dukex/uhura-api/database"
	"bitbucket.org/dukex/uhura-api/models"
	"github.com/hjr265/too"
	"github.com/jrallison/go-workers"
)

func recommendations(message *workers.Msg) {
	defer reporter(message)

	p := database.NewPostgresql()
	var users []int64

	p.Table(models.User{}.TableName()).Pluck("id", &users)

	te, err := too.New(os.Getenv("RREDIS_URL"), "channels")
	checkError(err)

	ops := make([]too.BatchRaterOp, 0)

	for i := 0; i < len(users); i++ {
		var subscriptions []models.Subscription
		items := make([]too.Item, 0)
		p.Table(models.Subscription{}.TableName()).Where("user_id = ?", users[i]).Find(&subscriptions)
		user := too.User(strconv.Itoa(int(users[i])))

		for j := 0; j < len(subscriptions); j++ {
			items = append(items, too.Item(strconv.Itoa(int(subscriptions[j].ChannelId))))
		}

		ops = append(ops, too.BatchRaterOp{
			User:  user,
			Items: items,
		})
		log.Println(ops)
	}

	te.Likes.Batch(ops)

	workers.EnqueueAt("recommendations", "recommendations", time.Now().Add(12*time.Hour), nil)
}
