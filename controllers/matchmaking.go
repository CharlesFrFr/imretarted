package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
	"github.com/zombman/server/socket"
)

/*
	fortnite/api/game/v2/matchmakingservice/ticket/player/e8f8db85-6a07-411f-80e2-260f4d7f6302
	?partyPlayerIds=e8f8db85-6a07-411f-80e2-260f4d7f6302
	&bucketId=6037427%3A0%3ANAE%3Aplaylist_defaultsolo
	&player.platform=Windows&player.subregions=VA%2COH
	&player.option.crossplayOptOut=false
	&party.WIN=true
	&input.KBM=true
	&player.input=KBM
	&player.playerGroups=e8f8db85-6a07-411f-80e2-260f4d7f6302
*/

// bucket 6037427:0:NAE:playlist_defaultsolo

type TicketBucket struct {
	PlaylistName string `json:"playlistName"`
	Region string `json:"region"`
	BuildId string `json:"buildId"`
}

var UserBuilds = make(map[string]string)

func MatchmakingTicket(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	accountId := c.Param("accountId")

	if accountId != user.AccountId {
		common.ErrorBadRequest(c)
		return
	}

	customKey := c.Param("player.option.customKey")
	if customKey != "" {
		common.ErrorBadRequest(c)
		return
	}

	bucketInformation := strings.Split(c.Query("bucketId"), ":")

	bucket := TicketBucket{
		PlaylistName: bucketInformation[3],
		Region: bucketInformation[2],
		BuildId: bucketInformation[0],
	}

	bucketString, _ := json.Marshal(bucket)

	_, ok := common.GameServers[bucket.PlaylistName + ":" + bucket.Region]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	UserBuilds[user.AccountId] = bucket.BuildId

	c.JSON(http.StatusOK, gin.H{
		"ticket-type": "mms-player",
		"signature": string(bucketString),
		"payload": accountId,
		"serviceUrl": "ws://" + common.IP + "/match",
	})
}

func GetMatchmakeSession(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	sessionId := c.Param("sessionId")

	if sessionId == "" {
		common.ErrorBadRequest(c)
		return
	}

	matchmakeInfo, ok := socket.Sessions[sessionId]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	gameServer, ok := common.GameServers[matchmakeInfo.PlaylistName + ":" + matchmakeInfo.Region]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	if matchmakeInfo.User.AccountId != user.AccountId {
		common.ErrorBadRequest(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": matchmakeInfo.SessionId,
		"ownerId": uuid.New().String(),
		"ownerName": "zombieman",
		"serverName": "zombieman",
		"serverAddress": gameServer.IP,
		"serverPort": gameServer.Port,
		"maxPublicPlayers": 100,
		"openPublicPlayers": 100,
		"maxPrivatePlayers": 0,
		"openPrivatePlayers": 0,
		"attributes": gin.H{
			"REGION_s": gameServer.Region,
			"GAMEMODE_s": "FORTATHENA",
			"ALLOWBROADCASTING_b": true,
			"SUBREGION_s": "",
			"DCID_s": "zombieman",
			"tenant_s": "Fortnite",
			"MATCHMAKINGPOOL_s": "Any",
			"STORMSHIELDDEFENSETYPE_i": 0,
			"HOTFIXVERSION_i": 0,
			"PLAYLISTNAME_s": matchmakeInfo.PlaylistName,
			"SESSIONKEY_s": uuid.New().String(),
			"TENANT_s": "Fortnite",
			"BEACONPORT_i": 15009,
		},
		"publicPlayers": []string{},
		"privatePlayers": []string{},
		"totalPlayers": gameServer.PlayersLeft,
		"allowJoinInProgress": false,
		"shouldAdvertise": false,
		"isDedicated": false,
		"usesStats": false,
		"allowInvites": false,
		"usesPresence": false,
		"allowJoinViaPresence": true,
		"allowJoinViaPresenceFriendsOnly": false,
		"buildUniqueId": UserBuilds[user.AccountId],
		"lastUpdated": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		"started": !gameServer.Joinable,
	})
}

func GetMatchmakingKey(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	sessionId := c.Param("sessionId")

	if sessionId == "" {
		common.ErrorBadRequest(c)
		return
	}

	matchmakeInfo, ok := socket.Sessions[sessionId]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	_, ok = common.GameServers[matchmakeInfo.PlaylistName + ":" + matchmakeInfo.Region]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	if matchmakeInfo.User.AccountId != user.AccountId {
		common.ErrorBadRequest(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key": "none",
		"accountId": user.AccountId,
		"sessionId": sessionId,
	})
}