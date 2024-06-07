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

	// CHAOS ZONE
	// Proceed with caution
	authGroup := router.Group("/auth")
	{
		authGroup.GET("", init.UserCtrl.Token)
		authGroup.POST("register", init.UserCtrl.Register)
		authGroup.POST("login", init.UserCtrl.Login)
		authGroup.GET("token", init.UserCtrl.Token) // TODO: Migrate to Account fn
		authGroup.DELETE("logout", init.UserCtrl.Logout)
	}

	accountGroup := router.Group("account")
	{
		accountGroup.GET("", init.UserCtrl.Account)
		accountGroup.PUT("", init.UserCtrl.Update)
		accountGroup.DELETE("", init.UserCtrl.Delete)
	}

	userGroup := router.Group("users")
	{
		userGroup.GET(":username", init.UserCtrl.AccountPublic)
		userGroup.GET(":username/decks", init.DeckCtrl.OwnedDecks)
		userGroup.GET(":username/subscribed", init.DeckCtrl.Subscriptions)
	}

	deckGroup := router.Group("decks")
	{
		deckGroup.POST("", init.DeckCtrl.Create)
		deckGroup.PUT(":deckID", init.DeckCtrl.Update)
		deckGroup.DELETE(":deckID", init.DeckCtrl.Delete)
		deckGroup.GET(":deckID", init.DeckCtrl.DeckDetails)

		deckGroup.POST("subs/:deckID", init.DeckCtrl.AddDeckSubscription)
		deckGroup.DELETE("subs/:deckID", init.DeckCtrl.RemoveDeckSubscription)
		deckGroup.GET("subs/:username/:deckID", init.DeckCtrl.DeckDetails)

		deckGroup.POST(":deckID/rating/:rating", init.DeckCtrl.SaveRating)
		deckGroup.GET(":deckID/rating", init.DeckCtrl.Rating)
		deckGroup.DELETE(":deckID/rating", init.DeckCtrl.DeleteRating)

		deckGroup.POST(":deckID", init.CardCtrl.Create)
		deckGroup.GET(":deckID/:cardID", init.CardCtrl.Card)
		deckGroup.GET(":deckID/cards", init.CardCtrl.Cards)
		deckGroup.PUT(":deckID/:cardID", init.CardCtrl.Update)
		deckGroup.DELETE(":deckID/:cardID", init.CardCtrl.Delete)
	}

	shopGroup := router.Group("shop")
	{
		shopGroup.GET(":deckID", init.DeckCtrl.DeckDetailsShop)
	}

	progressGroup := router.Group("progress")
	{
		progressGroup.POST("", init.ProgressCtrl.Create)
		progressGroup.GET(":cardID", init.ProgressCtrl.Progress)
		progressGroup.PUT("", init.ProgressCtrl.Update)
		progressGroup.DELETE("", init.ProgressCtrl.Delete)
	}

	pictureGroup := router.Group("pics")
	{
		pictureGroup.GET(":picID", init.PictureCtrl.Picture)
	}

	return router

}

func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
