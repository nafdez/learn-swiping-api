package deck

import (
	"errors"
	"fmt"
	"learn-swiping-api/erro"
	deck "learn-swiping-api/internal/deck/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeckController interface {
	Create(*gin.Context)                 // POST
	Deck(*gin.Context)                   // GET
	OwnedDecks(*gin.Context)             // GET
	Subscriptions(*gin.Context)          // GET
	Update(*gin.Context)                 // PUT
	Delete(*gin.Context)                 // DELETE
	AddDeckSubscription(*gin.Context)    // POST
	RemoveDeckSubscription(*gin.Context) // DELETE
}

type DeckControllerImpl struct {
	service DeckService
}

func NewDeckController(service DeckService) DeckController {
	return &DeckControllerImpl{service: service}
}

// Creates a deck
// Method: POST
func (c *DeckControllerImpl) Create(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	var request deck.CreateRequest
	request.Token = token
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	deckID, err := c.service.Create(request)
	if err != nil {
		if errors.Is(err, erro.ErrAccountNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrDeckExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrInvalidToken) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"deck_id": deckID})
}

// Retrieves a deck by it's id
// Method: GET
func (c *DeckControllerImpl) Deck(ctx *gin.Context) {
	token := ctx.GetHeader("Token") // Optional
	idParam := ctx.Param("deckID")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	var request deck.ReadOneRequest
	_ = ctx.ShouldBindJSON(&request) // Not checking on errors since this field is optional
	request.DeckID = int64(id)

	deck, err := c.service.Deck(request, token)
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
	token := ctx.GetHeader("Token") // Optional
	var request deck.ReadOwnedRequest
	ctx.ShouldBindJSON(&request) // No point checking for errors since it's optional
	if request.Username == "" {
		request.Username = ctx.Param("username")
	}

	decks, err := c.service.OwnedDecks(request, token)
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
// Method: POST
func (c *DeckControllerImpl) Subscriptions(ctx *gin.Context) {
	token := ctx.GetHeader("Token") // Optional
	var request deck.ReadRequest
	ctx.ShouldBindJSON(&request)
	if request.Username == "" {
		request.Username = ctx.Param("username")
	}

	decks, err := c.service.Suscriptions(request, token)
	if err != nil {
		if errors.Is(err, erro.ErrBadField) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrDeckNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": erro.ErrNotSuscribed.Error()})
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

	// File fetching need to be the first or throwns an error
	// that there is no file in the body
	file, err := ctx.FormFile("picture")
	if err == nil {
		request.Img = file
	} else {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusTeapot, gin.H{"error": err.Error()})
		return
	}

	// Header Token needed to update
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	// Param deckID needed to update
	deckIDSTR := ctx.Param("deckID")
	deckID, err := strconv.Atoi(deckIDSTR)
	if deckID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}
	request.DeckID = int64(deckID)

	_ = ctx.ShouldBindJSON(&request)

	if err := c.service.Update(request, token); err != nil {
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
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	var request deck.DeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	if err := c.service.Delete(request, token); err != nil {
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

	ctx.JSON(http.StatusOK, gin.H{})
}

// Subscribes a user to a deck
// Method: POST
func (c *DeckControllerImpl) AddDeckSubscription(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	var request deck.DeckSuscriptionRequest
	request.Token = token
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	if err := c.service.AddDeckSubscription(request); err != nil {
		if errors.Is(err, erro.ErrAlreadySuscribed) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrInvalidToken) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	var request deck.DeckSuscriptionRequest
	request.Token = token
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	if err := c.service.RemoveDeckSubscription(request); err != nil {
		if errors.Is(err, erro.ErrNotSuscribed) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrInvalidToken) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
