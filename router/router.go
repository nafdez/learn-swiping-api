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
		authGroup.POST("", init.UserCtrl.Token)
		authGroup.POST("register", init.UserCtrl.Register)
		authGroup.POST("login", init.UserCtrl.Login)
		authGroup.POST("token", init.UserCtrl.Token) // TODO: Migrate to Account fn
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
		userGroup.GET(":username/decks", init.DeckCtrl.OwnedDecks)
		userGroup.GET(":username/subscribed", init.DeckCtrl.Subscriptions)
	}

	deckGroup := router.Group("decks")
	{
		deckGroup.POST("", init.DeckCtrl.Create)
		deckGroup.PUT("", init.DeckCtrl.Update)
		deckGroup.DELETE("", init.DeckCtrl.Delete)
		deckGroup.GET(":deckID", init.DeckCtrl.Deck)

		deckGroup.POST("subscription", init.DeckCtrl.AddDeckSubscription)
		deckGroup.DELETE("subscription", init.DeckCtrl.RemoveDeckSubscription)

		deckGroup.POST(":deckID", init.CardCtrl.Create)
		deckGroup.GET(":deckID/:cardID", init.CardCtrl.Card)
		deckGroup.GET(":deckID/cards", init.CardCtrl.Cards)
		deckGroup.PUT(":deckID", init.CardCtrl.Update)
		deckGroup.DELETE(":deckID", init.CardCtrl.Delete)
	}

	return router

}

func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
