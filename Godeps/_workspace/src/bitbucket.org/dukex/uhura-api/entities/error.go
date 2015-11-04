package entities

type Error struct {
	Key     string `json:"key"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}
