package entities

import "time"

type User struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Locale     string    `json:"locale"`
	Email      string    `json:"email"`
	ProviderId string    `json:"provider_id"`
	ApiToken   string    `json:"token"`
	ProfileKey string    `json:"profile_key"`
	OptIn      bool      `json:"optin"`
	OptInAt    time.Time `json:"-"`
	DeletedAt  time.Time `json:"-"`
}
