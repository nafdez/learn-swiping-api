package card

type DeleteRequest struct {
	Id     int64 `json:"card_id" binding:"required"`
	DeckID int64 // Provided in GET Params
}
