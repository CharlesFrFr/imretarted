package all

import (
	"os"

	"github.com/zombman/server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Postgres *gorm.DB

func ConnectToDatabase() {
	var err error

	dsn := os.Getenv("DATABASE_URL")
	Postgres, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

func AutoMigrate() {
	Postgres.AutoMigrate(&models.User{})
	Postgres.AutoMigrate(&models.ClientToken{})
	Postgres.AutoMigrate(&models.AccessToken{})
	Postgres.AutoMigrate(&models.RefreshToken{})

	Postgres.AutoMigrate(&models.UserProfile{})
	Postgres.AutoMigrate(&models.UserLoadout{})
}