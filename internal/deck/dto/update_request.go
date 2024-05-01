package deck

type UpdateRequest struct {
	Token       string `json:"token" binding:"required"`
	DeckID      int64  `json:"deck_id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Visible     *bool  `json:"visible"` // pointer to check if empty or not
}
