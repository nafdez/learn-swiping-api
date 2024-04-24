package card

type DeleteRequest struct {
	Id int64 `json:"card_id" binding:"required"`
}
