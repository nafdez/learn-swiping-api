package card

import (
	"errors"
	"learn-swiping-api/erro"
	card "learn-swiping-api/internal/card/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardController interface {
	Create(*gin.Context) // POST
	Card(*gin.Context)   // GET
	Cards(*gin.Context)  // GET
	Update(*gin.Context) // PUT
	Delete(*gin.Context) // DELETE
}

type CardControllerImpl struct {
	service CardService
}

func NewCardController(service CardService) CardController {
	return &CardControllerImpl{service: service}
}

// Creates a card
// Method: POST
func (c *CardControllerImpl) Create(ctx *gin.Context) {
	var request card.CreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	deckID, err := strconv.Atoi(ctx.Param("deckID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	request.DeckID = int64(deckID)

	if _, err := c.service.Create(request); err != nil {
		if errors.Is(err, erro.ErrBadField) {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, erro.ErrDeckNotFound) || errors.Is(err, erro.ErrCardNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

// Retrieves a card by it's id inside a deck
// Method: GET
func (c *CardControllerImpl) Card(ctx *gin.Context) {
	cardID, err := strconv.Atoi(ctx.Param("cardID"))
	deckID, derr := strconv.Atoi(ctx.Param("deckID"))
	if err != nil || derr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	card, err := c.service.Card(int64(cardID), int64(deckID))
	if err != nil {
		if errors.Is(err, erro.ErrCardNotFound) || errors.Is(err, erro.ErrWrongNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, card)
}

// Retrieves a list of cards based on it's deckID
// Method: GET
func (c *CardControllerImpl) Cards(ctx *gin.Context) {
	deckID, err := strconv.Atoi(ctx.Param("deckID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	cards, err := c.service.Cards(int64(deckID))
	if err != nil {
		if errors.Is(err, erro.ErrCardNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cards)
}

// Updates a card or it's wrong answers
// Method: PUT
func (c *CardControllerImpl) Update(ctx *gin.Context) {
	var request card.UpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	deckID, err := strconv.Atoi(ctx.Param("deckID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	request.DeckID = int64(deckID)

	if err := c.service.Update(request); err != nil {
		if errors.Is(err, erro.ErrCardNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// Deletes a card and it's wrong answers
// Method: DELETE
func (c *CardControllerImpl) Delete(ctx *gin.Context) {
	var request card.DeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	deckID, err := strconv.Atoi(ctx.Param("deckID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	request.DeckID = int64(deckID)

	if err := c.service.Delete(request.Id, request.DeckID); err != nil {
		if errors.Is(err, erro.ErrCardNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
