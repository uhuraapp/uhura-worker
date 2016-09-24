package helpers

import (
	"github.com/jinzhu/gorm"
	"github.com/uhuraapp/uhura-api/entities"
)

func UserSubscriptions(userID string, db *gorm.DB, subscriptionsTableName, channelsTableName, profileID string) ([]entities.Subscription, []string) {
	var ids []int
	uids := make([]string, 0)
	subscriptions := make([]entities.Subscription, 0)

	db.Table(subscriptionsTableName).Where("user_id = ?", userID).
		Order("channel_id").
		Pluck("channel_id", &ids)

	if len(ids) > 0 {
		db.Table(channelsTableName).Where("id in (?)", ids).Order("title ASC").Find(&subscriptions)
	}

	for i := range subscriptions {
		subscriptions[i].Uri = MakeUri(subscriptions[i].Title)
		//	go subscriptions[i].SetSubscribed(userId)
		//	subscriptions[i].SetEpisodesIds()
		//	subscriptions[i].ToView = subscriptions[i].GetToView(s.DB, userId)
		uids = append(uids, subscriptions[i].Uri)
	}

	return subscriptions, uids
}
