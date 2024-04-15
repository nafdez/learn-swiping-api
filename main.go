package main

import (
	"learn-swiping-api/config"
	"learn-swiping-api/config/database"
	"learn-swiping-api/router"
	"log"

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

	router.Run(":8080")
}
