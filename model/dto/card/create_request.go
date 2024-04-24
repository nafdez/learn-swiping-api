package card

import "learn-swiping-api/model"

type CreateRequest struct {
	DeckID   int64               `json:"deck_id" binding:"required"`
	Study    string              `json:"study" binding:"required"`
	Question string              `json:"question" binding:"required"`
	Answer   string              `json:"answer" binding:"required"`
	Wrong    []model.WrongAnswer `json:"wrong" binding:"required"`
}
