package model

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type Ticket struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
