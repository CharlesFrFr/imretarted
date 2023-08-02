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
	pathToProfile := "default/" + profileId + ".json"

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

	if profileId == "athena" {
		CreateLoadoutForUser(user.AccountId, "sandbox_loadout")
	}

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

func CreateLoadoutForUser(accountId string, loadoutName string) {
	file, err := os.Open("default/loadout.json")
	if err != nil {
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return
	}
	str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

	var loadout models.Loadout
	err = json.Unmarshal([]byte(str), &loadout)
	if err != nil {
		return
	}

	loadout.Attributes.LockerName = loadoutName

	marshal, err := json.Marshal(loadout)
	if err != nil {
		return
	}

	helpers.Postgres.Create(&models.UserLoadout{
		AccountId: accountId,
		Loadout:   string(marshal),
	})

	helpers.PrintGreen([]string{"created loadout", loadoutName, "for", accountId})
}

func AppendLoadoutsToProfile(profile *models.Profile, accountId string) {
	var loadouts []models.UserLoadout
	result := helpers.Postgres.Model(&models.UserLoadout{}).Where("account_id = ?", accountId).Find(&loadouts)

	if result.Error != nil {
		return
	}

	for _, loadout := range loadouts {
		var loadoutData models.Loadout
		err := json.Unmarshal([]byte(loadout.Loadout), &loadoutData)
		if err != nil {
			return
		}

		profile.Items[loadoutData.Attributes.LockerName] = loadoutData
		profile.Stats.Attributes.Loadouts = append(profile.Stats.Attributes.Loadouts, loadoutData.Attributes.LockerName)
	}

	SaveProfileToUSer(accountId, *profile)
}