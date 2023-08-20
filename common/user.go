package common

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/zombman/server/all"
	"github.com/zombman/server/models"
)

func OnlyAllowCharacters(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		all.PrintGreen([]any{"serre", err})
		return ""
	}
	return reg.ReplaceAllString(s, "")
}

func CreateUser(username string, password string, level int) (models.User, error) {
	user := models.User{
		Username:  OnlyAllowCharacters(username),
		Password:  all.HashString(password),
		AccountId: uuid.New().String(),
		AccessLevel: level,
	}
	result := all.Postgres.Create(&user)
	
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

	result := all.Postgres.Where("account_id = ? AND banned = false", accountId).First(&user)

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

	result := all.Postgres.Where("username = ? AND banned = false", username).First(&user)

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

	result := all.Postgres.Where("username = ? AND password = ? AND banned = false", username, all.HashString(password)).First(&user)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	if user.ID == 0 {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func GetUserByUsernameAndHashPassword(username string, password string) (models.User, error) {
	var user models.User

	result := all.Postgres.Where("username = ? AND password = ? AND banned = false", username, password).First(&user)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	if user.ID == 0 {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}