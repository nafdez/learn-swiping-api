package card

type UpdateRequest struct {
	DeckID   int64                // Provided in GET params
	CardID   int64                `json:"card_id"`
	Front    string               `json:"front"`
	Back     string               `json:"back"`
	Question string               `json:"question"`
	Answer   string               `json:"answer"`
	Wrong    []updateWrongRequest `json:"wrong"`
}

type updateWrongRequest struct {
	WrongID int64  `json:"wrong_id" binding:"required"`
	Answer  string `json:"answer" binding:"required"` // Nothing to update if id is provided but not the new answer
}
