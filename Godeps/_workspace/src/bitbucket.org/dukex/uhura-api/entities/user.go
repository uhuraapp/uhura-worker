package entities

type User struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Picture    string `json:"image_url"`
	Locale     string `json:"locale"`
	Email      string `json:"email"`
	ProviderId string `json:"provider_id"`
	ApiToken   string `json:"token"`
}
