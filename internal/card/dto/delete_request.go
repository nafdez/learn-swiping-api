package card

type DeleteRequest struct {
	CardID int64
	DeckID int64 // Provided in GET Params
}
