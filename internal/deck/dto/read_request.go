package deck

// Suscriptions request
type ReadRequest struct {
	Username string `json:"username"`
}

// Owned decks request
type ReadOwnedRequest struct {
	AccID    int64  `json:"acc_id"`
	Username string `json:"username"`
}

// Only one request
type ReadOneRequest struct {
	DeckID int64 `json:"deck_id"`
}
