package entities

import "time"

type Subscription struct {
	Id          int64     `json:"raw_id"`
	Uri         string    `json:"id"`
	Title       string    `json:"title"`
	ImageUrl    string    `json:"image_url"`
	ToView      int64     `json:"to_view"`
	Subscribed  bool      `json:"subscribed"`
	Description string    `json:"description"`
	Copyright   string    `json:"copyright"`
	UpdatedAt   time.Time `json:"updated_at"`
	Enabled     bool      `json:"enabled"`
}
