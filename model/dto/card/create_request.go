package card

import "learn-swiping-api/model"

type CreateRequest struct {
	DeckID   int64               // Providen in GET params
	Front    string              `json:"front" binding:"required"`
	Back     string              `json:"back" binding:"required"`
	Question string              `json:"question" binding:"required"`
	Answer   string              `json:"answer" binding:"required"`
	Wrong    []model.WrongAnswer `json:"wrong" binding:"required"`
}
