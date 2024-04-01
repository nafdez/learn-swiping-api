package controller

import (
	"learn-swiping-api/model/dto/user"
	"learn-swiping-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(*gin.Context) // POST
	Login(*gin.Context)    // POST
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
	}

	user, err := c.service.Register(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserControllerImpl) Login(ctx *gin.Context) {
}

func (c *UserControllerImpl) Logout(ctx *gin.Context) {
}

func (c *UserControllerImpl) Account(ctx *gin.Context) {
}

func (c *UserControllerImpl) User(ctx *gin.Context) {
}

func (c *UserControllerImpl) Update(ctx *gin.Context) {
}

func (c *UserControllerImpl) Delete(ctx *gin.Context) {
}
