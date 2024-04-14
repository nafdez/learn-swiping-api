package main

import (
	"learn-swiping-api/controller"
	"learn-swiping-api/database"
	"learn-swiping-api/repository"
	"learn-swiping-api/service"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Init
	userRepo := repository.NewUserRepository(db)
	userSrvc := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userSrvc)

	// ROUTER
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = append(config.AllowMethods, "OPTIONS")

	router.Use(cors.New(config))

	router.GET("/ping", ping)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("", userCtrl.Token)
		authGroup.POST("register", userCtrl.Register)
		authGroup.POST("login", userCtrl.Login)
		authGroup.POST("logout", userCtrl.Logout)
	}

	accountGroup := router.Group("account")
	{
		accountGroup.POST("", userCtrl.Account)
		accountGroup.PUT("", userCtrl.Update)
		accountGroup.DELETE("", userCtrl.Delete)
	}

	userGroup := router.Group("users")
	{
		userGroup.GET(":username", userCtrl.User)

	}

	router.Run(":8080")

}

func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
