package model

type Card struct {
	CardID   int64         `json:"card_id"`
	DeckID   int64         `json:"deck_id"`
	Study    string        `json:"study"`
	Question string        `json:"question"`
	Answer   string        `json:"answer"`
	Wrong    []WrongAnswer `json:"wrong,omitempty"`
}

type WrongAnswer struct {
	WrongID int64  `json:"wrong_id,omitempty"`
	CardID  int64  `json:"card_id,omitempty"`
	Answer  string `json:"answer" binding:"required"` // required for create request

}
