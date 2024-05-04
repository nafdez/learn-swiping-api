package progress

import (
	"errors"
	"learn-swiping-api/erro"
	progress "learn-swiping-api/internal/progress/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestWithToken interface {
	SetToken(string)
}

type ProgressController interface {
	Create(*gin.Context)
	Progress(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type ProgressControllerImpl struct {
	service ProgressService
}

func NewProgressController(service ProgressService) ProgressController {
	return &ProgressControllerImpl{service: service}
}

// Creates a progress record
// Method: POST
func (c *ProgressControllerImpl) Create(ctx *gin.Context) {
	var req progress.AccessRequest
	if err := request(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Create(req)
	if err != nil {
		if errors.Is(err, erro.ErrCardNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrProgressExists) {
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

	ctx.JSON(http.StatusCreated, gin.H{})
}

// Retrieves a progress record
// Method: GET
func (c *ProgressControllerImpl) Progress(ctx *gin.Context) {
	var req progress.AccessRequest
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrInvalidToken.Error()})
		return
	}

	cardID, err := strconv.Atoi(ctx.Param("cardID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField.Error()})
		return
	}

	req.Token = token
	req.CardID = int64(cardID)

	progress, err := c.service.Progress(req)
	if err != nil {
		if errors.Is(err, erro.ErrProgressNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, progress)
}

// Updates a progress record
// Method: PUT
func (c *ProgressControllerImpl) Update(ctx *gin.Context) {
	var req progress.UpdateRequest
	if err := request(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Update(req)
	if err != nil {
		if errors.Is(err, erro.ErrBadField) || errors.Is(err, erro.ErrInvalidToken) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrProgressNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// Deletes a progress record
// Method: DELETE
func (c *ProgressControllerImpl) Delete(ctx *gin.Context) {
	var req progress.AccessRequest
	if err := request(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Delete(req); err != nil {
		if errors.Is(err, erro.ErrInvalidToken) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, erro.ErrProgressNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

// Binds the token header and a JSON body to a struct that implements
// RequestWithToken interface
func request(ctx *gin.Context, req RequestWithToken) error {
	token := ctx.GetHeader("Token")
	if token == "" {
		return erro.ErrInvalidToken
	}
	req.SetToken(token)

	if err := ctx.ShouldBindJSON(req); err != nil {
		return erro.ErrBadField
	}
	return nil
}
