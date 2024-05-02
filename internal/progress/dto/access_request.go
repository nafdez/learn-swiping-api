package progress

type AccessRequest struct {
	Token  string
	CardID int64 `json:"card_id" binding:"required"`
}

func (req *AccessRequest) SetToken(token string) {
	req.Token = token
}
