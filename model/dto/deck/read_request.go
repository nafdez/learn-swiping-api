package deck

// Suscriptions request
type ReadRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// Owned decks request
type ReadOwnedRequest struct {
	AccID    int64  `json:"acc_id"`
	Username string `json:"username"`
	Token    string `json:"token"` // To show hidden ones
}

// Only one request
type ReadOneRequest struct {
	DeckID int64  `json:"deck_id"`
	Token  string `json:"token"`
	// Same as suscription one, but id is not required
	// In fact it is, but is gathered from GET params
	// Just included here for maintain consistency in deck_service
}
