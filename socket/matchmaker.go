package socket

import (
	"encoding/json"
	"math/rand"
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

var (
	matchmakerUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	MatchmakeQueue = make(map[string]MatchmakeInfo)
	Sessions = make(map[string]MatchmakeInfo)
	RemoteAddressToAccountId = make(map[string]string)

	secondTimer = 0
	fakePlayersToInflateETA = 0
)

type MatchmakeInfo struct {
	User models.User
	PlaylistName string
	Region string
	BuildId string
	Authenticated bool
	PositionInQueue int
	SessionId string
	Connection *websocket.Conn
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

	SendStringJSONToClient(conn, websocket.TextMessage, gin.H{
		"payload": gin.H{
			"state": "Waiting",
			"totalPlayers": 1,
			"connectedPlayers": 1,
		},
		"name": "StatusUpdate",
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(websocket.TextMessage, []byte("pong"))
	}

	defer func() {
		all.PrintRed([]any{"[Matchmaker] ", "Client disconnected"})
		conn.Close()
		accountId := RemoteAddressToAccountId[conn.RemoteAddr().String()]
		delete(MatchmakeQueue, accountId)
		delete(RemoteAddressToAccountId, conn.RemoteAddr().String())
	}()
}

func calculateETA(matchmakeInfo *MatchmakeInfo) int {
	gameServers := common.GetAllGameServers(matchmakeInfo.PlaylistName, matchmakeInfo.Region)
	playersQueued := len(MatchmakeQueue) + fakePlayersToInflateETA
	infront := playersQueued - matchmakeInfo.PositionInQueue
	
	return int((infront * 10) / len(gameServers))
}

func sendStatusUpdates() {
	for {
		time.Sleep(250 * time.Millisecond)
		secondTimer += 250

		if secondTimer % 10000 == 0 {
			secondTimer = 0
			fakePlayersToInflateETA -= rand.Intn(10)
		}

		for _, matchmakeInfo := range MatchmakeQueue {
			shouldLoad := false
			etaSeconds := calculateETA(&matchmakeInfo)
			all.MarshPrintJSON(RemoteAddressToAccountId)
			playersQueued := len(MatchmakeQueue)

			if matchmakeInfo.PositionInQueue < 0 {
				matchmakeInfo.PositionInQueue = playersQueued
				MatchmakeQueue[matchmakeInfo.User.AccountId] = matchmakeInfo
			}

			infront := playersQueued - matchmakeInfo.PositionInQueue
			
			SendStringJSONToClient(matchmakeInfo.Connection, websocket.TextMessage, gin.H{
				"payload": gin.H{
					"state": "Queued",
					"ticketId": "ticketId",
					"estimatedWaitSec": etaSeconds + 1,
					"queuedPlayers": infront,
					"status": gin.H{},
				},
				"name": "StatusUpdate",
			})
			
			all.PrintMagenta([]any{"[Matchmaker] ", "Queued: ", infront, " ETA: ", etaSeconds})

			gameServerWantingToJoin := common.GetGameServer(matchmakeInfo.PlaylistName, matchmakeInfo.Region)
			if gameServerWantingToJoin.Joinable && gameServerWantingToJoin.PlayersLeft < 100 && infront < 1 {
				etaSeconds = 1
				shouldLoad = true
			}

			if shouldLoad {
				sendPlayMessage(&matchmakeInfo)
				continue
			}
		}
	}
}

func sendPlayMessage(matchmakeInfo *MatchmakeInfo) {
	SendStringJSONToClient(matchmakeInfo.Connection, websocket.TextMessage, gin.H{
		"payload": gin.H{
			"state": "SessionAssignment",
			"matchId": matchmakeInfo.PlaylistName + ":" + matchmakeInfo.Region,
		},
		"name": "StatusUpdate",
	})

	matchmakeInfo.PositionInQueue = -1
	MatchmakeQueue[matchmakeInfo.User.AccountId] = *matchmakeInfo

	SendStringJSONToClient(matchmakeInfo.Connection, websocket.TextMessage, gin.H{
		"payload": gin.H{
			"matchId": matchmakeInfo.PlaylistName + ":" + matchmakeInfo.Region,
			"sessionId": matchmakeInfo.SessionId,
			"joinDelaySec": 0,
		},
		"name": "Play",
	})
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
		Connection: client,
	}
	MatchmakeQueue[user.AccountId] = *matchmakeInfo
	RemoteAddressToAccountId[client.RemoteAddr().String()] = user.AccountId
	Sessions[matchmakeInfo.SessionId] = *matchmakeInfo

	all.PrintGreen([]any{"[Matchmaker] ", "New matchmake request from ", user.AccountId, " (", client.RemoteAddr().String(), ")"})
	all.MarshPrintJSON(RemoteAddressToAccountId)

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
	client.WriteMessage(messageType, data)
}

func InitMatchmaker() {
	go sendStatusUpdates()
}