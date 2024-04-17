package service

import "learn-swiping-api/model/dto/deck"

type DeckService interface {
	Create(deck.CreateRequest)
	Read(deck.ReadRequest)
	Update(deck.UpdateRequest)
	Delete(deck.DeleteRequest)
}
