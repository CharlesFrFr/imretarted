package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/models"
)

func AddProfileToUser(user models.User, profileId string) {
	pathToProfile := "data/" + profileId + ".json"

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

	unmarshaledProfile := models.Profile{}
	err = json.Unmarshal([]byte(str), &unmarshaledProfile)
	if err != nil {
		return
	}

	if profileId == "common_core" {
		vbucksAmount, err := strconv.Atoi(os.Getenv("USER_STARTING_VBUCKS"))
		if err == nil { 
			SetUserVBucks(user.AccountId, &unmarshaledProfile, vbucksAmount)
		}
	}

	profileData, err := json.Marshal(unmarshaledProfile)
	if err != nil {
		return
	}

	all.Postgres.Create(&models.UserProfile{
		AccountId: user.AccountId,
		ProfileId: profileId,
		Profile:   string(profileData),
	})

	if profileId == "athena" {
		all.PrintBlue([]any{"creating loadouts on athena profile for", user.Username})
		CreateLoadoutForUser(user.AccountId, "sandbox_loadout")
		CreateLoadoutForUser(user.AccountId, "zombie_loadout")

		athenaProfile, err := ReadProfileFromUser(user.AccountId, "athena")
		if err != nil {
			return
		}

		newAthenaProfile, err := ConvertProfileToAthena(athenaProfile)
		if err != nil {
			return
		}

		newAthenaProfile.Stats.Attributes.Level = 1
		newAthenaProfile.Stats.Attributes.Xp = 0
		newAthenaProfile.Stats.Attributes.BookLevel = 1
		newAthenaProfile.Stats.Attributes.BookXp = 0
		newAthenaProfile.Stats.Attributes.LifetimeWins = 10

		defaultProfile, err := ConvertAthenaToDefault(newAthenaProfile)
		if err != nil {
			return
		}

		AddItemToProfile(&defaultProfile, "AthenaCharacter:CID_001_Athena_Commando_F_Default", user.AccountId)
		SaveProfileToUser(user.AccountId, defaultProfile)
	}

	all.PrintGreen([]any{profileId, "profile added to", user.Username})
}

func ReadProfileFromUser(accountId string, profileId string) (models.Profile, error) {
	var userProfile models.UserProfile
	result := all.Postgres.Model(&models.UserProfile{}).Where("account_id = ? AND profile_id = ?", accountId, profileId).First(&userProfile)

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

func ConvertProfileToCommonCore(profile models.Profile) (models.CommonCoreProfile, error) {
	marshalledProfile, err := json.Marshal(profile)
	if err != nil {
		return models.CommonCoreProfile{}, err
	}

	var newProfile models.CommonCoreProfile
	err = json.Unmarshal(marshalledProfile, &newProfile)
	if err != nil {
		return models.CommonCoreProfile{}, err
	}

	return newProfile, nil
}

func ConvertProfileToAthena(profile models.Profile) (models.AthenaProfile, error) {
	marshalledProfile, err := json.Marshal(profile)
	if err != nil {
		return models.AthenaProfile{}, err
	}

	var newProfile models.AthenaProfile
	err = json.Unmarshal(marshalledProfile, &newProfile)
	if err != nil {
		return models.AthenaProfile{}, err
	}

	return newProfile, nil
}

func ConvertAthenaToDefault(profile models.AthenaProfile) (models.Profile, error) {
	marshalledProfile, err := json.Marshal(profile)
	if err != nil {
		return models.Profile{}, err
	}

	var newProfile models.Profile
	err = json.Unmarshal(marshalledProfile, &newProfile)
	if err != nil {
		return models.Profile{}, err
	}

	return newProfile, nil
}

func ConvertCommonCoreToDefault(profile models.CommonCoreProfile) (models.Profile, error) {
	marshalledProfile, err := json.Marshal(profile)
	if err != nil {
		return models.Profile{}, err
	}

	var newProfile models.Profile
	err = json.Unmarshal(marshalledProfile, &newProfile)
	if err != nil {
		return models.Profile{}, err
	}

	return newProfile, nil
}

func SaveProfileToUser(accountId string, profile models.Profile) error {
	profileData, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	result := all.Postgres.Model(&models.UserProfile{}).Where("account_id = ? AND profile_id = ?", accountId, profile.ProfileId).Update("profile", string(profileData))
	if result.Error != nil {
		return result.Error
	}

	all.PrintRed([]any{"saved profile", profile.ProfileId, "for", accountId})

	return nil
}

func CreateLoadoutForUser(accountId string, loadoutName string) {
	file, err := os.Open("data/loadout.json")
	if err != nil {
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		all.PrintRed([]any{"error reading loadout.json"})
		return
	}
	str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

	var loadout models.Loadout
	err = json.Unmarshal([]byte(str), &loadout)
	if err != nil {
		all.PrintRed([]any{err.Error()})
		return
	}

	loadout.Attributes.LockerName = loadoutName

	marshal, err := json.Marshal(loadout)
	if err != nil {
		return
	}

	all.Postgres.Create(&models.UserLoadout{
		AccountId: accountId,
		Loadout:   string(marshal),
		LoadoutName: loadoutName,
	})

	all.PrintGreen([]any{"created loadout", loadoutName, "for", accountId})
}

func AppendLoadoutToProfileNoSave(profile *models.Profile, loadout *models.Loadout, accountId string) {
	var userLoadout models.UserLoadout
	result := all.Postgres.Model(&models.UserLoadout{}).Where("account_id = ? AND loadout_name = ?", accountId, loadout.Attributes.LockerName).First(&userLoadout)

	if result.Error != nil {
		return
	}

	profile.Items[loadout.Attributes.LockerName] = *loadout

	var marshaledLoadout []byte
	marshaledLoadout, err := json.Marshal(loadout)
	if err != nil {
		return
	}

	result = all.Postgres.Model(&models.UserLoadout{}).Where("account_id = ? AND loadout_name = ?", accountId, loadout.Attributes.LockerName).Update("loadout", string(marshaledLoadout))
	if result.Error != nil {
		return
	}
}

func AppendLoadoutsToProfileNoSave(profile *models.Profile, accountId string) {
	var loadouts []models.UserLoadout
	result := all.Postgres.Model(&models.UserLoadout{}).Where("account_id = ?", accountId).Find(&loadouts)

	if result.Error != nil {
		return
	}

	loadoutIds := []string{}

	for _, loadout := range loadouts {
		var loadoutData models.Loadout
		err := json.Unmarshal([]byte(loadout.Loadout), &loadoutData)
		if err != nil {
			return
		}

		loadoutIds = append(loadoutIds, loadoutData.Attributes.LockerName)

		profile.Items[loadoutData.Attributes.LockerName] = loadoutData
		profile.Stats.Attributes["loadouts"] = loadoutIds
		profile.Stats.Attributes["active_loadout_index"] = len(loadoutIds) - 1
		profile.Stats.Attributes["last_applied_loadout"] = loadoutData.Attributes.LockerName
	}
}

func AppendLoadoutToProfile(profile *models.Profile, loadout *models.Loadout, accountId string) {
	AppendLoadoutToProfileNoSave(profile, loadout, accountId)
	SaveProfileToUser(accountId, *profile)
}

func AppendLoadoutsToProfile(profile *models.Profile, accountId string) {
	AppendLoadoutsToProfileNoSave(profile, accountId)
	SaveProfileToUser(accountId, *profile)
}

func GetLoadout(loadoutId string, accountId string) (models.Loadout, error) {
	var loadouts []models.UserLoadout
	result := all.Postgres.Model(&models.UserLoadout{}).Where("account_id = ?", accountId).Find(&loadouts)
	if result.Error != nil {
		return models.Loadout{}, result.Error
	}

	for _, loadout := range loadouts {
		var loadoutData models.Loadout
		err := json.Unmarshal([]byte(loadout.Loadout), &loadoutData)
		if err != nil {
			return models.Loadout{}, err
		}

		if loadoutData.Attributes.LockerName == loadoutId {
			return loadoutData, nil
		}
	}

	return models.Loadout{}, errors.New("loadout not found")
}

func AddItemToProfile(profile *models.Profile, itemId string, accountId string) {
	profile.Items[itemId] = models.Item{
		TemplateId: itemId,
		Attributes: models.ItemAttributes{
			MaxLevelBonus: 0,
			Level: 1,
			ItemSeen: true,
			Variants: []models.ItemVariant{},
			Favorite: false,
			Xp: 0,
		},
		Quantity: 1,
	}
	AppendLoadoutsToProfileNoSave(profile, accountId)
}

func AddItemsToProfile(profile *models.Profile, itemIds []string, accountId string) {
	for _, itemId := range itemIds {
		profile.Items[itemId] = models.Item{
			TemplateId: itemId,
			Attributes: models.ItemAttributes{
				MaxLevelBonus: 0,
				Level: 1,
				ItemSeen: true,
				Variants: []models.ItemVariant{},
				Favorite: false,
				Xp: 0,
			},
			Quantity: 1,
		}
	}
	AppendLoadoutsToProfileNoSave(profile, accountId)
}

func RemoveItemFromProfile(profile *models.Profile, itemId string, accountId string) {
	delete(profile.Items, itemId)
	AppendLoadoutsToProfileNoSave(profile, accountId)
}

func RemoveItemsFromProfile(profile *models.Profile, itemIds []string, accountId string) {
	for _, itemId := range itemIds {
		delete(profile.Items, itemId)
	}
	AppendLoadoutsToProfileNoSave(profile, accountId)
}

func AddEverythingToProfile(profile *models.Profile, accountId string) {
	pathToAllItems := "data/shop/all.json"

	file, err := os.Open(pathToAllItems)
	if err != nil {
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return
	}
	str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

	var itemsData []models.BeforeStoreItem
	err = json.Unmarshal([]byte(str), &itemsData)
	if err != nil {
		return
	}

	var itemIds []string
	for _, item := range itemsData {
		itemIds = append(itemIds, item.BackendType + ":" + item.ID)
	}

	commonCore, err := ReadProfileFromUser(accountId, "common_core")
	if err != nil {
		return
	}

	gift := models.CommonCoreItem{
		TemplateId: "GiftBox:gb_default",
		Attributes: gin.H{
			"fromAccountId": "Server",
			"lootList": []gin.H{},
			"params": gin.H{
				"userMessage": "Enjoy this gift from the server!",
			},
			"level": 1,
			"giftedOn": time.Now().Format("2006-01-02T15:04:05.999Z"),
		},
		Quantity: 1,
	}

	for _, item := range itemIds {
		gift.Attributes["lootList"] = append(gift.Attributes["lootList"].([]gin.H), gin.H{
			"itemType": item,
			"itemGuid": item,
			"itemProfile": "athena",
			"quantity": 1,
		})
	}

	commonCore.Items["GiftBox:gb_default"] = gift
	AddItemsToProfile(profile, itemIds, accountId)
	SaveProfileToUser(accountId, commonCore)

	all.PrintGreen([]any{"added all items to profile", accountId})
}

func RemoveEverythingFromProfile(profile *models.Profile, accountId string) {
	for itemId := range profile.Items {
		delete(profile.Items, itemId)
	}
	AppendLoadoutsToProfileNoSave(profile, accountId)

	commonCore, err := ReadProfileFromUser(accountId, "common_core")
	if err != nil {
		return
	}

	gift := models.CommonCoreItem{
		TemplateId: "GiftBox:gb_default",
		Attributes: gin.H{
			"fromAccountId": "Server",
			"lootList": []gin.H{
				{
					"itemType": "AthenaCharacter:CID_001_Athena_Commando_F_Default",
					"itemGuid": "AthenaCharacter:CID_001_Athena_Commando_F_Default",
					"itemProfile": "athena",
					"quantity": 1,
				},
				{
					"itemType": "AthenaPickaxe:DefaultPickaxe",
					"itemGuid": "AthenaPickaxe:DefaultPickaxe",
					"itemProfile": "athena",
					"quantity": 1,
				},
				{
					"itemType": "AthenaGlider:DefaultGlider",
					"itemGuid": "AthenaGlider:DefaultGlider",
					"itemProfile": "athena",
					"quantity": 1,
				},
				{
					"itemType": "AthenaDance:EID_DanceMoves",
					"itemGuid": "AthenaDance:EID_DanceMoves",
					"itemProfile": "athena",
					"quantity": 1,
				},
			},
			"params": gin.H{
				"userMessage": "Server has removed all items from your account. Enjoy!",
			},
			"level": 1,
			"giftedOn": time.Now().Format("2006-01-02T15:04:05.999Z"),
		},
		Quantity: 1,
	}

	commonCore.Items["GiftBox:gb_default"] = gift
	SaveProfileToUser(accountId, commonCore)
}

func SetUserVBucks(accountId string, profile *models.Profile, amount int) {
	_, err := GetUserByAccountId(accountId)
	if err != nil {
		return
	}

	wantedAmount := amount

	all.Postgres.Model(&models.User{}).Where("account_id = ?", accountId).Update("v_bucks", wantedAmount)
	
	profile.Items["Currency:MtxPurchased"] = models.CommonCoreItem{
		TemplateId: "Currency:MtxPurchased",
		Attributes: map[string]any {
			"platform": "EpicPC",
		},
		Quantity: wantedAmount,
	}

	AppendLoadoutsToProfileNoSave(profile, accountId)
}

func TakeUserVBucks(accountId string, profile *models.Profile, amount int) {
	user, err := GetUserByAccountId(accountId)
	if err != nil {
		return
	}

	wantedAmount := user.VBucks - amount

	all.Postgres.Model(&models.User{}).Where("account_id = ?", accountId).Update("v_bucks", wantedAmount)
	
	profile.Items["Currency:MtxPurchased"] = models.CommonCoreItem{
		TemplateId: "Currency:MtxPurchased",
		Attributes: map[string]any {
			"platform": "EpicPC",
		},
		Quantity: wantedAmount,
	}

	AppendLoadoutsToProfileNoSave(profile, accountId)
}

func AddUserVBucks(accountId string, profile *models.Profile, amount int) {
	user, err := GetUserByAccountId(accountId)
	if err != nil {
		return
	}

	wantedAmount := user.VBucks + amount

	all.Postgres.Model(&models.User{}).Where("account_id = ?", accountId).Update("v_bucks", wantedAmount)
	
	profile.Items["Currency:MtxPurchased"] = models.CommonCoreItem{
		TemplateId: "Currency:MtxPurchased",
		Attributes: map[string]any {
			"platform": "EpicPC",
		},
		Quantity: wantedAmount,
	}

	AppendLoadoutsToProfileNoSave(profile, accountId)
}

func GetItemFromProfile(profile *models.Profile, itemId string) (models.Item, error) {
	item, ok := profile.Items[itemId]
	if !ok {
		return models.Item{}, errors.New("item not found")
	}

	marshal, err := json.Marshal(item)
	if err != nil {
		return models.Item{}, err
	}

	var unmarshal models.Item
	err = json.Unmarshal(marshal, &unmarshal)
	if err != nil {
		return models.Item{}, err
	}

	return unmarshal, nil
}

func FindVariant(item *models.Item, channel string) (models.ItemVariant, error) {
	for _, variant := range item.Attributes.Variants {
		if variant.Channel == channel {
			return variant, nil
		}
	}

	return models.ItemVariant{}, errors.New("variant not found")
}

func SetVariantInItem(item *models.Item, variant models.ItemVariant) (models.ItemVariant, error) {
	foundVariant, err := FindVariant(item, variant.Channel)
	if err != nil {
		item.Attributes.Variants = append(item.Attributes.Variants, variant)
		return variant, nil
	}

	foundVariant.Active = variant.Active
	foundVariant.Channel = variant.Channel
	foundVariant.Owned = []string{variant.Active}

	for i, v := range item.Attributes.Variants {
		if v.Channel == variant.Channel {
			item.Attributes.Variants[i] = foundVariant
		}
	}

	return foundVariant, nil
}

func GetFullAthenaProfile(accountId string) models.Profile {
	profile, err := ReadProfileFromUser(accountId, "athena")
	if err != nil {
		return models.Profile{}
	}

	loadout, err := GetLoadout("zombie_loadout", accountId)
	if err != nil {
		return models.Profile{}
	}
	profile.Items["zombie_loadout"] = loadout

	loadout2, err := GetLoadout("sandbox_loadout", accountId)
	if err != nil {
		return models.Profile{}
	}
	profile.Items["sandbox_loadout"] = loadout2

	return profile
}