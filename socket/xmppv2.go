package socket

import (
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)
var (
	XMPPClients = make(map[string]*ClientInfo)
	AccountIdToXMPPClientKey = make(map[string]string)
)

func InitXMPPServer() {
	router := xmpp.NewRouter()
	router.HandleFunc("auth", handleAuth)
}

func handleAuth(s xmpp.Sender, p stanza.Packet) {
	
}

//  HI ID DO NOT WANT TO MAKE ANOTHER 1 BUT THER IS A PREBUILT LIB FOR THIS
// SO IF I NEED TO REWORK I WILL