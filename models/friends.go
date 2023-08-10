package models

import (
	"gorm.io/gorm"
)

type FriendAction struct {
	gorm.Model
	ForAccountId string
	AccountId string
	Action string  // ACCEPTED, INCOMING, OUTGOING, BLOCKED
}