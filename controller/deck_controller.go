package controller

import (
	"errors"
	"learn-swiping-api/erro"
	"learn-swiping-api/model/dto/deck"
	"learn-swiping-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeckController interface {
	Create(*gin.Context)                 // POST
	Deck(*gin.Context)                   // GET
	OwnedDecks(*gin.Context)             // POST (to accept token)
	Suscriptions(*gin.Context)           // POST (from other users)
	Update(*gin.Context)                 // PUT
	Delete(*gin.Context)                 // DELETE
	AddDeckSubscription(*gin.Context)    // POST
	RemoveDeckSubscription(*gin.Context) // DELETE
}

type DeckControllerImpl struct {
	service service.DeckService
}

func NewDeckController(service service.DeckService) DeckController {
	return &DeckControllerImpl{service: service}
}

// Creates a deck
// Method: POST
func (c *DeckControllerImpl) Create(ctx *gin.Context) {
	var request deck.CreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField})
		return
	}

	_, err := c.service.Create(request)
	if err != nil {
		if errors.Is(err, erro.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

// Retrieves a deck by it's id
// Method: GET
func (c *DeckControllerImpl) Deck(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField})
		return
	}

	var request deck.ReadOneRequest
	_ = ctx.ShouldBindJSON(&request) // Not checking on errors since this field is optional
	request.Id = int64(id)

	deck, err := c.service.Deck(request)
	if err != nil {
		if errors.Is(err, erro.ErrBadField) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrDeckNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, deck)
}

// Retrieves a list of decks a user has created
// Method: POST
func (c *DeckControllerImpl) OwnedDecks(ctx *gin.Context) {
	var request deck.ReadOwnedRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": erro.ErrBadField})
		return
	}

	decks, err := c.service.OwnedDecks(request)
	if err != nil {
		if errors.Is(err, erro.ErrBadField) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrDeckNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, decks)
}

// Retrieves a list of decks that a user has been subscribed for
// Method: PPOST
func (c *DeckControllerImpl) Suscriptions(ctx *gin.Context) {
	var request deck.ReadRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": erro.ErrBadField})
		return
	}

	decks, err := c.service.Suscriptions(request)
	if err != nil {
		if errors.Is(err, erro.ErrBadField) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrDeckNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": erro.ErrNotSuscribed})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, decks)
}

// Updates a deck
// Method: PUT
func (c *DeckControllerImpl) Update(ctx *gin.Context) {
	var request deck.UpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField})
		return
	}

	if err := c.service.Update(request); err != nil {
		if errors.Is(err, erro.ErrBadField) || errors.Is(err, erro.ErrInvalidToken) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrDeckNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// Deletes a deck
// Method: DELETE
func (c *DeckControllerImpl) Delete(ctx *gin.Context) {
	var request deck.DeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField})
		return
	}

	if err := c.service.Delete(request); err != nil {
		if errors.Is(err, erro.ErrInvalidToken) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, erro.ErrDeckNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// Subscribes a user to a deck
// Method: POST
func (c *DeckControllerImpl) AddDeckSubscription(ctx *gin.Context) {
	var request deck.DeckSuscriptionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField})
		return
	}

	if err := c.service.AddDeckSubscription(request); err != nil {
		if errors.Is(err, erro.ErrAlreadySuscribed) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// Unsubscribes a deck from a user
// Method: DELETE
func (c *DeckControllerImpl) RemoveDeckSubscription(ctx *gin.Context) {
	var request deck.DeckSuscriptionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField})
		return
	}

	if err := c.service.RemoveDeckSubscription(request); err != nil {
		if errors.Is(err, erro.ErrNotSuscribed) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
