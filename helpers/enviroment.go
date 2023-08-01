package helpers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnviroment() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}