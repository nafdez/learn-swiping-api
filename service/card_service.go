package service

import (
	"learn-swiping-api/model"
	"learn-swiping-api/model/dto/card"
	"learn-swiping-api/repository"
)

type CardService interface {
	Create(card.CreateRequest) (int64, error)
	Read(card.ReadRequest) ([]model.Card, error)
	Update(card.UpdateRequest) error
	Delete(card.DeleteRequest) error
}

type CardServiceImpl struct {
	repository repository.CardRepository
}

func NewCardService(repository repository.CardRepository) CardService {
	return &CardServiceImpl{repository: repository}
}

func (s *CardServiceImpl) Create(request card.CreateRequest) (int64, error) {
	return 0, nil
}

func (s *CardServiceImpl) Read(request card.ReadRequest) ([]model.Card, error) {
	return []model.Card{}, nil
}

func (s *CardServiceImpl) Update(request card.UpdateRequest) error {
	return nil
}

func (s *CardServiceImpl) Delete(request card.DeleteRequest) error {
	return nil
}
