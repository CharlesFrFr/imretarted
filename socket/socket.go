package socket

import (
	"encoding/base64"
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zombman/server/all"
	"github.com/zombman/server/middleware"
	"github.com/zombman/server/models"
)

var (
	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	clients = make(map[string]*ClientInfo)
)

type ClientInfo struct {
	UUID string
	Authenticated bool
	User models.User
	SocketID string
}

func Handler(w http.ResponseWriter, r *http.Request){
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		all.PrintRed([]any{ "Failed to set websocket upgrade: ", err, })
		return
	}

	remoteAddress := conn.RemoteAddr().String()
	all.PrintMagenta([]any{ "new client connected:", remoteAddress, })

	clientInfo := &ClientInfo{
		UUID: uuid.New().String(),
		Authenticated: false,
		User: models.User{},
		SocketID: "",
	}
	clients[remoteAddress] = clientInfo

	for {
		messageType, messageData, err := conn.ReadMessage()
		if err != nil {
			break
		}

		all.PrintCyan([]any{string(messageData)})

		if xml.Unmarshal([]byte(messageData), &models.OpenXML{}) == nil {
			HandleOpen(conn, messageData, messageType, clientInfo)
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.AuthXML{}) == nil {
			HandleAuth(conn, messageData, messageType, clientInfo)
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.BindIQXML{}) == nil {
			var iq models.BindIQXML
			xml.Unmarshal([]byte(messageData), &iq)

			if iq.ID == "_xmpp_bind1" {
				HandleBindIQ(conn, messageData, messageType, clientInfo)
				continue
			}
			HandleSessionIQ(conn, messageData, messageType, clientInfo)
			
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.MessageWithBodyXML{}) == nil {
			HandleMessage(conn, messageData, messageType, clientInfo)
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.PresenceXML{}) == nil {
			HandlePresence(conn, messageData, messageType, clientInfo)
			continue
		}

		if xml.Unmarshal([]byte(messageData), &models.CloseXML{}) == nil {
			all.PrintYellow([]any{"HandleCloseMessage"})
			conn.Close()
			continue
		}


		all.PrintRed([]any{"unknown message type"})

		conn.WriteMessage(messageType, []byte(`CRASH`))
	}

	defer func() {
		delete(clients, remoteAddress)
		conn.Close()
	}()
}

func HandlePresence(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	all.PrintYellow([]any{"HandlePresence"})

	var presence models.PresenceXML
	xml.Unmarshal([]byte(message), &presence)

	
}


func HandleMessage(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	all.PrintYellow([]any{"HandleMessage"})

	var msg models.MessageWithBodyXML
	xml.Unmarshal([]byte(message), &msg)

	conn.WriteMessage(messageType, []byte(`
		<message from="`+ clientInfo.SocketID +`" id="`+ msg.ID +`" to="`+ strings.Split(msg.To, "/")[0] +`" xmlns="jabber:client">
			<body>`+ msg.Body.Value +`</body>
		</message>
	`))
}

func HandleSessionIQ(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	all.PrintYellow([]any{"HandleSessionIQ"})
		conn.WriteMessage(messageType, []byte(`
		<iq to="`+ clientInfo.SocketID +`" from="prod.ol.epicgames.com" id="_xmpp_session1" xmlns="jabber:client" type="result" />
	`))
}

func HandleBindIQ(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	all.PrintYellow([]any{"HandleBindIQ"})

	var iq models.BindIQXML
	xml.Unmarshal([]byte(message), &iq)

	clientInfo.SocketID = clientInfo.User.AccountId + "@prod.ol.epicgames.com/" + iq.Bind.Resource
	all.PrintMagenta([]any{"user socket id:", clientInfo.SocketID})
	
	conn.WriteMessage(messageType, []byte(`
		<iq to="` + clientInfo.SocketID + `" id="_xmpp_bind1" xmlns="jabber:client" type="result">
			<bind xmlns="urn:ietf:params:xml:ns:xmpp-bind">
				<jid>` + clientInfo.SocketID + `</jid>
			</bind>
		</iq>
	`))
}

func HandleOpen(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	all.PrintYellow([]any{"HandleOpenMessage"})

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
	all.PrintRed([]any{"sending auth data", clientInfo.Authenticated})

	conn.WriteMessage(messageType, xmlDataToMarshal)
}

func HandleAuth(conn *websocket.Conn, message []byte, messageType int, clientInfo *ClientInfo) {
	all.PrintYellow([]any{"HandleAuth"})

	var auth models.AuthXML
	xml.Unmarshal([]byte(message), &auth)

	decoded, err := base64.StdEncoding.DecodeString(auth.Value)
	if err != nil {
		all.PrintRed([]any{"failed to decode base64", err})
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

	// auth.Value
	
	conn.WriteMessage(
		messageType, 
		[]byte(`<success xmlns="urn:ietf:params:xml:ns:xmpp-sasl" />`),
	)
}	