package deck

type DeckSuscriptionRequest struct {
	UserID int64 `json:"acc_id" binding:"required"`
	DeckID int64 `json:"deck_id" binding:"required"`
}
