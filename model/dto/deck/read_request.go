package deck

type ReadRequest struct {
	Id    int64  `json:"id" binding:"required"`
	Token string `json:"token"`
}
