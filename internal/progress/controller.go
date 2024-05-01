package progress

import "github.com/gin-gonic/gin"

type ProgressController interface {
	Create(*gin.Context)
	Progress(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type ProgressControllerImpl struct {
}
