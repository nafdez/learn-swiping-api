package router

import (
	"learn-swiping-api/config"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(init *config.Initialization) *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = append(config.AllowMethods, "OPTIONS")

	router.Use(cors.New(config))

	router.GET("/ping", ping)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("", init.UserCtrl.Token)
		authGroup.POST("register", init.UserCtrl.Register)
		authGroup.POST("login", init.UserCtrl.Login)
		authGroup.POST("logout", init.UserCtrl.Logout)
	}

	accountGroup := router.Group("account")
	{
		accountGroup.POST("", init.UserCtrl.Account)
		accountGroup.PUT("", init.UserCtrl.Update)
		accountGroup.DELETE("", init.UserCtrl.Delete)
	}

	userGroup := router.Group("users")
	{
		userGroup.GET(":username", init.UserCtrl.User)

	}

	return router

}

func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
