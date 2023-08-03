package controllers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

func ProfileActionHandler(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	profileId, _ := c.GetQuery("profileId")
	action := c.Param("action")

	profile, err := common.ReadProfileFromUser(user.AccountId, profileId)
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	profile.Rvn += 1
	profile.CommandRevision += 1
	profile.AccountId = user.AccountId
	common.SaveProfileToUser(user.AccountId, profile)

	switch action {
		case "QueryProfile":
			break
		case "SetMtxPlatform":
			break
		case "ClientQuestLogin":
			break
		case "EquipBattleRoyaleCustomization":
			EquipBattleRoyaleCustomization(c, user, &profile)
		default:
			all.PrintRed([]string{"unknown action", action})
			common.ErrorBadRequest(c)
			return
	}

	athenaProfile, nerr := common.ReadProfileFromUser(user.AccountId, "athena")
	if nerr != nil {
		common.ErrorBadRequest(c)
		return
	}

	profile.Stats.Attributes.SeasonNum = common.Season
	athenaProfile.Stats.Attributes.SeasonNum = common.Season
	athenaProfile.Stats.Attributes.LastAppliedLoadout = athenaProfile.Stats.Attributes.Loadouts[0]

	AddItemToProfile(&athenaProfile, "AthenaCharacter:CID_024_Athena_Commando_F", user.AccountId)
	AddItemToProfile(&athenaProfile, "AthenaBackpack:BID_003_RedKnight", user.AccountId)
	AddItemToProfile(&athenaProfile, "AthenaPickaxe:Pickaxe_ID_015_HolidayCandyCane", user.AccountId)
	AddItemToProfile(&athenaProfile, "AthenaGlider:Umbrella_Platinum", user.AccountId)
	AddItemToProfile(&athenaProfile, "AthenaSkyDiveContrail:Trails_ID_003_Fire", user.AccountId)
	AddItemToProfile(&athenaProfile, "AthenaItemWrap:Wrap_004_DurrBurgerPJs", user.AccountId)
	
	c.JSON(200, gin.H{
		"profileRevision": profile.Rvn,
		"profileId": profileId,
		"profileChangesBaseRevision": profile.Rvn,
		"profileChanges": []gin.H{{
			"changeType": "fullProfileUpdate",
			"profile": profile,
		}},
		"serverTime": time.Now().Format("2006-01-02T15:04:05.999Z"),
		"multiUpdate": []gin.H{{
			"profileRevision": athenaProfile.Rvn,
			"profileId": athenaProfile.ProfileId,
			"profileChangesBaseRevision": athenaProfile.Rvn,
			"profileChanges": []gin.H{{
				"changeType": "fullProfileUpdate",
				"profile": athenaProfile,
			}},
			"profileCommandRevision": athenaProfile.CommandRevision,
		}},
	})
}

func EquipBattleRoyaleCustomization(c *gin.Context, user models.User, profile *models.Profile) {
	if profile.ProfileId != "athena" {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}
	var body struct {
		SlotName string `json:"slotName"` //"slotName": "Character",
		ItemToSlot string `json:"itemToSlot"` // "itemToSlot": "AthenaCharacter:CID_008_Athena_Commando_M_Default",
		IndexWithinSlot int `json:"indexWithinSlot"`
		VariantUpdates []map[string]interface{} `json:"variantUpdates"`
	}

	if err := c.ShouldBind(&body); err != nil {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}


	foundItem := profile.Items[body.ItemToSlot]
	if foundItem == nil {
		all.PrintRed([]string{"could not find item", body.ItemToSlot})
		common.ErrorItemNotFound(c)
		c.Abort()
		return
	}

	activeLoadoutId := profile.Stats.Attributes.Loadouts[profile.Stats.Attributes.ActiveLoadoutIndex]
	activeLoadout, err := common.GetLoadout(activeLoadoutId, user.AccountId)
	if err != nil {
		common.ErrorItemNotFound(c)
		c.Abort()
		return
	}

	lowercaseItemType := strings.ToLower(body.SlotName)

	switch lowercaseItemType {
		case "character":
			profile.Stats.Attributes.FavoriteCharacter = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Character"].Items[0] = body.ItemToSlot
		case "backpack":
			profile.Stats.Attributes.FavoriteBackpack = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Backpack"].Items[0] = body.ItemToSlot
		case "pickaxe":
			profile.Stats.Attributes.FavoritePickaxe = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Pickaxe"].Items[0] = body.ItemToSlot
		case "glider":
			profile.Stats.Attributes.FavoriteGlider = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Glider"].Items[0] = body.ItemToSlot
		case "skydivecontrail":
			profile.Stats.Attributes.FavoriteSkyDiveContrail = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["SkyDiveContrail"].Items[0] = body.ItemToSlot
		case "loadingscreen":
			profile.Stats.Attributes.FavoriteLoadingScreen = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["LoadingScreen"].Items[0] = body.ItemToSlot
		case "musicpack":
			profile.Stats.Attributes.FavoriteMusicPack = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["MusicPack"].Items[0] = body.ItemToSlot
		case "dance":
			profile.Stats.Attributes.FavoriteDance[body.IndexWithinSlot] = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Dance"].Items[body.IndexWithinSlot] = body.ItemToSlot
		case "itemwrap":
			profile.Stats.Attributes.FavoriteItemWraps[body.IndexWithinSlot] = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["ItemWrap"].Items[body.IndexWithinSlot] = body.ItemToSlot
		default:
			all.PrintRed([]string{"unknown item type", lowercaseItemType})
			common.ErrorBadRequest(c)
			c.Abort()
	}

	profile.Rvn += 1
	profile.CommandRevision += 1
	profile.Updated = time.Now().Format("2006-01-02T15:04:05.999Z")

	common.AppendLoadoutToProfile(profile, &activeLoadout, user.AccountId)
}

func AddItemToProfile(profile *models.Profile, itemId string, accountId string) {
	profile.Items[itemId] = models.Item{
		TemplateId: itemId,
		Attributes: models.ItemAttributes{
			MaxLevelBonus: 0,
			Level: 1,
			ItemSeen: true,
			Variants: []string{},
			Favorite: false,
			Xp: 0,
		},
	}
	common.AppendLoadoutsToProfile(profile, accountId)
}