package socket

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

/*
fortnite/api/game/v2/matchmakingservice/ticket/player/e8f8db85-6a07-411f-80e2-260f4d7f6302
?partyPlayerIds=e8f8db85-6a07-411f-80e2-260f4d7f6302
&bucketId=6037427%3A0%3AEU%3Aplaylist_defaultsolo#
&player.platform=Windows
&player.subregions=IE%2CGB%2CDE%2CFR
&player.option.crossplayOptOut=false
&party.WIN=true&input.KBM=true&player.input=KBM
&player.playerGroups=e8f8db85-6a07-411f-80e2-260f4d7f6302
*/

var (
	matchmakerUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	MatchmakeQueue = make(map[string]MatchmakeInfo)
	Sessions = make(map[string]MatchmakeInfo)
	RemoteAddressToAccountId = make(map[string]string)
)

type MatchmakeInfo struct {
	User models.User
	PlaylistName string
	Region string
	BuildId string
	Authenticated bool
	PositionInQueue int
	SessionId string
}

type TicketBucket struct {
	PlaylistName string `json:"playlistName"`
	Region string `json:"region"`
	BuildId string `json:"buildId"`
}

func MatchmakerHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := matchmakerUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	SendStringJSONToClient(conn, websocket.TextMessage, gin.H{
		"payload": gin.H{
			"state": "Connecting",
		},
		"name": "StatusUpdate",
	})

	initErr := MHandleInit(r.Header.Get("Authorization"), conn)
	if initErr != nil {
		all.PrintRed([]any{"[Matchmaker] ", "Error: ", initErr.Error()})
		SendError(conn, "errors.com.epicgames.common.oauth.invalid_token", "Invalid token")
		return
	}

	// logic for server checking!
	// if !common.IsServerValid(matchmakeInfo.PlaylistName, matchmakeInfo.Region, matchmakeInfo.BuildId) { etc

	SendStringJSONToClient(conn, websocket.TextMessage, gin.H{
		"payload": gin.H{
			"state": "Waiting",
			"totalPlayers": 1,
			"connectedPlayers": 1,
		},
		"name": "StatusUpdate",
	})
	go sendStatusUpdates(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(websocket.TextMessage, []byte("pong"))
	}

	defer func() {
		all.PrintRed([]any{"[Matchmaker] ", "Client disconnected"})
		accountId := RemoteAddressToAccountId[conn.RemoteAddr().String()]
		delete(MatchmakeQueue, accountId)
		delete(RemoteAddressToAccountId, conn.RemoteAddr().String())
		conn.Close()
	}()
}

func sendStatusUpdates(client *websocket.Conn) {
	accountId := RemoteAddressToAccountId[client.RemoteAddr().String()]
	matchmakeInfo := MatchmakeQueue[accountId]
	
	playersQueued := len(MatchmakeQueue)
	infront := playersQueued - matchmakeInfo.PositionInQueue
	etaSeconds := (infront%100) + 1

	servers := common.GameServers
	var leastPlayersServer common.GameServer
	for _, server := range servers {
		if leastPlayersServer.IP == "" || server.PlayersLeft < leastPlayersServer.PlayersLeft {
			leastPlayersServer = server
		}
	}
	etaSeconds += leastPlayersServer.PlayersLeft * 10
	
	if leastPlayersServer.IP == "" {
		etaSeconds = 0
	}

	for {
		playersQueued := len(MatchmakeQueue)
		time.Sleep(250 * time.Millisecond)

		if matchmakeInfo.PositionInQueue < 0 {
			matchmakeInfo.PositionInQueue = playersQueued
			MatchmakeQueue[accountId] = matchmakeInfo
		}

		if RemoteAddressToAccountId[client.RemoteAddr().String()] != accountId {
			break
		}

		infront := playersQueued - matchmakeInfo.PositionInQueue
		
		SendStringJSONToClient(client, websocket.TextMessage, gin.H{
			"payload": gin.H{
				"state": "Queued",
				"ticketId": "ticketId",
				"estimatedWaitSec": etaSeconds,
				"queuedPlayers": infront,
				"status": gin.H{},
			},
			"name": "StatusUpdate",
		})
		
		all.PrintMagenta([]any{"[Matchmaker] ", "Queued: ", infront, " ETA: ", etaSeconds})
		
		if infront != 0 {
			continue
		}

		gameServerWantingToJoin := common.GameServers[matchmakeInfo.PlaylistName + ":" + matchmakeInfo.Region]

		if gameServerWantingToJoin.Joinable {
			etaSeconds = 0
			break
		}
	}

	SendStringJSONToClient(client, websocket.TextMessage, gin.H{
		"payload": gin.H{
			"state": "SessionAssignment",
			"matchId": matchmakeInfo.PlaylistName + ":" + matchmakeInfo.Region,
		},
		"name": "StatusUpdate",
	})

	matchmakeInfo.PositionInQueue = -1
	MatchmakeQueue[accountId] = matchmakeInfo

	SendStringJSONToClient(client, websocket.TextMessage, gin.H{
		"payload": gin.H{
			"matchId": matchmakeInfo.PlaylistName + ":" + matchmakeInfo.Region,
			"sessionId": matchmakeInfo.SessionId,
			"joinDelaySec": 0,
		},
		"name": "Play",
	})

	all.PrintGreen([]any{"closing queue stream"})

	client.Close()
	delete(MatchmakeQueue, RemoteAddressToAccountId[client.RemoteAddr().String()])
	delete(RemoteAddressToAccountId, client.RemoteAddr().String())
}

func MHandleInit(authHeader string, client *websocket.Conn) error {
	matchmakeData := strings.Split(strings.Split(authHeader, "  ")[1], " ")
	user, err := common.GetUserByAccountId(matchmakeData[0])
	if err != nil {
		return err
	}

	var bucket TicketBucket
	json.Unmarshal([]byte(matchmakeData[1]), &bucket)
	
	matchmakeInfo := &MatchmakeInfo{
		User: user,
		PlaylistName: bucket.PlaylistName,
		Region: bucket.Region,
		BuildId: bucket.BuildId,
		Authenticated: true,
		PositionInQueue: -1,
		SessionId: uuid.New().String(),
	}
	MatchmakeQueue[user.AccountId] = *matchmakeInfo
	RemoteAddressToAccountId[client.RemoteAddr().String()] = user.AccountId
	Sessions[matchmakeInfo.SessionId] = *matchmakeInfo

	return nil
}

func SendError(client *websocket.Conn, errName string, errMessage string) {
	SendStringJSONToClient(client, 1, gin.H{
		"payload": gin.H{
			"state": "Error",
			"error": errName,
			"errorMessage": errMessage,
		},
	})
}

func SendStringJSONToClient(client *websocket.Conn, messageType int, messageData gin.H) {
	data, _ := json.Marshal(messageData)

	err := client.WriteMessage(messageType, data)
	if err != nil {
		all.PrintRed([]any{"CLOSSING CONNECTION", err.Error()})
		client.Close()
		delete(MatchmakeQueue, RemoteAddressToAccountId[client.RemoteAddr().String()])
		delete(RemoteAddressToAccountId, client.RemoteAddr().String())
	}
}