package deck

type DeckSuscriptionRequest struct {
	Token  string // acc_id gathered with token
	DeckID int64  `json:"deck_id" binding:"required"`
}
