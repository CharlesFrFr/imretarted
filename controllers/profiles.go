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
	if profileId == "athena" { profile.Stats.Attributes.LastAppliedLoadout = profile.Stats.Attributes.Loadouts[0] }
	common.SaveProfileToUSer(user.AccountId, profile)

	switch action {
		case "QueryProfile":
			break
		case "SetMtxPlatform":
			break
		default:
			break
	}

	profile.Stats.Attributes.SeasonNum = 6
	
	c.JSON(200, gin.H{
		"profileRevision": profile.Rvn,
		"profileId": profileId,
		"profileChangesBaseRevision": profile.Rvn,
		"profileChanges": []gin.H{{
			"changeType": "fullProfileUpdate",
			"profile": profile,
		}},
		"serverTime": time.Now().Format(time.RFC3339),
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