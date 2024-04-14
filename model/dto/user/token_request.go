package user

type TokenRequest struct {
	Token string `json:"token" binding:"required"`
}
