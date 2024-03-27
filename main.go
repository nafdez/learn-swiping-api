package main

import (
	"learn-swiping-api/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	_, err = database.Connect()
	if err != nil {
		log.Fatalln(err)
	}

}
