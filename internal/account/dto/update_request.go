package account

import "mime/multipart"

type UpdateRequest struct {
	Token    string
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Img      *multipart.FileHeader
}
