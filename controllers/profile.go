package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

func ProfileQuery(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	profileId, _ := c.GetQuery("profileId")
	action, _ := c.GetQuery("action")

	profile, err := common.ReadProfileFromUser(user.AccountId, profileId)
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	athenaProfile, nerr := common.ReadProfileFromUser(user.AccountId, "athena")
	if nerr != nil {
		common.ErrorBadRequest(c)
		return
	}

	profile.Rvn += 1
	profile.CommandRevision += 1
	profile.AccountId = user.AccountId
	common.SaveProfileToUSer(user.AccountId, profile)

	athenaProfile.Stats.Attributes.LastAppliedLoadout = athenaProfile.Stats.Attributes.Loadouts[0]
	addblackknight(&athenaProfile, user.AccountId)
	common.SaveProfileToUSer(user.AccountId, athenaProfile)

	switch action {
		case "QueryProfile":
			break
		case "SetMtxPlatform":
			break
		default:
			break
	}

	profile.Stats.Attributes.SeasonNum = common.Season
	
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

func addblackknight(profile *models.Profile, accountId string) {
	profile.Items["CID_035_Athena_Commando_M_Medieval"] = models.Item{
		TemplateId: "AthenaCharacter:CID_035_Athena_Commando_M_Medieval",
		Attributes: models.ItemAttributes{
			MaxLevelBonus: 0,
			Level: 1,
			ItemSeen: true,
			Variants: []string{},
			Favorite: false,
			Xp: 0,
		},
		Quantity: 1,
	}
	common.AppendLoadoutsToProfile(profile, accountId)
}