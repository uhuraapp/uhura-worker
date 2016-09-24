package models

import "time"

type Subscription struct {
  Id        int64
  UserId    int64
  ChannelId int64
  CreatedAt time.Time
  UpdatedAt time.Time
}

func (s Subscription) TableName() string {
  return "user_channels"
}
