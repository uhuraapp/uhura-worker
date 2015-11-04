package entities

type Profile struct {
	Key      string   `json:"id"`
	UserID   string   `json:"-"`
	Channels []string `json:"channels"`
}
