package socket

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/middleware"
	"github.com/zombman/server/models"
)

var (
	xmppUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ActiveXMPPClients = make(map[string]*ClientInfo)
	AccountIdToXMPPRemoteAddress = make(map[string]string)
)

type ClientInfo struct {
	UUID string
	SocketID string
	Resource string
	Connection *websocket.Conn
	Authenticated bool
	User models.User
	Status string
}

func XMPPHandler(w http.ResponseWriter, r *http.Request){
	conn, err := xmppUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	remoteAddress := conn.RemoteAddr().String()

	clientInfo := &ClientInfo{
		UUID: uuid.New().String(),
		Authenticated: false,
		User: models.User{},
		SocketID: "",
		Connection: conn,
	}
	ActiveXMPPClients[remoteAddress] = clientInfo

	for {
		_, ok := AccountIdToXMPPRemoteAddress[clientInfo.User.AccountId]
		if clientInfo.Authenticated && !ok {
			AccountIdToXMPPRemoteAddress[clientInfo.User.AccountId] = remoteAddress
		}

		messageType, messageData, err := conn.ReadMessage()
		if err != nil {
			break
		}

		all.PrintGreen([]any{
			string(messageData),
		})

		if xml.Unmarshal([]byte(messageData), &models.OpenXML{}) == nil {
			XHandleOpen(conn, messageData, messageType, clientInfo)
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.AuthXML{}) == nil {
			XHandleAuth(conn, messageData, messageType, clientInfo)
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.BindIQXML{}) == nil {
			var iq models.BindIQXML
			xml.Unmarshal([]byte(messageData), &iq)

			if iq.ID == "_xmpp_bind1" {
				XHandleBindIQ(conn, messageData, messageType, clientInfo)
				continue
			}
			
			XHandleSessionIQ(conn, messageData, messageType, clientInfo, iq.ID)
			
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.MessageWithBodyXML{}) == nil {
			XHandleMessage(conn, messageData, messageType, clientInfo)
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.PartyPresenceXML{}) == nil {
			var presence models.PartyPresenceXML
			xml.Unmarshal([]byte(messageData), &presence)

			if presence.To == "" {
				XHandlePresence(conn, messageData, messageType, clientInfo)
				continue
			}
	
			XHandlePartyPresence(conn, messageData, messageType, clientInfo)
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.CloseXML{}) == nil {
			conn.Close()
			continue
		}

		all.PrintRed([]any{
			"unkown message type",
			string(messageData),
		})

		conn.WriteMessage(messageType, []byte(""))
	}

	defer func() {
		all.PrintGreen([]any{"user logged out:", clientInfo.User.Username})
		delete(ActiveXMPPClients, remoteAddress)
		delete(AccountIdToXMPPRemoteAddress, clientInfo.User.AccountId)
		conn.Close()
	}()
}

func XHandleMessage(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	var msg models.MessageWithBodyXML
	xml.Unmarshal([]byte(message), &msg)

	messageRecipientAccountId := strings.Split(msg.To, "@")[0]
	recipientClient, err := XGetClientFromAccountId(messageRecipientAccountId)

	if err != nil {
		return
	}

	if msg.Type == "chat" {
		recipientClient.Connection.WriteMessage(messageType, []byte(`
			<message to="`+ recipientClient.SocketID +`" from="`+ clientInfo.SocketID +`" id="`+ msg.ID +`" xmlns="jabber:client" type="chat">
				<body>`+ msg.Body.Value +`</body>
			</message>
		`))
		return
	}

	recipientClient.Connection.WriteMessage(messageType, []byte(`
		<message to="`+ recipientClient.SocketID +`" from="`+ clientInfo.SocketID +`" id="`+ msg.ID +`" xmlns="jabber:client">
			<body>`+ msg.Body.Value +`</body>
		</message>
	`))
}

func XHandleSessionIQ(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo, id string) {
	conn.WriteMessage(messageType, []byte(`
		<iq to="`+ clientInfo.SocketID +`" from="prod.ol.epicgames.com" id="`+ id +`" xmlns="jabber:client" type="result" />
	`))
}

func XHandleBindIQ(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	var iq models.BindIQXML
	xml.Unmarshal([]byte(message), &iq)

	fmt.Println(iq.Bind.Resource)

	clientInfo.SocketID = clientInfo.User.AccountId + "@prod.ol.epicgames.com/" + iq.Bind.Resource

	
	conn.WriteMessage(messageType, []byte(`
		<iq to="` + clientInfo.SocketID + `" id="_xmpp_bind1" xmlns="jabber:client" type="result">
			<bind xmlns="urn:ietf:params:xml:ns:xmpp-bind">
				<jid>` + clientInfo.SocketID + `</jid>
			</bind>
		</iq>
	`))
}

func XHandleOpen(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	conn.WriteMessage(messageType, []byte(`
		<open xmlns="urn:ietf:params:xml:ns:xmpp-framing" from="prod.ol.epicgames.com" version="1.0" id="` + clientInfo.UUID  + `" />
	`))

	dataNeedAuth := models.StreamFeaturesNeedAuth{
		Stream: "http://etherx.jabber.org/streams",
		Mechanisms: models.Mechanisms{
			Mechanism: "PLAIN",
		},
		Ver: models.Ver{},
		StartTLS: models.StartTLS{},
		Compression: models.Compression{
			Method: "zlib",
		},
		Auth: models.Auth{},
	}

	dataAuth := models.StreamFeatures{
		Stream: "http://etherx.jabber.org/streams",
		Ver: models.Ver{},
		StartTLS: models.StartTLS{},
		Compression: models.Compression{
			Method: "zlib",
		},
		Session: models.Session{},
		Bind: models.Bind{},
	}

	var xmlDataToMarshal []byte

	if clientInfo.Authenticated { xmlDataToMarshal, _ = xml.Marshal(dataAuth) } else  { xmlDataToMarshal, _ = xml.Marshal(dataNeedAuth) }

	conn.WriteMessage(messageType, xmlDataToMarshal)
}

func XHandleAuth(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	var auth models.AuthXML
	xml.Unmarshal([]byte(message), &auth)

	decoded, err := base64.StdEncoding.DecodeString(auth.Value)
	if err != nil {
		conn.Close()
	}
	authData := strings.Split(string(decoded), "eg1~")

	user, err := middleware.VerifyAccessTokenXMPP(authData[1])
	if err != nil {
		all.PrintRed([]any{authData[0], authData[1]})
		conn.Close()
		return
	}

	clientInfo.Authenticated = true
	clientInfo.User = user

	all.PrintGreen([]any{"user logged in:", user.Username})
	
	conn.WriteMessage(
		messageType, 
		[]byte(`<success xmlns="urn:ietf:params:xml:ns:xmpp-sasl" />`),
	)
}

func XHandlePresence(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	all.PrintRed([]any{"normal presence"})
	var presence models.PresenceXML
	xml.Unmarshal([]byte(message), &presence)

	var status models.StatusJSON
	json.Unmarshal([]byte(presence.Status.Value), &status)

	// if status.Status == "" {
	// }
	XGetFriendStatus(clientInfo)

	clientInfo.Status = presence.Status.Value

	friends := common.GetAllAcceptedFriends(clientInfo.User.AccountId)
	for _, friend := range friends {
		friendClient, err := XGetClientFromAccountId(friend.AccountId)
		if err != nil {
			continue
		}

		friendClient.Connection.WriteMessage(1, []byte(`
			<presence to="`+ friendClient.SocketID +`" xmlns="jabber:client" from="`+ clientInfo.SocketID +`" type="available">
				<status>`+ presence.Status.Value +`</status>
			</presence>
		`))
	}

	XMPPUpdateStatus(clientInfo.User.AccountId, clientInfo.User.AccountId)
}

func XHandlePartyPresence(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	var presence models.PartyPresenceXML
	xml.Unmarshal([]byte(message), &presence)
	
	all.PrintMagenta([]any{"party presence", presence})
	// Party-2ad8b220-4f13-4456-a0a3-cbde2bcbfcfd@muc.prod.ol.epicgames.com/admin:571f16e7-c6aa-41f5-b24c-edc70fc88406:V2:Fortnite:WIN::E0EB415645D78EC5C252798418B1548A
	// to := strings.Split(presence.To, "/")

	// clientInfo.Connection.WriteMessage(1, []byte(`
	// 	<presence to="`+ clientInfo.SocketID +`" from="`+ presence.To +`" xmlns="jabber:client" type="unavailable">
	// 		<x xmlns="http://jabber.org/protocol/muc#user">
	// 			<item nick="`+ clientInfo.User.Username +`" jid="`+ clientInfo.SocketID +`" role="participant"/>
	// 			<status code="110"/>
	// 			<status code="100"/>
	// 			<status code="170"/>
	// 		</x>
	// 	</presence>
	// `))

	foundPartyId, ok := common.AccountIdToPartyId[clientInfo.User.AccountId]
	if !ok {
		return
	}

	// party, ok := common.ActiveParties[foundPartyId]
	// if !ok {
	// 	return
	// }

	if presence.Type == "unavailable" {
		clientInfo.Connection.WriteMessage(1, []byte(`
			<presence to="`+ clientInfo.SocketID +`" from="`+ presence.To +`" xmlns="jabber:client" type="unavailable">
				<x xmlns="http://jabber.org/protocol/muc#user">
					<item nick="`+ partyNick(clientInfo.User, clientInfo.SocketID) +`" jid="`+ partySocketId(foundPartyId, partyNick(clientInfo.User, clientInfo.SocketID)) +`" role="none"/>
					<status code="110"/>
					<status code="100"/>
					<status code="170"/>
				</x>
			</presence>
		`))
		return
	}

		clientInfo.Connection.WriteMessage(1, []byte(`
		<presence to="`+ clientInfo.SocketID +`" from="`+ presence.To +`" xmlns="jabber:client">
			<x xmlns="http://jabber.org/protocol/muc#user">
				<item nick="`+ partyNick(clientInfo.User, clientInfo.SocketID) +`" jid="`+ partySocketId(foundPartyId, partyNick(clientInfo.User, clientInfo.SocketID)) +`" role="none"/>
				<status code="110"/>
				<status code="100"/>
				<status code="170"/>
			</x>
		</presence>
	`))

	party, ok := common.ActiveParties[foundPartyId]
	if !ok {
		return
	}

	for _, member := range party.Members {
		partyMemberClient, err := XGetClientFromAccountId(member.AccountId)
		if err != nil {
			continue
		}
				
		clientInfo.Connection.WriteMessage(1, []byte(`
			<presence to="`+ clientInfo.SocketID +`" from="`+ partyMemberClient.SocketID +`" xmlns="jabber:client">
				<x xmlns="http://jabber.org/protocol/muc#user">
					<item nick="`+ partyNick(partyMemberClient.User, partyMemberClient.SocketID) +`" jid="`+ partyMemberClient.SocketID +`" role="participant" affiliation="none"/>
				</x>
			</presence>
		`))

		all.PrintMagenta([]any{`
			<presence to="`+ clientInfo.SocketID +`" from="`+ clientInfo.SocketID +`" xmlns="jabber:client">
				<x xmlns="http://jabber.org/protocol/muc#user">
					<item nick="`+ partyNick(partyMemberClient.User, partyMemberClient.SocketID) +`" jid="`+ partyMemberClient.SocketID +`" role="participant" affiliation="none"/>
				</x>
			</presence>
		`})

		if partyMemberClient.User.AccountId == clientInfo.User.AccountId {
			continue
		}

		partyMemberClient.Connection.WriteMessage(1, []byte(`
			<presence to="`+ partyMemberClient.SocketID +`" from="`+ clientInfo.SocketID +`" xmlns="jabber:client">
				<x xmlns="http://jabber.org/protocol/muc#user">
					<item nick="`+ partyNick(clientInfo.User, clientInfo.SocketID) +`" jid="`+ clientInfo.SocketID +`" role="participant" affiliation="none"/>
				</x>
			</presence>
		`))

		all.PrintMagenta([]any{`
			<presence to="`+ partyMemberClient.SocketID +`" from="`+ partySocketId(foundPartyId, partyNick(clientInfo.User, clientInfo.SocketID)) +`" xmlns="jabber:client">
				<x xmlns="http://jabber.org/protocol/muc#user">
					<item nick="`+ partyNick(clientInfo.User, clientInfo.SocketID) +`" jid="`+ clientInfo.SocketID +`" role="participant" affiliation="none"/>
				</x>
			</presence>
		`})
	}
}


func partyNick(user models.User, socketId string) string {
	return user.Username + ":" + user.AccountId + ":" + strings.Split(socketId, "/")[1]
}

func partySocketId(partyId string, socketId string) string {
	return "Party-" + partyId + "@muc.prod.ol.epicgames.com/" + socketId
}


func XGetFriendStatus(clientInfo *ClientInfo) {
	friends := common.GetAllAcceptedFriends(clientInfo.User.AccountId)
	for _, friend := range friends {
		friendClient, err := XGetClientFromAccountId(friend.AccountId)
		if err != nil {
			continue
		}

		// XMPPUpdateStatusSingle(friendClient.User.AccountId, clientInfo.User.AccountId)

		clientInfo.Connection.WriteMessage(1, []byte(`
			<presence to="`+ clientInfo.SocketID +`" xmlns="jabber:client" from="`+ friendClient.SocketID +`" type="available">
				<status>`+ friendClient.Status +`</status>
			</presence>
		`))
	}
}

func XGetClientFromAccountId(accountId string) (*ClientInfo, error) {
	clientRemoteAddress, ok := AccountIdToXMPPRemoteAddress[accountId]
	if !ok {
		return nil, fmt.Errorf("failed to find client remote address")
	}

	client, ok := ActiveXMPPClients[clientRemoteAddress]
	if !ok {
		return nil, fmt.Errorf("failed to find client")
	}

	return client, nil
}

func XMPPSendBodyToAll(body map[string]interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		return
	}

	for _, client := range ActiveXMPPClients {
		client.Connection.WriteMessage(1, []byte(`
			<message xmlns="jabber:client" from="xmpp-admin@prod.ol.epicgames.com" to="`+ client.SocketID +`">
				<body>`+ string(data) +`</body>
			</message>
		`))
	}
}

func XMPPSendBody(body map[string]interface{}, accountId string) {
	data, err := json.Marshal(body)
	if err != nil {
		return
	}

	client, err := XGetClientFromAccountId(accountId)
	if err != nil {
		return
	}

	client.Connection.WriteMessage(1, []byte(`
		<message xmlns="jabber:client" from="xmpp-admin@prod.ol.epicgames.com" to="`+ client.SocketID +`">
			<body>`+ string(data) +`</body>
		</message>
	`))
}

func XMPPUpdateStatus(accountId string, friendId string) {
	mainClient, err := XGetClientFromAccountId(accountId)
	if err != nil {
		return
	}

	friendClient, err := XGetClientFromAccountId(friendId)
	if err != nil {
		return
	}

	friendClient.Connection.WriteMessage(1, []byte(`
		<presence to="`+ friendClient.SocketID +`" xmlns="jabber:client" from="`+ mainClient.SocketID +`" type="available">
			<status>`+ mainClient.Status +`</status>
		</presence>
	`))

	mainClient.Connection.WriteMessage(1, []byte(`
		<presence to="`+ mainClient.SocketID +`" xmlns="jabber:client" from="`+ friendClient.SocketID +`" type="available">
			<status>`+ friendClient.Status +`</status>
		</presence>
	`))
}

func SendJoinPartyRequest(accountId string, partyId string, ac string) {
	client, err := XGetClientFromAccountId(accountId)
	if err != nil {
		return
	}

	joinPresence := gin.H{
		"Status": "",
		"bIsJoinable": false,
		"bIsPlaying": false,
		"bHasVoiceSupport": false,
		"SessionId": "",
		"Properties": gin.H{
			"party.joininfodata.286331153_j": gin.H{
				"sourceId": client.User.AccountId,
				"sourceDisplayName": client.User.Username,
				"sourcePlatform": "WIN",
				"partyId": partyId,
				"partyTypeId": 286331153,
				"key": ac,
				"appId": "Fortnite",
				"buildId": "1:1:",
				"partyFlags": 6,
				"notAcceptingReason": 0,
			},
		},
	}

	jsonPresence, err := json.Marshal(joinPresence)
	if err != nil {
		return
	}

	friends := common.GetAllAcceptedFriends(accountId)
	for _, friend := range friends {
		friendClient, err := XGetClientFromAccountId(friend.AccountId)
		if err != nil {
			continue
		}

		friendClient.Connection.WriteMessage(1, []byte(`
			<presence to="`+ friendClient.SocketID +`" xmlns="jabber:client" from="`+ client.SocketID +`" type="available">
				<status>`+ string(jsonPresence) +`</status>
			</presence>
		`))

		all.PrintMagenta([]any{`
			<presence to="`+ client.SocketID +`" from="`+ friendClient.SocketID +`" xmlns="jabber:client">
				<status>`+ string(jsonPresence) +`</status>
			</presence>
		`})
	}
}