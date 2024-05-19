package main

import (
	"learn-swiping-api/config"
	"learn-swiping-api/config/database"
	"learn-swiping-api/router"
	"log"
	"os"

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

	init := config.NewInitialization(db)
	router := router.NewRouter(init)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "9999"
	}
	router.Run(":" + port)
}
