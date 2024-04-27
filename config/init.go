package config

import (
	"database/sql"
	"learn-swiping-api/controller"
	"learn-swiping-api/repository"
	"learn-swiping-api/service"
)

type Initialization struct {
	UserCtrl controller.UserController
	// userSrvc service.UserService
	// userRepo repository.UserRepository
	DeckCtrl controller.DeckController
	// deckSrvc service.DeckService
	// deckRepo repository.DeckRepository
	CardCtrl controller.CardController
	// cardSrvc service.CardService
	// cardRepo repository.CardRepository
}

func NewInitialization(db *sql.DB) *Initialization {
	userRepo := repository.NewUserRepository(db)
	userSrvc := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userSrvc)

	deckRepo := repository.NewDeckRepository(db)
	deckSrvc := service.NewDeckService(deckRepo)
	deckCtrl := controller.NewDeckController(deckSrvc)

	cardRepo := repository.NewCardRepository(db)
	cardSrvc := service.NewCardService(cardRepo)
	cardCtrl := controller.NewCardController(cardSrvc)

	return &Initialization{
		// userRepo: userRepo,
		// userSrvc: userSrvc,
		UserCtrl: userCtrl,
		// deckRepo: deckRepo,
		// deckSrvc: deckSrvc,
		DeckCtrl: deckCtrl,
		CardCtrl: cardCtrl,
	}
}
