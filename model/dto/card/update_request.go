package card

type UpdateRequest struct {
	ID       int64                `json:"card_id"`
	Study    string               `json:"study"`
	Question string               `json:"question"`
	Answer   string               `json:"answer"`
	Wrong    []updateWrongRequest `json:"wrong"`
}

type updateWrongRequest struct {
	WrongID int64  `json:"wrong_id"`
	Answer  string `json:"answer"`
}
