package deck

type DeleteRequest struct {
	Token  string `json:"token" binding:"required"`
	DeckID int64  `json:"deck_id" binding:"required"`
}
