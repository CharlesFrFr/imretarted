package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var IsProduction bool

func LoadEnviroment() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("PRODUCTION") == "true" {
		IsProduction = true
	}
}