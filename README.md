![template(7)(1)](https://github.com/zombman/backend/assets/121964555/a1624e90-ed7a-4a8d-a00c-b6707c69b4ac)

An open source, performant fortnite backend server with a built in web-interface and desktop launcher! Thank you to @mhtsotakis for help with party v2! _Note: Some features are not yet implemented, please see the [Roadmap](#roadmap) for more information._

## Branches

- **[Frontend](https://github.com/zombman/server/tree/frontend)** - Launcher and web interface for backend.
- **[CLI](https://github.com/zombman/server/tree/cli)** - The CLI that the launcher interacts with to launch the game.

## Features

- **Blazing fast:** Written in Go, this server is extremely fast and lightweight.
- **Easy to use:** Designed to be easy to use and setup.
- **Open source:** Completely free to use, share and modify!
- **Web interface:** Built in web interface to manage your server.

## Setup

- Watch the quick setup guide [here](https://www.youtube.com/watch?v=WvWrgmEH6ZI&t=189s&ab_channel=ulnk).
- To keep your backend updated with the latest updates, instead of downloading the project, use the command `git clone https://github.com/zombman/backend` to download. You may need to install [git](https://git-scm.com/) for this. Now whenever I update the repo, just use the command `git pull`!

## Roadmap

Roadmap is in order of completion, so closer to the top of list means it will be done before items lower down the list. View the trello board [here](https://trello.com/b/7AKhxa5T/zombie-server) to see a more detailed view of this roadmap and to view any bugs.

- [x] Basic user creation and authentication
- [x] Access the fortnite lobby
- [x] Profiles per user (not shared/full locker)
- [x] Rotating item shop with correct prices and items from season or below
- [x] Gifting & Buying from item shop
- [x] Equip items and variants (works with superhero skins)
- [x] Cloud storage for settings
- [x] Control panel for server admins
- [x] Friends (friends v2)
- [x] Matchmaker & XMPP
- [x] Parties (party v2 is kinda scuffed)
- [x] Daily rewards
- [x] Change lobby background in chapter 2+
- [ ] Battle Pass & Levelling Up
- [ ] Challenges (or Quests)

## Compatibility

- Everything is functional in **Chapter 1**
- Unkown issue causing matchmaker and parties to not work on **Season 16+**

## Prerequisites

- A basic understanding of how to build and run applications. Now made easier with the [scripts folder!](https://github.com/zombman/backend/tree/master/_setup)
- [GoLang](https://go.dev)
