package controllers

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
	"github.com/zombman/server/socket"
)

var (
	ActivePings []models.Ping
)

func PostPing(c *gin.Context) {
	sentBy, _ := common.GetUserByAccountId(c.Param("pingerId"))
	sentTo, err := common.GetUserByAccountId(c.Param("accountId"))

	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	var body gin.H
	c.BindJSON(&body)

	ping := models.Ping{
		SentBy: sentBy.AccountId,
		SentTo: sentTo.AccountId,
		SentAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		ExpiresAt: time.Now().Add(time.Minute * 5).Format("2006-01-02T15:04:05.999Z"),
		Meta: body,
	}

	ActivePings = append(ActivePings, ping)

	socket.XMPPSendBodyToAccountId(gin.H{
		"pinger_id": sentBy.AccountId,
		"pinger_dn": sentBy.Username,
		"sent": ping.SentAt,
		"expires": ping.ExpiresAt,
		"meta": ping.Meta,
		"ns": "Fortnite",
		"type": "com.epicgames.social.party.notification.v0.PING",
	}, sentTo.AccountId)

	c.JSON(200, ping)
}

func PostPartyPing(c *gin.Context) {
	sentBy, _ := common.GetUserByAccountId(c.Param("pingerId"))
	sentTo, err := common.GetUserByAccountId(c.Param("accountId"))
	party, ok := common.ActiveParties[c.Param("partyId")]

	if err != nil || !ok {
		common.ErrorBadRequest(c)
		return
	}

	var body gin.H
	c.BindJSON(&body)

	ping := models.Ping{
		SentBy: sentBy.AccountId,
		SentTo: sentTo.AccountId,
		SentAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		ExpiresAt: time.Now().Add(time.Minute * 5).Format("2006-01-02T15:04:05.999Z"),
		Meta: body,
	}
	ActivePings = append(ActivePings, ping)

	socket.XMPPSendBodyToAccountId(gin.H{
		"party_id": party.ID,
		"inviter_id": sentBy.AccountId,
		"inviter_dn": sentBy.Username,
		"sent": ping.SentAt,
		"expires": ping.ExpiresAt,
		"meta": ping.Meta,
		"ns": "Fortnite",
		"type": "com.epicgames.social.party.notification.v0.INITIAL_INVITE",
	}, sentTo.AccountId)

	c.JSON(200, ping)
}

func DeletePing(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	pingId := c.Param("pingId")

	for i, ping := range ActivePings {
		if ping.SentBy == user.AccountId && pingId == ping.SentTo {
			ActivePings = append(ActivePings[:i], ActivePings[i+1:]...)
			c.Status(204)
			return
		}
	}

	common.ErrorBadRequest(c)
}

func GetPings(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	var pings []models.Ping

	for _, ping := range ActivePings {
		if ping.SentTo == user.AccountId  {
			pings = append(pings, ping)
		}
	}

	c.JSON(200, pings)
}

func GetPartyPings(c *gin.Context) {
	sentBy, _ := common.GetUserByAccountId(c.Param("pingerId"))

	var parties []models.V2Party
	for _, party := range common.ActiveParties {
		for _, member := range party.Members {
			if member.AccountId == sentBy.AccountId {
				parties = append(parties, party)
			}
		}
	}

	all.MarshPrintJSON(parties)

	c.JSON(200, parties)
}

func JoinPing(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	sentBy, err := common.GetUserByAccountId(c.Param("pingerId"))

	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	var body struct {
		Connection struct {
			Id string `json:"id"`
			Meta struct {
				UrnEpicConnPlatformS string `json:"urn:epic:conn:platform_s"`
			} `json:"meta"`
		} `json:"connection"`
		Meta struct {
			UrnEpicMemberDnS string `json:"urn:epic:member:dn_s"`
			UrnJoinRequestUsers string `json:"urn:epic:member:joinrequestusers_j"`
		} `json:"meta"`
	}

	if err := c.BindJSON(&body); err != nil {
		common.ErrorBadRequest(c)
		return
	}

	var pings []models.Ping
	for _, ping := range ActivePings {
		if ping.SentBy == sentBy.AccountId && ping.SentTo == user.AccountId {
			pings = append(pings, ping)
		}
	}

	if len(pings) == 0 {
		common.ErrorBadRequest(c)
		return
	}

	ActivePings = append(ActivePings[:0], ActivePings[1:]...)

	var party models.V2Party
	for _, p := range common.ActiveParties {
		for _, member := range p.Members {
			if member.AccountId == sentBy.AccountId {
				party = p
				break
			}
		}
	}
	
	if party.ID == "" {
		common.ErrorBadRequest(c)
		return
	}
	
	var captain models.V2PartyMember
	for _, member := range party.Members {
		if member.Role == "CAPTAIN" {
			captain = member
			break
		}
	}

	connectionMeta := make(map[string]interface{})
	connectionMeta["urn:epic:conn:platform_s"] = body.Connection.Meta.UrnEpicConnPlatformS
	connection := models.V2PartyConnection{
		ID: body.Connection.Id,
		Meta: connectionMeta,
		YieldLeadership: false,
		ConnectedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
	}

	partyMemberMeta := make(map[string]interface{})
	partyMemberMeta["urn:epic:member:dn_s"] = body.Meta.UrnEpicMemberDnS
	partyMemberMeta["urn:epic:member:joinrequestusers_j"] = body.Meta.UrnJoinRequestUsers

	partyMember := models.V2PartyMember{
		AccountId: user.AccountId,
		Meta: partyMemberMeta,
		Connections: []models.V2PartyConnection{connection},
		Role: "MEMBER",
		Revision: 0,
		JoinedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
	}

	party.Members = append(party.Members, partyMember)
	
	captianJoinRequests := models.V2CaptainJoinRequestUsers{
		Users: []models.V2CaptainJoinRequestUser{},
	}
	rawSquadAssignments := models.V2RawSquadAssignments{
		RawSquadAssignments: []models.V2RawSquadAssignment{},
	}

	for i, member := range party.Members {
		captianJoinRequests.Users = append(captianJoinRequests.Users, models.V2CaptainJoinRequestUser{
			ID: member.AccountId,
			DisplayName: member.Meta["urn:epic:member:dn_s"].(string),
			Platform: "WIN",
			Data: "{\"CrossplayPreference_i\":\"1\"}",
		})
		rawSquadAssignments.RawSquadAssignments = append(rawSquadAssignments.RawSquadAssignments, models.V2RawSquadAssignment{
			MemberId: member.AccountId,
			AbsoluteMemberIdx: i,
		})
	}

	for i, member := range party.Members {
		if member.Role == "CAPTAIN" {
			party.Members = append(party.Members[:i], party.Members[i+1:]...)
			break
		}
	}

	captianJoinRequestsRaw, _ := json.Marshal(captianJoinRequests)
	captain.Meta["urn:epic:member:joinrequestusers_j"] = string(captianJoinRequestsRaw)

	rawSquadAssignmentsRaw, _ := json.Marshal(rawSquadAssignments)
	party.Meta["RawSquadAssignments_j"] = string(rawSquadAssignmentsRaw)

	party.Members = append(party.Members, captain)
	common.ActiveParties[party.ID] = party
	common.AccountIdToPartyId[user.AccountId] = party.ID

	for _, member := range party.Members {
		memberClient, err := socket.XGetClientFromAccountId(member.AccountId)
		if err != nil {
			continue
		}

		socket.XMPPSendBody(gin.H{
			"account_id": partyMember.AccountId,
			"account_dn": partyMember.Meta["urn:epic:member:dn_s"],
			"member_state_updated": partyMember.Meta,
			"party_id": party.ID,
			"updated_at": partyMember.UpdatedAt,
			"joined_at": partyMember.JoinedAt,
			"sent": time.Now().Format("2006-01-02T15:04:05.000Z"),
			"revision": party.Revision,
			"ns": "Fortnite",
			"type": "com.epicgames.social.party.notification.v0.MEMBER_JOINED",
		}, memberClient)
	}
	
	c.JSON(201, gin.H{
		"status": "JOINED",
		"party_id": party.ID,
	})

	deleteAnyEmptyParties()
}