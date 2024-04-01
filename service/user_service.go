package service

import (
	"learn-swiping-api/model"
	"learn-swiping-api/model/dto/user"
)

type UserService interface {
	Register(user.RegisterRequest) (*model.User, error)
	Login(user.LoginRequest) (*model.User, error)
	Token(user.TokenRequest) (*model.User, error) // Login with token
	Logout(user.TokenRequest) error
	Account(user.TokenRequest) (*model.User, error)
	User(username string) (*user.Public, error)
	Update(user.UpdateRequest) error
	Delete(user.TokenRequest) error
}
