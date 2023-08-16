package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/zombman/server/models"
)

func CreateParty(activeParties *map[string]models.V2Party, accountIdToPartyId *map[string]string, captain models.User) models.V2Party {
	partyMeta := make(map[string]interface{})
	partyMeta["urn:epic:cfg:build-id_s"] = "1:1:1"
	partyMeta["urn:epic:cfg:party-type-id_s"] = "default"
	partyMeta["Default:PartyState_s"] = "BattleRoyaleView"
	partyMeta["urn:epic:cfg:join-request-action_s"] = "Manual"
	partyMeta["urn:epic:cfg:accepting-members_b"] = true
	party := models.V2Party{
		ID: uuid.New().String(),
		Meta: partyMeta,
		Config: models.V2PartyConfig{
			MaxSize: 16,
			SubType: "default",
			Type: "default",
			Joinability: "OPEN",
			Discoverability: "ALL",
			InviteTtl: 86400,
			JoinConfirmation: false,
		},
		Members: []models.V2PartyMember{},
		Revision: 0,
		Invites: []any{},
		Intentions: []any{},	
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		CreatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
	}

	(*activeParties)[party.ID] = party
	(*accountIdToPartyId)[captain.AccountId] = party.ID

	return party
}