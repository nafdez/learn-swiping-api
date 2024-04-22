package user

type DeckSuscriptionRequest struct {
	UserID int64 `json:"acc_id"`
	DeckID int64 `json:"deck_id"`
}
