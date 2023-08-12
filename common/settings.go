package common

type GameServer struct {
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Playlist    string `json:"platform"`
	Region      string `json:"region"`
	PlayersLeft int    `json:"playersLeft"`
	Joinable    bool   `json:"busLeft"`
}

var IP string = "127.0.0.1:3000"
var Season int = 8
var Chapter int = 1
var LoadShopFromJson bool = false

var GameServers = make(map[string]GameServer)

func InitGameServers() {
	GameServers["playlist_defaultsolo:EU"] = GameServer{
		IP:          "127.0.0.1",
		Port:        7777,
		Playlist:    "playlist_defaultsolo",
		Region:      "EU",
		PlayersLeft: 0,
		Joinable:    true,
	}
}
