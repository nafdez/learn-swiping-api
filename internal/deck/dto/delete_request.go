package deck

type DeleteRequest struct {
	DeckID int64 `json:"deck_id" binding:"required"`
}
