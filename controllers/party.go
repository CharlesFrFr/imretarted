package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

var (
	ActiveParties = make(map[string]models.V2Party)
	AccountIdToPartyId = make(map[string]string)
)

func PartyGetUser(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	all.PrintMagenta([]any{"PartyGetUser"})
	
	partyId := AccountIdToPartyId[user.AccountId]
	party, ok := ActiveParties[partyId]
	if !ok {
		party = common.CreateParty(&ActiveParties, &AccountIdToPartyId, user)
	}
	
	c.JSON(200, gin.H{
		"current": party,
	})
	all.MarshPrintJSON(party)
}

func PartyGetFriendPartyPings(c *gin.Context) {
	friend, err := common.GetUserByAccountId(c.Param("friendId"))
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	partyId, ok := AccountIdToPartyId[friend.AccountId]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	party, ok := ActiveParties[partyId]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	c.JSON(200, []models.V2Party{party})
}

func PartyPost(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	all.PrintMagenta([]any{"PartyPost for user", user.Username})

	var body struct {
		Config struct {
			JoinConfirmation bool `json:"join_confirmation"`
			Joinability string `json:"joinability"`
			MaxSize int `json:"max_size"`
		} `json:"config"`
		JoinInfo struct {
			Connection struct {
				Id string `json:"id"`
				Meta struct {
					UrnEpicConnPlatformS string `json:"urn:epic:conn:platform_s"`
					UrnEpicConnTypeS string `json:"urn:epic:conn:type_s"`
				} `json:"meta"`
			} `json:"connection"`
			Meta struct {
				UrnEpicMemberDnS string `json:"urn:epic:member:dn_s"`
				UrnEpicMemberTypeS string `json:"urn:epic:member:type_s"`
				UrnEpicMemberPlatformS string `json:"urn:epic:member:platform_s"`
			} `json:"meta"`
		} `json:"join_info"`
		Meta map[string]interface{} `json:"meta"`
	}

	if err := c.BindJSON(&body); err != nil {
		common.ErrorBadRequest(c)
		return
	}

	all.MarshPrintJSON(body)
	party := common.CreateParty(&ActiveParties, &AccountIdToPartyId, user)
	ActiveParties[party.ID] = party
	AccountIdToPartyId[user.AccountId] = party.ID

	connectionMeta := make(map[string]interface{})
	connectionMeta["urn:epic:conn:platform_s"] = body.JoinInfo.Connection.Meta.UrnEpicConnPlatformS
	connectionMeta["urn:epic:conn:type_s"] = body.JoinInfo.Connection.Meta.UrnEpicConnTypeS
	connection := models.V2PartyConnection{
		ID: body.JoinInfo.Connection.Id,
		Meta: connectionMeta,
		YieldLeadership: false,
		ConnectedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
	}

	partyMemberMeta := make(map[string]interface{})
	partyMemberMeta["urn:epic:member:dn_s"] = user.Username
	partyMemberMeta["urn:epic:member:joinrequestusers_j"] = "{\"users\":[{\"id\":\""+ user.AccountId +"\",\"dn\":\""+ user.Username +"\",\"plat\":\"WIN\",\"data\":\"{\\\"CrossplayPreference_i\\\":\\\"1\\\",\\\"SubGame_u\\\":\\\"1\\\"}\"}]}"
	partyMember := models.V2PartyMember{
		AccountId: user.AccountId,
		Meta: partyMemberMeta,
		Connections: []models.V2PartyConnection{connection},
		Role: "CAPTAIN",
		Revision: 0,
		JoinedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
	}

	party.Config.JoinConfirmation = body.Config.JoinConfirmation
	party.Config.Joinability = body.Config.Joinability
	party.Config.MaxSize = body.Config.MaxSize
	party.Members = []models.V2PartyMember{partyMember}

	for key, metaItem := range body.Meta {
		party.Meta[key] = metaItem
	}

	ActiveParties[party.ID] = party

	c.JSON(200, party)

	deleteAnyEmptyParties()
}

func PartyPatch(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var body struct {
		Config struct {
			JoinConfirmation bool `json:"join_confirmation"`
			Joinability string `json:"joinability"`
			MaxSize int `json:"max_size"`
		} `json:"config"`
		JoinInfo struct {
			Connection struct {
				Id string `json:"id"`
				Meta struct {
					UrnEpicConnPlatformS string `json:"urn:epic:conn:platform_s"`
					UrnEpicConnTypeS string `json:"urn:epic:conn:type_s"`
				} `json:"meta"`
			} `json:"connection"`
			Meta struct {
				UrnEpicMemberDnS string `json:"urn:epic:member:dn_s"`
				UrnEpicMemberTypeS string `json:"urn:epic:member:type_s"`
				UrnEpicMemberPlatformS string `json:"urn:epic:member:platform_s"`
			} `json:"meta"`
		} `json:"join_info"`
		Meta struct {
			Update map[string]interface{} `json:"update"`
			Delete []string `json:"delete"`
		} `json:"meta"`
	}

	if err := c.BindJSON(&body); err != nil {
		common.ErrorBadRequest(c)
		return
	}

	partyId := AccountIdToPartyId[user.AccountId]
	party, ok := ActiveParties[partyId]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	for _, key := range body.Meta.Delete {
		delete(party.Meta, key)
	}

	for key, metaItem := range body.Meta.Update {
		party.Meta[key] = metaItem
	}

	party.Config.JoinConfirmation = body.Config.JoinConfirmation
	party.Config.Joinability = body.Config.Joinability
	party.Config.MaxSize = body.Config.MaxSize

	ActiveParties[partyId] = party
	c.JSON(200, party)
}

func PartyPatchMemberMeta(c *gin.Context) {
	partyId := c.Param("partyId")
	memberId := c.Param("memberId")

	party, ok := ActiveParties[partyId]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	var body struct {
		Update map[string]interface{} `json:"update"`
		Delete []string `json:"delete"`
	}

	if err := c.BindJSON(&body); err != nil {
		common.ErrorBadRequest(c)
		return
	}

	for _, member := range party.Members {
		if member.AccountId == memberId {
			for _, key := range body.Delete {
				delete(member.Meta, key)
			}

			for key, metaItem := range body.Update {
				member.Meta[key] = metaItem
			}

			break
		}
	} 

	ActiveParties[partyId] = party
	c.JSON(200, party)
}

func PartyGet(c *gin.Context) {
	partyId := c.Param("partyId")
	party, ok := ActiveParties[partyId]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	all.PrintMagenta([]any{
		"PartyGet",
		partyId,
	})

	c.JSON(200, party)
}

func PartyDeleteMember(c *gin.Context) {
	partyId := c.Param("partyId")
	memberId := c.Param("memberId")

	party, ok := ActiveParties[partyId]
	if !ok {
		common.ErrorBadRequest(c)
		return
	}

	for i, member := range party.Members {
		if member.AccountId == memberId {
			party.Members = append(party.Members[:i], party.Members[i+1:]...)
			break
		}
	}

	if len(party.Members) == 0 {
		delete(ActiveParties, partyId)
	}

	delete(AccountIdToPartyId, memberId)
	ActiveParties[partyId] = party
	c.JSON(200, party)
}

func deleteAnyEmptyParties() {
	for partyId, party := range ActiveParties {
		if len(party.Members) == 0 {
			delete(ActiveParties, partyId)
		}
	}
}