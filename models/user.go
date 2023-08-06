package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AccountId   string `gorm:"uniqueIndex;default:null"`
	Username    string `gorm:"unique;default:null"`
	Password    string `gorm:"default:null"`
	AccessLevel int    `gorm:"default:0"`
	DiscordId   uint32 `gorm:"default:0"`
	Banned      bool   `gorm:"default:false"`
	VBucks 			int    `gorm:"default:0"`
}