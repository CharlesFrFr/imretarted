# Zombie Backend

An open source, performant fortnite backend server with a built in web-interface and desktop launcher!

_Note: Some features are not yet implemented, please see the [Roadmap](#roadmap) for more information._

# Branches

- **[Frontend](https://github.com/zombman/server/tree/frontend)** - Launcher and web interface for backend.
- **[CLI](https://github.com/zombman/server/tree/cli)** - The CLI that the launcher interacts with to launch the game.

## Features

- **Blazing fast:** Written in Go, this server is extremely fast and lightweight.
- **Easy to use:** Designed to be easy to use and setup.
- **Open source:** Completely free to use, share and modify!
- **Web interface:** Built in web interface to manage your server. (not implemented yet)

## Roadmap

- [x] Basic user creation and authentication
- [x] Fortnite uses oauth tokens
- [x] Profiles to access lobby and item shop
- [x] Random item shop
- [x] Buy from item shop
- [x] Equip items and variants
- [x] Cloud storage for settings
- [x] Control panel for server admins
- [x] Game launcher and web interface
- [x] Friends
- [x] Matchmaker
  - [ ] Not smart atm. Currently making a custom game server so that the backend can talk to each game server and retrieve the required information to make a smart matchmaker. (e.g. only let first 100 in the queue enter match, properly calculated eta, etc.)
- [x] Xmpp
- [x] Parties
- [ ] Party v2
