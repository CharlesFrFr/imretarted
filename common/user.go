package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/zombman/server/all"
	"github.com/zombman/server/models"
)

func OnlyAllowCharacters(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
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

func VerifyGoogleRecaptcha(token string) bool {
	secret := os.Getenv("GOOGLE_RECAPTCHA_SECRET_KEY")
	if secret == "" {
		return false
	}

	if secret == "OFF" {
		return true
	}

	response, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify?secret="+ secret +"&response="+ token, url.Values{})
	if err != nil {
		return false
	}

	if response.StatusCode != 200 {
		return false
	}

	bodyclone := response.Body

	// body2 := make([]byte, 10000)
	// response.Body.Read(body2)
	// fmt.Println(string(body2))

	// read response body
	defer response.Body.Close()
	var body struct {
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(bodyclone).Decode(&body); err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(body.Success)

	return body.Success
}