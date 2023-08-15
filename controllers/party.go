package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

var (
	ActiveParties = make(map[string]models.V2Party)
	AccountIdToPartyId = make(map[string]string)
)

func PartyGetUser(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	partyId, ok := AccountIdToPartyId[user.AccountId]
	if !ok {
		c.JSON(400, gin.H{"error": "User not in party"})
		return
	}

	party, ok := ActiveParties[partyId]
	if !ok {
		c.JSON(400, gin.H{"error": "Party not found"})
		return
	}

	c.JSON(200, party)
}

func PartyGetFriendPartyPings(c *gin.Context) {
	friend, err := common.GetUserByAccountId(c.Param("friendId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	partyId, ok := AccountIdToPartyId[friend.AccountId]
	if !ok {
		c.JSON(400, gin.H{"error": "User not in party"})
		return
	}

	party, ok := ActiveParties[partyId]
	if !ok {
		c.JSON(400, gin.H{"error": "Party not found"})
		return
	}

	c.JSON(200, []models.V2Party{party})
}

func PartyLeave(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	partyId, ok := AccountIdToPartyId[user.AccountId]
	if !ok {
		c.JSON(400, gin.H{"error": "User not in party"})
		return
	}

	party, ok := ActiveParties[partyId]
	if !ok {
		c.JSON(400, gin.H{"error": "Party not found"})
		return
	}

	delete(AccountIdToPartyId, user.AccountId)

	for memberIndex, partyMember := range party.Members {
		if partyMember.AccountId == user.AccountId {
			party.Members = append(party.Members[:memberIndex], party.Members[memberIndex+1:]...)
			break
		}
	}

	ActiveParties[partyId] = party
	c.Status(204)
}