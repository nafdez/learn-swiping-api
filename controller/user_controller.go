package controller

import (
	"errors"
	"learn-swiping-api/erro"
	"learn-swiping-api/model/dto/user"
	"learn-swiping-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(*gin.Context) // POST
	Login(*gin.Context)    // POST
	Token(*gin.Context)    // POST
	Logout(*gin.Context)   // POST
	Account(*gin.Context)  // POST
	User(*gin.Context)     // GET
	Update(*gin.Context)   // PUT
	Delete(*gin.Context)   // DELETE
}

type UserControllerImpl struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &UserControllerImpl{service: service}
}

func (c *UserControllerImpl) Register(ctx *gin.Context) {
	var request user.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.Register(request)
	if err != nil {
		if errors.Is(err, erro.ErrUserExists) {
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

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserControllerImpl) Login(ctx *gin.Context) {
	var request user.LoginRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.Login(request)
	if err != nil {
		if errors.Is(err, erro.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if errors.Is(err, erro.ErrBadField) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserControllerImpl) Token(ctx *gin.Context) {
	var request user.TokenRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.Token(request)
	if err != nil {
		if errors.Is(err, erro.ErrInvalidToken) || errors.Is(err, erro.ErrUserNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserControllerImpl) Logout(ctx *gin.Context) {
	var request user.TokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Logout(request)
	if err != nil {
		if errors.Is(err, erro.ErrInvalidToken) || errors.Is(err, erro.ErrUserNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *UserControllerImpl) Account(ctx *gin.Context) {
	var request user.TokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.Account(request)
	if err != nil {
		if errors.Is(err, erro.ErrInvalidToken) || errors.Is(err, erro.ErrUserNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserControllerImpl) User(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := c.service.User(user.PublicRequest{Username: username})
	if err != nil {
		if errors.Is(err, erro.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserControllerImpl) Update(ctx *gin.Context) {
	var request user.UpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Update(request)
	if err != nil {
		if errors.Is(err, erro.ErrInvalidToken) || errors.Is(err, erro.ErrUserNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": erro.ErrInvalidToken.Error()})
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

func (c *UserControllerImpl) Delete(ctx *gin.Context) {
	var request user.TokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Delete(request)
	if err != nil {
		if errors.Is(err, erro.ErrInvalidToken) || errors.Is(err, erro.ErrUserNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": erro.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
