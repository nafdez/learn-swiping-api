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
	Read(deck.ReadRequest) (model.Deck, error)
	ReadFromUser(deck.ReadRequest) ([]model.Deck, error)
	ReadUserSuscriptions(deck.ReadRequest) ([]model.Deck, error) // Should be only accepting ID but for the sake of consistency
	Update(deck.UpdateRequest) error
	Delete(deck.DeleteRequest) error
}

type DeckServiceImpl struct {
	repository repository.DeckRepository
}

func NewDeckService(repository repository.DeckRepository) DeckService {
	return &DeckServiceImpl{repository: repository}
}

func (s *DeckServiceImpl) Create(request deck.CreateRequest) (int64, error) {
	if request.Owner == 0 || request.Title == "" {
		return 0, erro.ErrBadField
	}

	deck := model.Deck{
		Owner:       request.Owner,
		Title:       request.Title,
		Description: request.Description,
		Visible:     &request.Visible,
	}

	return s.repository.Create(deck)
}

// Wondering what kind of mistakes I have made in my life to be doing this stuff
func (s *DeckServiceImpl) Read(request deck.ReadRequest) (model.Deck, error) {
	return s.repository.ById(request)
}

func (s *DeckServiceImpl) ReadFromUser(request deck.ReadRequest) ([]model.Deck, error) {
	return s.repository.ByOwner(request)
}

func (s *DeckServiceImpl) ReadUserSuscriptions(request deck.ReadRequest) ([]model.Deck, error) {
	return s.repository.ByUserId(request)
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
