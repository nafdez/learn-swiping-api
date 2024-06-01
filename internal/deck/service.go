package deck

import (
	"bytes"
	"fmt"
	"io"
	"learn-swiping-api/erro"
	deck "learn-swiping-api/internal/deck/dto"
	"learn-swiping-api/internal/picture"
	"path/filepath"
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
	request.PicID = "default_deck_pic_1.png"

	deckID, err := s.repository.Create(request)
	if err != nil {
		return 0, err
	}

	// TODO: Check if error and rollback
	s.repository.AddDeckSubscription(request.Token, deckID)

	return 0, nil
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
	// If all fields are empty, throw an error
	if request.Title == "" && request.Description == "" && request.Visible == nil && request.Img == nil {
		return erro.ErrBadField
	}

	deck := Deck{
		Title:       request.Title,
		Description: request.Description,
		Visible:     request.Visible,
		UpdatedAt:   time.Now(),
	}

	// Check if image file isn't empty, stores it
	// and then binds the PicID to the user
	if request.Img != nil {
		// Necesary to remove the previous pic
		fmt.Println(request.DeckID, token)
		oldDeck, err := s.repository.ById(request.DeckID, token)
		if err != nil {
			return err
		}

		img, err := request.Img.Open()
		if err != nil {
			return erro.ErrBadField
		}
		defer img.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, img); err != nil {
			return erro.ErrBadField
		}

		picture.Remove(oldDeck.PicID)
		picID, err := picture.Store(filepath.Ext(request.Img.Filename), buf.Bytes())
		if err != nil {
			return err
		}
		deck.PicID = picID
	}

	// TODO: Only update if requested deck owner matched with the token provided directly into the update query if possible
	if s.repository.CheckOwnership(request.DeckID, token) {
		return s.repository.Update(request.DeckID, deck)
	}

	return erro.ErrInvalidToken
}

func (s *DeckServiceImpl) Delete(request deck.DeleteRequest, token string) error {
	// Doesn't work as intended. revisar
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
