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
	AllFromDeck(deckID int64) ([]model.Card, error)
	ReadOne(id int64) (model.Card, error)
	Update(card.UpdateRequest) error
	Delete(id int64) error
}

type CardServiceImpl struct {
	repository repository.CardRepository
}

func NewCardService(repository repository.CardRepository) CardService {
	return &CardServiceImpl{repository: repository}
}

func (s *CardServiceImpl) Create(request card.CreateRequest) (int64, error) {
	if request.DeckID == 0 || request.Study == "" || request.Question == "" || request.Answer == "" || len(request.Wrong) != 3 {
		return 0, erro.ErrBadField
	}

	card := model.Card{
		DeckID:   request.DeckID,
		Study:    request.Study,
		Question: request.Question,
		Answer:   request.Answer,
	}

	return s.repository.Create(card, request.Wrong)
}

func (s *CardServiceImpl) AllFromDeck(deckID int64) ([]model.Card, error) {
	if deckID == 0 {
		return nil, erro.ErrCardNotFound
	}

	// Only getting all results from the query directly
	// without adding the wrong answers to the return
	// so we don't retrieve lots of rows at once.
	// Wrong answers should only be needed when viewing one
	// card at most
	return s.repository.ByDeckId(deckID)
}

func (s *CardServiceImpl) ReadOne(id int64) (model.Card, error) {
	if id == 0 {
		return model.Card{}, erro.ErrCardNotFound
	}

	card, err := s.repository.ById(id)
	if err != nil {
		return model.Card{}, err
	}

	card.Wrong, err = s.repository.WrongById(id)
	if err != nil {
		return model.Card{}, err
	}

	return card, nil
}

func (s *CardServiceImpl) Update(request card.UpdateRequest) error {
	if request.ID == 0 {
		return erro.ErrCardNotFound
	}

	if request.Study != "" || request.Question != "" || request.Answer != "" {
		card := model.Card{
			Study:    request.Study,
			Question: request.Question,
			Answer:   request.Answer,
		}

		err := s.repository.Update(request.ID, card)
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

func (s *CardServiceImpl) Delete(id int64) error {
	if id != 0 {
		return erro.ErrBadField
	}
	return s.repository.Delete(id)
}
