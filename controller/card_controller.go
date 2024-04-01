package controller

import (
	"learn-swiping-api/service"

	"github.com/gin-gonic/gin"
)

type CardController interface {
	Create(*gin.Context)       // POST
	Card(*gin.Context)         // GET
	Update(*gin.Context)       // PUT
	Delete(*gin.Context)       // DELETE
	CreateAnswer(*gin.Context) // POST
	Answers(*gin.Context)      // GET
	UpdateAnswer(*gin.Context) // PUT
	DeleteAnswer(*gin.Context) // DELETE
}

type CardControllerImpl struct {
	service service.CardService
}

func NewCardController(service service.CardService) CardController {
	return &CardControllerImpl{service: service}
}

func (c *CardControllerImpl) Create(ctx *gin.Context) {
}

func (c *CardControllerImpl) Card(ctx *gin.Context) {
}

func (c *CardControllerImpl) Update(ctx *gin.Context) {
}

func (c *CardControllerImpl) Delete(ctx *gin.Context) {
}

func (c *CardControllerImpl) CreateAnswer(ctx *gin.Context) {
}

func (c *CardControllerImpl) Answers(ctx *gin.Context) {
}

func (c *CardControllerImpl) UpdateAnswer(ctx *gin.Context) {
}

func (c *CardControllerImpl) DeleteAnswer(ctx *gin.Context) {
}
