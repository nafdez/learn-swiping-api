package picture

import (
	"learn-swiping-api/erro"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PictureController struct {
}

func NewPictureController() PictureController {
	return PictureController{}
}

func (c *PictureController) Picture(ctx *gin.Context) {
	picID := ctx.Param("picID")
	if picID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": erro.ErrBadField})
		return
	}

	ctx.File("./data/pictures/" + picID)
}
