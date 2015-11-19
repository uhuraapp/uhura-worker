package main

import (
	"os"

	"bitbucket.org/dukex/uhura-worker/setup"
)

func main() {
	setup.Worker(os.Getenv("REDIS_URL"), true)
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
