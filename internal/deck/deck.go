package deck

import "time"

type Deck struct {
	ID          int64     `json:"deck_id"`
	Owner       int64     `json:"owner"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Visible     *bool     `json:"visible"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}
