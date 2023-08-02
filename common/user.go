package common

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/zombman/server/helpers"
	"github.com/zombman/server/models"
)

func CreateUser(username string, password string) (models.User, error) {
	user := models.User{
		Username:  username,
		Password:  helpers.HashString(password),
		AccountId: uuid.New().String(),
	}

	result := helpers.Postgres.Create(&user)
	
	if result.Error != nil {
		return models.User{}, result.Error
	}

	AddProfileToUser(user, "athena")
	AddProfileToUser(user, "campaign")
	AddProfileToUser(user, "collection_book_people0")
	AddProfileToUser(user, "collection_book_schematics0")
	AddProfileToUser(user, "collections")
	AddProfileToUser(user, "common_core")
	AddProfileToUser(user, "common_public")
	AddProfileToUser(user, "creative")
	AddProfileToUser(user, "metadata")
	AddProfileToUser(user, "outpost0")
	AddProfileToUser(user, "profile0")
	AddProfileToUser(user, "theater0")

	return user, nil
}

func GetUserByAccountId(accountId string) (models.User, error) {
	var user models.User

	result := helpers.Postgres.Where("account_id = ? AND banned = false", accountId).First(&user)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	if user.ID == 0 {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User

	result := helpers.Postgres.Where("username = ? AND banned = false", username).First(&user)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	if user.ID == 0 {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func GetUserByUsernameAndPlainPassword(username string, password string) (models.User, error) {
	var user models.User

	result := helpers.Postgres.Where("username = ? AND password = ? AND banned = false", username, helpers.HashString(password)).First(&user)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	if user.ID == 0 {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func AddProfileToUser(user models.User, profileId string) {
	pathToProfile := "profiles/" + profileId + ".json"

	file, err := os.Open(pathToProfile)
	if err != nil {
		return
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	helpers.Postgres.Create(&models.UserProfile{
		AccountId: user.AccountId,
		ProfileId: profileId,
		Profile:   string(fileData),
	})

	helpers.PrintGreen([]string{profileId, "profile added to", user.Username})
}