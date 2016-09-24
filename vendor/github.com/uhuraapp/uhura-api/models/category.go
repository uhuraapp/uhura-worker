package models

import "time"

type Category struct {
  Id        int64
  Name      string
  CreatedAt time.Time
  UpdatedAt time.Time
}

func (_ Category) TableName() string {
  return "categories"
}

type Categoriable struct {
  ChannelId  int64
  CategoryId int64
}

func (_ Categoriable) TableName() string {
  return "channel_categories"
}
