package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/zombman/server/models"
)

func CreateParty(activeParties *map[string]models.V2Party, accountIdToPartyId *map[string]string, captainId string) models.V2Party {
	partyMemberMeta := make(map[string]interface{})
	partyMember := models.V2PartyMember{
		AccountId: captainId,
		Meta: partyMemberMeta,
		Connections: []models.V2PartyConnection{{}},
		Role: "CAPTAIN",
		Revision: 0,
		JoinedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
	}

	partyMeta := make(map[string]interface{})
	party := models.V2Party{
		ID: uuid.New().String(),
		Meta: partyMeta,
		Config: models.V2PartyConfig{
			Joinability: "OPEN",
			MaxSize: 16,
			SubType: "default",
			Type: "default",
			Discoverability: "ALL",
			InviteTtl: 86400,
			JoinConfirmation: "false",
		},
		Members: []models.V2PartyMember{partyMember},
		Revision: 0,
		Invites: []any{},
		Intentions: []any{},	
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		CreatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
	}

	(*activeParties)[party.ID] = party
	(*accountIdToPartyId)[captainId] = party.ID

	return party
}