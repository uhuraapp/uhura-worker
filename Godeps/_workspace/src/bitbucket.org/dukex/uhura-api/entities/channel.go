package entities

import "time"

type Channel struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	Id          int64     `json:"raw_id"`
	Uri         string    `json:"id"`
	Subscribed  bool      `json:"subscribed"`
	Copyright   string    `json:"copyright"`
	Episodes    []int64   `json:"episodes" sql:"-"`
	UpdatedAt   time.Time `json:"updated_at"`
}
