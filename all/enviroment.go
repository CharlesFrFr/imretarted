package all

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	var mode string
	if IsProduction {
		mode = gin.ReleaseMode
	} else {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)
}