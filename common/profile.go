package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/zombman/server/helpers"
	"github.com/zombman/server/models"
)

func AddProfileToUser(user models.User, profileId string) {
	pathToProfile := "profiles/" + profileId + ".json"

	file, err := os.Open(pathToProfile)
	if err != nil {
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return
	}
	str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

	helpers.Postgres.Create(&models.UserProfile{
		AccountId: user.AccountId,
		ProfileId: profileId,
		Profile:   str,
	})

	helpers.PrintGreen([]string{profileId, "profile added to", user.Username})
}

func ReadProfileFromUser(accountId string, profileId string) (models.Profile, error) {
	var userProfile models.UserProfile
	result := helpers.Postgres.Model(&models.UserProfile{}).Where("account_id = ? AND profile_id = ?", accountId, profileId).First(&userProfile)

	if result.Error != nil {
		return models.Profile{}, result.Error
	}

	if userProfile.ID == 0 {
		return models.Profile{}, errors.New("profile not found")
	}

	var profileData models.Profile
	err := json.Unmarshal([]byte(userProfile.Profile), &profileData)
	if err != nil {
		return models.Profile{}, err
	}

	return profileData, nil
}

func SaveProfileToUSer(accountId string, profile models.Profile) error {
	profileData, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	result := helpers.Postgres.Model(&models.UserProfile{}).Where("account_id = ? AND profile_id = ?", accountId, profile.ProfileId).Update("profile", string(profileData))
	if result.Error != nil {
		return result.Error
	}

	return nil
}