package progress

import (
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

func (c *ProgressControllerImpl) Create(ctx *gin.Context) {
	var req progress.AccessRequest
	if err := request(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call srvc
}

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
	req.CardID = int64(cardID)

	// Call srvc
}

func (c *ProgressControllerImpl) Update(ctx *gin.Context) {
	var req progress.UpdateRequest
	if err := request(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call srvc
}

func (c *ProgressControllerImpl) Delete(ctx *gin.Context) {
	var req progress.AccessRequest
	if err := request(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call srvc
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
