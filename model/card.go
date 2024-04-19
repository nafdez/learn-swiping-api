package model

type Card struct {
	ID       int64         `json:"card_id"`
	DeckID   int64         `json:"deck_id"`
	Study    string        `json:"study"`
	Question string        `json:"question"`
	Answer   string        `json:"answer"`
	Wrong    []WrongAnswer `json:"wrong"`
}

type WrongAnswer struct {
	ID     int64  `json:"wrong_id"`
	CardID int64  `json:"card_id"`
	Answer string `json:"answer"`
}
