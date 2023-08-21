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

		clientRemoteAddress, _ := socket.AccountIdToXMPPRemoteAddress[wantedFriend]
		client, _ := socket.ActiveXMPPClients[clientRemoteAddress]

		friendRemoteAddress, _ := socket.AccountIdToXMPPRemoteAddress[user.AccountId]
		friend, _ := socket.ActiveXMPPClients[friendRemoteAddress]

		client.Connection.WriteMessage(1, []byte(`
			<presence to="`+ client.JID +`" xmlns="jabber:client" from="`+ friend.JID +`" type="available">
				<status>`+ friend.Status +`</status>
			</presence>
		`))

		friend.Connection.WriteMessage(1, []byte(`
			<presence to="`+ friend.JID +`" xmlns="jabber:client" from="`+ client.JID +`" type="available">
				<status>`+ client.Status +`</status>
			</presence>
		`))
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

type FriendSummaryItem struct {
	AccountId string `json:"accountId"`
	Alias string `json:"alias"`
	Mutual int `json:"mutual"`
	Note string `json:"note"`
	Groups []string `json:"groups"`
	Favorite bool `json:"favorite"`
	Created string `json:"created"`
}

type FriendSummary struct {
	Friends []FriendSummaryItem `json:"friends"`
	Blocklist []FriendSummaryItem `json:"blocklist"`
	Incoming []FriendSummaryItem `json:"incoming"`
	Outgoing []FriendSummaryItem `json:"outgoing"`
	Suggested []FriendSummaryItem `json:"suggested"`
	Settings gin.H `json:"settings"`
}

func FriendsSummary(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	response := FriendSummary{
		Friends: []FriendSummaryItem{},
		Blocklist: []FriendSummaryItem{},
		Incoming: []FriendSummaryItem{},
		Outgoing: []FriendSummaryItem{},
		Suggested: []FriendSummaryItem{},
		Settings: gin.H{},
	}

	friends := common.GetAllAcceptedFriends(user.AccountId)
	pending := common.GetPendingFriendsList(user.AccountId)
	blocked := common.GetBlockedFriendsList(user.AccountId)

	for _, friend := range friends {
		friendAccount, _ := common.GetUserByAccountId(friend.AccountId)

		response.Friends = append(response.Friends, FriendSummaryItem{
			AccountId: friendAccount.AccountId,
			Alias: friendAccount.Username,
			Mutual: 0,
			Note: "",
			Groups: []string{},
			Favorite: false,
			Created: friend.Created,
		})
	}

	for _, friend := range pending {
		friendAccount, _ := common.GetUserByAccountId(friend.AccountId)

		if friend.Direction == "INBOUND" {
			response.Incoming = append(response.Friends, FriendSummaryItem{
				AccountId: friendAccount.AccountId,
				Alias: friendAccount.Username,
				Mutual: 0,
				Note: "",
				Groups: []string{},
				Favorite: false,
				Created: friend.Created,
			})
		}

		if friend.Direction == "OUTBOUND" {
			response.Outgoing = append(response.Friends, FriendSummaryItem{
				AccountId: friendAccount.AccountId,
				Alias: friendAccount.Username,
				Mutual: 0,
				Note: "",
				Groups: []string{},
				Favorite: false,
				Created: friend.Created,
			})
		}
	}

	for _, friend := range blocked {
		friendAccount, _ := common.GetUserByAccountId(friend.AccountId)

		response.Blocklist = append(response.Blocklist, FriendSummaryItem{
			AccountId: friendAccount.AccountId,
			Alias: friendAccount.Username,
			Mutual: 0,
			Note: "",
			Groups: []string{},
			Favorite: false,
			Created: friend.Created,
		})
	}

	c.JSON(200, response)
}