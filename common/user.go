package common

import (
	"errors"

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