package controller

import (
	"learn-swiping-api/service"

	"github.com/gin-gonic/gin"
)

type DeckController interface {
	Create(*gin.Context)     // POST
	Deck(*gin.Context)       // GET
	OwnedDecks(*gin.Context) // GET
	Decks(*gin.Context)      // GET (from other users)
	Update(*gin.Context)     // PUT
	Delete(*gin.Context)     // DELETE
}

type DeckControllerImpl struct {
	service service.DeckService
}

func NewDeckController(service service.DeckService) DeckController {
	return &DeckControllerImpl{service: service}
}

func (c *DeckControllerImpl) Create(ctx *gin.Context) {
}

func (c *DeckControllerImpl) Deck(ctx *gin.Context) {
}

func (c *DeckControllerImpl) OwnedDecks(ctx *gin.Context) {
}

func (c *DeckControllerImpl) Decks(ctx *gin.Context) {
}

func (c *DeckControllerImpl) Update(ctx *gin.Context) {
}

func (c *DeckControllerImpl) Delete(ctx *gin.Context) {
}
