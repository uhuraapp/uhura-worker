package entities

import "time"

type Subscription struct {
	Id          int64     `json:"raw_id"`
	Uri         string    `json:"id"`
	Title       string    `json:"title"`
	ImageUrl    string    `json:"image_url"`
	ToView      int64     `json:"to_view"`
	Subscribed  bool      `json:"subscribed"`
	ProfileID   string    `json:"profile_id"`
	Description string    `json:"description"`
	Copyright   string    `json:"copyright"`
	UpdatedAt   time.Time `json:"updated_at"`
	Enabled     bool      `json:"enabled"`
}

// func (s *Subscription) GetToView(database gorm.DB, userId string) int64 {
// 	var (
// 		listenedCount int64
// 	)
//
// 	database.Table(models.Listened{}.TableName()).
// 		Where("channel_id = ? AND user_id = ?", s.Id, userId).
// 		Count(&listenedCount)
//
// 	return models.Episode{}.CountByChannel(database, s.Id) - listenedCount
// }
