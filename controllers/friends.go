package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
	"github.com/zombman/server/socket"
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
	user := c.MustGet("user").(models.User)
	wantedFriend := c.Param("friendId")

	if wantedFriend == user.AccountId {
		common.ErrorBadRequest(c)
		return
	}

	if common.IsFriend(user.AccountId, wantedFriend) {
		common.ErrorBadRequest(c)
		return
	}

	if common.IsBlocked(user.AccountId, wantedFriend) {
		common.ErrorBadRequest(c)
		return
	}

	if common.IsPending(user.AccountId, wantedFriend) {
		common.ErrorBadRequest(c)
		return
	}

	res := common.CreateFriend(user.AccountId, wantedFriend)

	if res == "ACCEPTED" {
		socket.XMPPSendBodyToAccountId(gin.H{
			"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"payload": gin.H{
				"accountId": user.AccountId,
				"status": "ACCEPTED",
				"direction": "INBOUND",
				"favorite": false,
				"created": time.Now().Format("2006-01-02T15:04:05.999Z"),
			},
		}, wantedFriend)

		socket.XMPPSendBodyToAccountId(gin.H{
			"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"payload": gin.H{
				"accountId": wantedFriend,
				"status": "ACCEPTED",
				"direction": "OUTBOUND",
				"favorite": false,
				"created": time.Now().Format("2006-01-02T15:04:05.999Z"),
			},
		}, user.AccountId)
	}

	if res == "PENDING" {
		socket.XMPPSendBodyToAccountId(gin.H{
			"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"payload": gin.H{
				"accountId": user.AccountId,
				"status": "PENDING",
				"direction": "INBOUND",
				"favorite": false,
				"created": time.Now().Format("2006-01-02T15:04:05.999Z"),
			},
		}, wantedFriend)

		socket.XMPPSendBodyToAccountId(gin.H{
			"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"payload": gin.H{
				"accountId": wantedFriend,
				"status": "PENDING",
				"direction": "OUTBOUND",
				"favorite": false,
				"created": time.Now().Format("2006-01-02T15:04:05.999Z"),
			},
		}, user.AccountId)
	}

	c.Status(204)
}

func DeleteFriend(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	wantedFriend := c.Param("friendId")

	if wantedFriend == user.AccountId {
		all.PrintRed([]any{"TRYING TO DELETE HIMSELF"})

		common.ErrorBadRequest(c)
		return
	}

	if !common.IsFriend(user.AccountId, wantedFriend) {
		all.PrintRed([]any{"NOT FRIENDS"})

		common.ErrorBadRequest(c)
		return
	}

	res := common.DeleteFriend(user.AccountId, wantedFriend)
	if res == "DELETED" {
		socket.XMPPSendBodyToAccountId(gin.H{
			"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
			"type": "com.epicgames.friends.core.apiobjects.FriendRemoval",
			"payload": gin.H{
				"accountId": user.AccountId,
				"reason": "DELETED",
			},
		}, wantedFriend)

		socket.XMPPSendBodyToAccountId(gin.H{
			"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
			"type": "com.epicgames.friends.core.apiobjects.FriendRemoval",
			"payload": gin.H{
				"accountId": wantedFriend,
				"reason": "DELETED",
			},
		}, user.AccountId)
	}

	c.Status(204)
}

func BlockFriend(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	wantedFriend := c.Param("friendId")

	if wantedFriend == user.AccountId {
		all.PrintRed([]any{"TRYING TO DELETE HIMSELF"})

		common.ErrorBadRequest(c)
		return
	}

	if !common.IsFriend(user.AccountId, wantedFriend) {
		all.PrintRed([]any{"NOT FRIENDS"})

		common.ErrorBadRequest(c)
		return
	}

	common.BlockFriend(user.AccountId, wantedFriend)
	
	c.Status(204)
}

func UnBlockFriend(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	wantedFriend := c.Param("friendId")

	if wantedFriend == user.AccountId {
		all.PrintRed([]any{"TRYING TO DELETE HIMSELF"})

		common.ErrorBadRequest(c)
		return
	}

	common.UnBlockFriend(user.AccountId, wantedFriend)
	
	c.Status(204)
}

type MatchUser struct {
	AccountId string `json:"accountId"`
	EpicMutuals int `json:"epicMutuals"`
	SortPosition int `json:"sortPosition"`
	MatchType string `json:"matchType"`
	Matches []Match `json:"matches"`
}

type Match struct {
	Value string `json:"value"`
	Platform string `json:"platform"`
}

func SearchForUser(c *gin.Context) {
	var users []MatchUser
	prefix := c.Query("prefix")
	
	var databaseMatches []models.User
	all.Postgres.Model(&models.User{}).Where("username LIKE ?", prefix + "%").Limit(10).Find(&databaseMatches)

	for i, match := range databaseMatches {
		users = append(users, MatchUser{
			AccountId: match.AccountId,
			EpicMutuals: 0,
			SortPosition: i,
			MatchType: "PREFIX",
			Matches: []Match{{
				Value: match.Username,
				Platform: "WIN",
			}},
		})
	}

	c.JSON(200, users)
}