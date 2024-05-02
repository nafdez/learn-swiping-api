package deck

import (
	"learn-swiping-api/erro"
	deck "learn-swiping-api/internal/deck/dto"
	"time"
)

type DeckService interface {
	Create(deck.CreateRequest) (int64, error)
	Deck(req deck.ReadOneRequest, token string) (Deck, error)
	OwnedDecks(req deck.ReadOwnedRequest, token string) ([]Deck, error)
	Suscriptions(req deck.ReadRequest, token string) ([]Deck, error) // Should be only accepting ID but for the sake of consistency
	Update(req deck.UpdateRequest, token string) error
	Delete(req deck.DeleteRequest, token string) error
	AddDeckSubscription(deck.DeckSuscriptionRequest) error
	RemoveDeckSubscription(deck.DeckSuscriptionRequest) error
}

type DeckServiceImpl struct {
	repository DeckRepository
}

func NewDeckService(repository DeckRepository) DeckService {
	return &DeckServiceImpl{repository: repository}
}

func (s *DeckServiceImpl) Create(request deck.CreateRequest) (int64, error) {
	return s.repository.Create(request)
}

// Wondering what kind of mistakes I have made in my life to be doing this stuff
func (s *DeckServiceImpl) Deck(request deck.ReadOneRequest, token string) (Deck, error) {
	if request.DeckID == 0 {
		return Deck{}, erro.ErrBadField
	}

	return s.repository.ById(request.DeckID, token)
}

func (s *DeckServiceImpl) OwnedDecks(request deck.ReadOwnedRequest, token string) ([]Deck, error) {
	if request.AccID == 0 && request.Username == "" {
		return []Deck{}, erro.ErrBadField
	}

	return s.repository.ByOwner(request.AccID, request.Username, token)
}

func (s *DeckServiceImpl) Suscriptions(request deck.ReadRequest, token string) ([]Deck, error) {
	if request.Username == "" {
		return []Deck{}, erro.ErrBadField
	}

	return s.repository.BySubsUsername(request.Username, token)
}

func (s *DeckServiceImpl) Update(request deck.UpdateRequest, token string) error {
	if request.Title == "" && request.Description == "" && request.Visible == nil {
		return erro.ErrBadField
	}

	// TODO: Only update if requested deck owner matched with the token provided
	deck := Deck{
		Title:       request.Title,
		Description: request.Description,
		Visible:     request.Visible,
		UpdatedAt:   time.Now(),
	}

	if s.repository.CheckOwnership(request.DeckID, token) {
		return s.repository.Update(request.DeckID, deck)
	}

	return erro.ErrInvalidToken
}

func (s *DeckServiceImpl) Delete(request deck.DeleteRequest, token string) error {
	if s.repository.CheckOwnership(request.DeckID, token) {
		return s.repository.Delete(request.DeckID)
	}
	return erro.ErrInvalidToken
}

func (s *DeckServiceImpl) AddDeckSubscription(request deck.DeckSuscriptionRequest) error {
	return s.repository.AddDeckSubscription(request.Token, request.DeckID)
}

func (s *DeckServiceImpl) RemoveDeckSubscription(request deck.DeckSuscriptionRequest) error {
	return s.repository.RemoveDeckSubscription(request.Token, request.DeckID)
}
