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
var LoadShopFromJson bool = true

var GameServers = make(map[string]GameServer)

func InitGameServers() {
	addGameServer("playlist_defaultsolo", "EU", "158.178.203.104", 7777)
}

func addGameServer(playlist string, region string, ip string, port int) {
	GameServers[playlist+":"+region] = GameServer{
		IP:          ip,
		Port:        port,
		Playlist:    playlist,
		Region:      region,
		PlayersLeft: 0,
		Joinable:    true,
	}
}