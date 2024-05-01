package config

import (
	"database/sql"
	"learn-swiping-api/internal/account"
	"learn-swiping-api/internal/card"
	"learn-swiping-api/internal/deck"
	"learn-swiping-api/internal/picture"
)

type Initialization struct {
	UserCtrl account.AccountController
	// userSrvc service.UserService
	// userRepo repository.UserRepository
	DeckCtrl deck.DeckController
	// deckSrvc service.DeckService
	// deckRepo repository.DeckRepository
	CardCtrl card.CardController
	// cardSrvc service.CardService
	// cardRepo repository.CardRepository
	PictureCtrl picture.PictureController
}

func NewInitialization(db *sql.DB) *Initialization {
	userRepo := account.NewAccountRepository(db)
	userSrvc := account.NewAccountService(userRepo)
	userCtrl := account.NewAccountController(userSrvc)

	deckRepo := deck.NewDeckRepository(db)
	deckSrvc := deck.NewDeckService(deckRepo)
	deckCtrl := deck.NewDeckController(deckSrvc)

	cardRepo := card.NewCardRepository(db)
	cardSrvc := card.NewCardService(cardRepo)
	cardCtrl := card.NewCardController(cardSrvc)

	pictureCtrl := picture.NewPictureController()

	return &Initialization{
		// userRepo: userRepo,
		// userSrvc: userSrvc,
		UserCtrl: userCtrl,
		// deckRepo: deckRepo,
		// deckSrvc: deckSrvc,
		DeckCtrl:    deckCtrl,
		CardCtrl:    cardCtrl,
		PictureCtrl: pictureCtrl,
	}
}
