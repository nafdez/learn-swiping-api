package account

import (
	"errors"
	"learn-swiping-api/erro"
	account "learn-swiping-api/internal/account/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController interface {
	Register(*gin.Context)      // POST
	Login(*gin.Context)         // POST
	Token(*gin.Context)         // POST
	Logout(*gin.Context)        // POST
	Account(*gin.Context)       // GET
	AccountPublic(*gin.Context) // GET
	Update(*gin.Context)        // PUT
	Delete(*gin.Context)        // DELETE
}

type AccountControllerImpl struct {
	service AccountService
}

func NewAccountController(service AccountService) AccountController {
	return &AccountControllerImpl{service: service}
}

// Creates a new account
// Method: POST
func (c *AccountControllerImpl) Register(ctx *gin.Context) {
	var request account.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	account, err := c.service.Register(request)
	if err != nil {
		if errors.Is(err, erro.ErrAccountExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrBadField) || errors.Is(err, erro.ErrInvalidEmail) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

// Retrieves an account if username and password are correct
// Method: POST
func (c *AccountControllerImpl) Login(ctx *gin.Context) {
	var request account.LoginRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	account, err := c.service.Login(request)
	if err != nil {
		if errors.Is(err, erro.ErrAccountNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// Retrieves an account if token is correct
// Method: POST
func (c *AccountControllerImpl) Token(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	account, err := c.service.Token(token)
	if err != nil {
		if errors.Is(err, erro.ErrInvalidToken) || errors.Is(err, erro.ErrAccountNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// Invalidates a token to restrict access to the account
// Method: POST
func (c *AccountControllerImpl) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	err := c.service.Logout(token)
	if err != nil {
		if errors.Is(err, erro.ErrAccountNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// Retrieves an account if provided token is correct
// Method: GET
func (c *AccountControllerImpl) Account(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	account, err := c.service.Account(token)
	if err != nil {
		if errors.Is(err, erro.ErrAccountNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// Retrieves a public profile of an account by it's username
// Method: GET
func (c *AccountControllerImpl) AccountPublic(ctx *gin.Context) {
	username := ctx.Param("username")

	account, err := c.service.account(account.PublicRequest{Username: username})
	if err != nil {
		if errors.Is(err, erro.ErrBadField) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if errors.Is(err, erro.ErrAccountNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// Updates an account
// Method: PUT
func (c *AccountControllerImpl) Update(ctx *gin.Context) {
	var request account.UpdateRequest

	// File fetching need to be the first or throwns an error
	// that there is no file in the body
	file, err := ctx.FormFile("picture")
	if err == nil {
		request.Img = file
	}

	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
	}

	_ = ctx.ShouldBindJSON(&request)

	request.Token = token

	err = c.service.Update(request)
	if err != nil {
		if errors.Is(err, erro.ErrInvalidToken) || errors.Is(err, erro.ErrAccountNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		if errors.Is(err, erro.ErrBadField) || errors.Is(err, erro.ErrInvalidEmail) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// Deletes an account
// Method: DELETE
func (c *AccountControllerImpl) Delete(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	err := c.service.Delete(token)
	if err != nil {
		if errors.Is(err, erro.ErrInvalidToken) || errors.Is(err, erro.ErrAccountNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
