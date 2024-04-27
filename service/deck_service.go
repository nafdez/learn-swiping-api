package service

import (
	"learn-swiping-api/erro"
	"learn-swiping-api/model"
	"learn-swiping-api/model/dto/deck"
	"learn-swiping-api/repository"
	"time"
)

type DeckService interface {
	Create(deck.CreateRequest) (int64, error)
	Deck(deck.ReadOneRequest) (model.Deck, error)
	OwnedDecks(deck.ReadOwnedRequest) ([]model.Deck, error)
	Suscriptions(deck.ReadRequest) ([]model.Deck, error) // Should be only accepting ID but for the sake of consistency
	Update(deck.UpdateRequest) error
	Delete(deck.DeleteRequest) error
	AddDeckSubscription(deck.DeckSuscriptionRequest) error
	RemoveDeckSubscription(deck.DeckSuscriptionRequest) error
}

type DeckServiceImpl struct {
	repository repository.DeckRepository
}

func NewDeckService(repository repository.DeckRepository) DeckService {
	return &DeckServiceImpl{repository: repository}
}

func (s *DeckServiceImpl) Create(request deck.CreateRequest) (int64, error) {
	deck := model.Deck{
		Owner:       request.Owner,
		Title:       request.Title,
		Description: request.Description,
		Visible:     &request.Visible,
	}

	return s.repository.Create(deck)
}

// Wondering what kind of mistakes I have made in my life to be doing this stuff
func (s *DeckServiceImpl) Deck(request deck.ReadOneRequest) (model.Deck, error) {
	if request.DeckID == 0 {
		return model.Deck{}, erro.ErrBadField
	}

	return s.repository.ById(request.DeckID, request.Token)
}

func (s *DeckServiceImpl) OwnedDecks(request deck.ReadOwnedRequest) ([]model.Deck, error) {
	if request.AccID == 0 && request.Username == "" {
		return []model.Deck{}, erro.ErrBadField
	}

	return s.repository.ByOwner(request.AccID, request.Username, request.Token)
}

func (s *DeckServiceImpl) Suscriptions(request deck.ReadRequest) ([]model.Deck, error) {
	if request.Username == "" {
		return []model.Deck{}, erro.ErrBadField
	}

	return s.repository.BySubsUsername(request.Username, request.Token)
}

func (s *DeckServiceImpl) Update(request deck.UpdateRequest) error {
	if request.Title == "" && request.Description == "" && request.Visible == nil {
		return erro.ErrBadField
	}

	// TODO: Only update if requested deck owner matched with the token provided
	deck := model.Deck{
		Title:       request.Title,
		Description: request.Description,
		Visible:     request.Visible,
		UpdatedAt:   time.Now(),
	}

	if s.repository.CheckOwnership(request.DeckID, request.Token) {
		return s.repository.Update(request.DeckID, deck)
	}

	return erro.ErrInvalidToken
}

func (s *DeckServiceImpl) Delete(request deck.DeleteRequest) error {
	if s.repository.CheckOwnership(request.DeckID, request.Token) {
		return s.repository.Delete(request.DeckID)
	}
	return erro.ErrInvalidToken
}

func (s *DeckServiceImpl) AddDeckSubscription(request deck.DeckSuscriptionRequest) error {
	return s.repository.AddDeckSubscription(request.UserID, request.DeckID)
}

func (s *DeckServiceImpl) RemoveDeckSubscription(request deck.DeckSuscriptionRequest) error {
	return s.repository.RemoveDeckSubscription(request.UserID, request.DeckID)
}
