package config

import (
	"database/sql"
	"learn-swiping-api/controller"
	"learn-swiping-api/repository"
	"learn-swiping-api/service"
)

type Initialization struct {
	UserCtrl controller.UserController
	userSrvc service.UserService
	userRepo repository.UserRepository
}

func NewInitialization(db *sql.DB) *Initialization {
	userRepo := repository.NewUserRepository(db)
	userSrvc := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userSrvc)

	return &Initialization{
		userRepo: userRepo,
		userSrvc: userSrvc,
		UserCtrl: userCtrl,
	}
}
