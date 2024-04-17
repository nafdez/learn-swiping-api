package card

import "learn-swiping-api/model"

type CreateRequest struct {
	DeckID   int64               `json:"deck_id"`
	Study    string              `json:"study"`
	Question string              `json:"question"`
	Answer   string              `json:"answer"`
	Wrong    []model.WrongAnswer `json:"wrong"`
}
