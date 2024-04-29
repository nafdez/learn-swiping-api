package service

import (
	"errors"
	"learn-swiping-api/erro"
	"learn-swiping-api/model"
	"learn-swiping-api/model/dto/card"
	"learn-swiping-api/repository"
	"strconv"
)

type CardService interface {
	Create(card.CreateRequest) (int64, error)
	Card(cardID int64, deckID int64) (model.Card, error)
	Cards(deckID int64) ([]model.Card, error)
	Update(card.UpdateRequest) error
	Delete(cardID int64, deckID int64) error
}

type CardServiceImpl struct {
	repository repository.CardRepository
}

func NewCardService(repository repository.CardRepository) CardService {
	return &CardServiceImpl{repository: repository}
}

func (s *CardServiceImpl) Create(request card.CreateRequest) (int64, error) {
	if len(request.Wrong) != 3 { // All cards should have four answers, one OK three bad
		return 0, erro.ErrBadField
	}

	card := model.Card{
		DeckID:   request.DeckID,
		Front:    request.Front,
		Back:     request.Back,
		Question: request.Question,
		Answer:   request.Answer,
		Wrong:    request.Wrong,
	}

	return s.repository.Create(card)
}

func (s *CardServiceImpl) Card(cardID int64, deckID int64) (model.Card, error) {
	card, err := s.repository.ById(cardID, deckID)
	if err != nil {
		return model.Card{}, err
	}

	card.Wrong, err = s.repository.WrongByCardId(cardID)
	if err != nil {
		return model.Card{}, err
	}

	return card, nil
}

func (s *CardServiceImpl) Cards(deckID int64) ([]model.Card, error) {
	// Wrong answers should only be needed when viewing one
	// card at most
	return s.repository.ByDeckId(deckID)
}

func (s *CardServiceImpl) Update(request card.UpdateRequest) error {
	if request.Front != "" || request.Back != "" || request.Question != "" || request.Answer != "" {
		card := model.Card{
			CardID:   request.CardID,
			DeckID:   request.DeckID,
			Front:    request.Front,
			Back:     request.Back,
			Question: request.Question,
			Answer:   request.Answer,
		}

		err := s.repository.Update(card)
		if err != nil {
			return err
		}
	}

	// TODO: redo
	if len(request.Wrong) > 0 {
		affected := 0
		for _, answer := range request.Wrong {
			if answer.WrongID != 0 && answer.Answer != "" {
				err := s.repository.UpdateWrong(answer.WrongID, model.WrongAnswer{Answer: answer.Answer})
				if err != nil {
					return err
				}
				affected++
			}
		}
		if affected < len(request.Wrong) {
			return errors.New("some answers were not updated. updated: " + strconv.Itoa(affected))
		}
	}

	return nil
}

func (s *CardServiceImpl) Delete(cardID int64, deckID int64) error {
	return s.repository.Delete(cardID, deckID)
}
