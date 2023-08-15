package common

import (
	"time"

	"github.com/zombman/server/models"
)

func CreateParty(activeParties *map[string]models.V2Party, accountIdToPartyId *map[string]string, captainId string) {
	partyMemberMeta := make(map[string]interface{})

	_ = models.V2PartyMember{
		AccountId: captainId,
		Meta: partyMemberMeta,
		Connections: []models.V2PartyConnection{{}},
		Role: "CAPTAIN",
		Revision: 0,
		JoinedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: time.Now().Format("2006-01-02T15:04:05.999Z"),
	}
}