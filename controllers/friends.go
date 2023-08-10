package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

func FriendsPublic(c *gin.Context) {
	all.PrintMagenta([]any{
		"GET /friends for",
		c.MustGet("user").(models.User).AccountId,
	})
	user := c.MustGet("user").(models.User)
	friends := common.GetFriendsList(user.AccountId)

	all.MarshPrintJSON(friends)

	c.JSON(200, friends)
}

func FriendsBlocked(c *gin.Context) {
	all.PrintMagenta([]any{
		"GET BLOCKED friends for",
		c.MustGet("user").(models.User).AccountId,
	})
	user := c.MustGet("user").(models.User)
	friends := common.GetBlockedFriendsList(user.AccountId)

	all.MarshPrintJSON(friends)

	c.JSON(200, friends)
}


func CreateFriend(c *gin.Context) {
	all.PrintMagenta([]any{
		"CREATE FRIENF for",
		c.MustGet("user").(models.User).AccountId,
	})
	user := c.MustGet("user").(models.User)
	wantedFriend := c.Param("friendId")

	if wantedFriend == user.AccountId {
		all.PrintRed([]any{
			"tried to friend self",
			user.AccountId,
		})
		common.ErrorBadRequest(c)
		return
	}

	if common.IsFriend(user.AccountId, wantedFriend) {
		all.PrintRed([]any{
			"already friends with",
			wantedFriend,
		})
		common.ErrorBadRequest(c)
		return
	}

	if common.IsBlocked(user.AccountId, wantedFriend) {
		all.PrintRed([]any{
			"blocked",
			wantedFriend,
		})
		common.ErrorBadRequest(c)
		return
	}

	if common.IsPending(user.AccountId, wantedFriend) {
		all.PrintRed([]any{
			"already pending",
			wantedFriend,
		})
		common.ErrorBadRequest(c)
		return
	}

	common.CreateFriend(user.AccountId, wantedFriend)
	// if err != nil {
	// 	common.ErrorInternalServer(c)
	// 	return
	// }

	all.PrintMagenta([]any{
		"sent friend request from",
		user.AccountId,
		"to",
		wantedFriend,
		})

	c.Status(204)
}