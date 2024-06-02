package deck

import "time"

type Details struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	PicID          string    `json:"pic_id"`
	IsSubscribed   bool      `json:"is_subscribed,omitempty"`
	Subscriptions  int       `json:"subscriptions,omitempty"`
	IsVisible      bool      `json:"is_visible,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	OwnerID        int64     `json:"owner_id,omitempty"`
	Owner          string    `json:"owner,omitempty"`
	Cards          int64     `json:"cards,omitempty"`
	TotalProgress  float32   `json:"total_progress,omitempty"`
	CardsRevised   int64     `json:"cards_revised,omitempty"`
	CardsRemaining int64     `json:"cards_remaining,omitempty"`
}
