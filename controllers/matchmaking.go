package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
	"github.com/zombman/server/socket"
)

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
		all.PrintRed([]any{"account id not match"})
		common.ErrorBadRequest(c)
		return
	}

	customKey := c.Param("player.option.customKey")
	if customKey != "" {
		all.PrintRed([]any{"custom key"})
		common.ErrorBadRequest(c)
		return
	}

	bucketInformation := strings.Split(c.Query("bucketId"), ":")
	fmt.Println(bucketInformation[2])
	bucket := TicketBucket{
		PlaylistName: bucketInformation[3],
		Region: bucketInformation[2],
		BuildId: bucketInformation[0],
	}

	bucketString, _ := json.Marshal(bucket)

	_, ok := common.GameServers[bucket.PlaylistName + ":" + bucket.Region]
	if !ok {
		all.PrintRed([]any{"servers not found", bucket.PlaylistName + ":" + bucket.Region})
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

	gameServer := common.GetGameServer(matchmakeInfo.PlaylistName, matchmakeInfo.Region)
	if gameServer.IP == "" {
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

func JoinMatchmakeSession(c *gin.Context) {
	sessionId := c.Param("sessionId")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"sessionId": sessionId,
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

func AddNewGameServer(c *gin.Context) {
	var body struct {
		IP string `json:"ip"`
		Port int `json:"port"`
		Region string `json:"region"`
		PlaylistName string `json:"playlistName"`
		MaxPlayers int `json:"maxPlayers"`
		PlayersLeft int `json:"playersLeft"`
		Joinable bool `json:"joinable"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		common.ErrorBadRequest(c)
		return
	}

	for idx, gameServerArray := range common.GameServers[body.PlaylistName + ":" + body.Region] {
		if gameServerArray.IP == body.IP && gameServerArray.Port == body.Port {
			common.GameServers[body.PlaylistName + ":" + body.Region] = append(common.GameServers[body.PlaylistName + ":" + body.Region][:idx], common.GameServers[body.PlaylistName + ":" + body.Region][idx+1:]...)
		}
	}

	common.GameServers[body.PlaylistName + ":" + body.Region] = append(common.GameServers[body.PlaylistName + ":" + body.Region], common.GameServer{
		IP: body.IP,
		Port: body.Port,
		Region: body.Region,
		Playlist: body.PlaylistName,
		PlayersLeft: body.PlayersLeft,
		Joinable: body.Joinable,
	})

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func RemoveGameServer(c *gin.Context) {
	var body struct {
		IP string `json:"ip"`
		Port int `json:"port"`
		Region string `json:"region"`
		PlaylistName string `json:"playlistName"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		common.ErrorBadRequest(c)
		return
	}

	for idx, gameServerArray := range common.GameServers[body.PlaylistName + ":" + body.Region] {
		if gameServerArray.IP == body.IP && gameServerArray.Port == body.Port {
			common.GameServers[body.PlaylistName + ":" + body.Region] = append(common.GameServers[body.PlaylistName + ":" + body.Region][:idx], common.GameServers[body.PlaylistName + ":" + body.Region][idx+1:]...)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}