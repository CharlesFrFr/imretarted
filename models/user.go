package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	AccountId   string `gorm:"uniqueIndex;default:null"`
	Username    string
	Password    string
	DiscordId   uint32 `gorm:"default:0"`
	Banned      bool   `gorm:"default:false"`
	LoginSecret string `gorm:"default:null"`
}