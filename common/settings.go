package common

import (
	"math/rand"
	"sort"
)

type GameServer struct {
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Playlist    string `json:"platform"`
	Region      string `json:"region"`
	PlayersLeft int    `json:"playersLeft"`
	Joinable    bool   `json:"joinable"`
}

var (
	GameServers             = make(map[string][]GameServer)
	IP               string = "127.0.0.1:3000"
	Season           int    = 7
	Chapter          int    = 2
	LoadShopFromJson bool   = false
)

func InitGameServers() {
	// addGameServer("playlist_defaultsolo", "EU", "158.178.203.104", 7777)
	// addGameServer("playlist_defaultsolo", "NAE", "158.178.203.104", 7777)
	// addGameServer("playlist_defaultsolo", "NAW", "158.178.203.104", 7777)
	addGameServer("playlist_defaultsolo", "EU", "127.0.0.1", 7777)
	addGameServer("playlist_defaultsolo", "NAE", "127.0.0.1", 7777)
	addGameServer("playlist_defaultsolo", "NAW", "127.0.0.1", 7777)
}

func GetGameServer(playlist string, region string) GameServer {
	for _, server := range GameServers[playlist+":"+region] {
		if server.Joinable {
			return server
		}
	}
	return GameServer{
		Joinable: false,
	}
}

func GetAllGameServers(playlist string, region string) []GameServer {
	return GameServers[playlist+":"+region]
}

func SortGameServersByPlayersLeft(playlist string, region string) []GameServer {
	gameServerClone := GetAllGameServers(playlist, region)
	sort.SliceStable(gameServerClone, func(i, j int) bool {
		return gameServerClone[i].PlayersLeft < gameServerClone[j].PlayersLeft
	})

	return gameServerClone
}

func addGameServer(playlist string, region string, ip string, port int) {
	GameServers[playlist+":"+region] = append(GameServers[playlist+":"+region], GameServer{
		IP:          ip,
		Port:        port,
		Playlist:    playlist,
		Region:      region,
		PlayersLeft: rand.Intn(100 - 43) + 43,
		Joinable:    true,
	})
}