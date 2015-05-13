package models

import "time"

type Listened struct {
	Id        int64
	UserId    int64
	ItemId    int64
	Viewed    bool
	ChannelId int64
	StoppedAt int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (l Listened) TableName() string {
	return "user_items"
}
