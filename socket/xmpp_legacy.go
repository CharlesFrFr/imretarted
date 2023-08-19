package socket

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

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
	JID string
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
		JID: "",
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

			XHandlePresence(conn, messageData, messageType, clientInfo)
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
			<message to="`+ recipientClient.JID +`" from="`+ clientInfo.JID +`" id="`+ msg.ID +`" xmlns="jabber:client" type="chat">
				<body>`+ msg.Body.Value +`</body>
			</message>
		`))
		return
	}

	recipientClient.Connection.WriteMessage(messageType, []byte(`
		<message to="`+ recipientClient.JID +`" from="`+ clientInfo.JID +`" id="`+ msg.ID +`" xmlns="jabber:client">
			<body>`+ msg.Body.Value +`</body>
		</message>
	`))
}

func XHandleSessionIQ(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo, id string) {
	conn.WriteMessage(messageType, []byte(`
		<iq to="`+ clientInfo.JID +`" from="prod.ol.epicgames.com" id="`+ id +`" xmlns="jabber:client" type="result" />
	`))
}

func XHandleBindIQ(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	var iq models.BindIQXML
	xml.Unmarshal([]byte(message), &iq)

	fmt.Println(iq.Bind.Resource)

	clientInfo.JID = clientInfo.User.AccountId + "@prod.ol.epicgames.com/" + iq.Bind.Resource
	clientInfo.Resource = iq.Bind.Resource
	
	conn.WriteMessage(messageType, []byte(`
		<iq to="` + clientInfo.JID + `" id="_xmpp_bind1" xmlns="jabber:client" type="result">
			<bind xmlns="urn:ietf:params:xml:ns:xmpp-bind">
				<jid>` + clientInfo.JID + `</jid>
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
			<presence to="`+ friendClient.JID +`" xmlns="jabber:client" from="`+ clientInfo.JID +`" type="available">
				<status>`+ presence.Status.Value +`</status>
			</presence>
		`))
	}

	XMPPUpdateStatus(clientInfo.User.AccountId, clientInfo.User.AccountId)
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
			<presence to="`+ clientInfo.JID +`" xmlns="jabber:client" from="`+ friendClient.JID +`" type="available">
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
			<message xmlns="jabber:client" from="xmpp-admin@prod.ol.epicgames.com" to="`+ client.JID +`">
				<body>`+ string(data) +`</body>
			</message>
		`))
	}
}

func XMPPSendBody(body map[string]interface{}, client *ClientInfo) {
	data, err := json.Marshal(body)
	if err != nil {
		return
	}

	client.Connection.WriteMessage(1, []byte(`
		<message xmlns="jabber:client" from="xmpp-admin@prod.ol.epicgames.com" to="`+ client.JID +`">
			<body>`+ string(data) +`</body>
		</message>
	`))
}

func XMPPSendBodyToAccountId(body map[string]interface{}, accountId string) {
	data, err := json.Marshal(body)
	if err != nil {
		return
	}

	client, err := XGetClientFromAccountId(accountId)
	if err != nil {
		return
	}

	client.Connection.WriteMessage(1, []byte(`
		<message xmlns="jabber:client" from="xmpp-admin@prod.ol.epicgames.com" to="`+ client.JID +`">
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
		<presence to="`+ friendClient.JID +`" xmlns="jabber:client" from="`+ mainClient.JID +`" type="available">
			<status>`+ mainClient.Status +`</status>
		</presence>
	`))

	mainClient.Connection.WriteMessage(1, []byte(`
		<presence to="`+ mainClient.JID +`" xmlns="jabber:client" from="`+ friendClient.JID +`" type="available">
			<status>`+ friendClient.Status +`</status>
		</presence>
	`))
}