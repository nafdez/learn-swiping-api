package user

type UpdateRequest struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}
