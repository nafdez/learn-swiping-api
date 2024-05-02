package config

import (
	"database/sql"
	"learn-swiping-api/internal/account"
	"learn-swiping-api/internal/card"
	"learn-swiping-api/internal/deck"
	"learn-swiping-api/internal/picture"
	"learn-swiping-api/internal/progress"
)

type Initialization struct {
	UserCtrl     account.AccountController
	DeckCtrl     deck.DeckController
	CardCtrl     card.CardController
	ProgressCtrl progress.ProgressController
	PictureCtrl  picture.PictureController
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

	progressRepo := progress.NewProgressRepository(db)
	progressSrvc := progress.NewProgressService(progressRepo)
	progressCtrl := progress.NewProgressController(progressSrvc)

	pictureCtrl := picture.NewPictureController()

	return &Initialization{
		UserCtrl:     userCtrl,
		DeckCtrl:     deckCtrl,
		CardCtrl:     cardCtrl,
		ProgressCtrl: progressCtrl,
		PictureCtrl:  pictureCtrl,
	}
}
