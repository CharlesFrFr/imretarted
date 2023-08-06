package models

import "gorm.io/gorm"

type ClientToken struct {
	gorm.Model
	IP string `json:"ip"`
	Token string `json:"token"`
}

type AccessToken struct {
	gorm.Model
	AccountId string `json:"accountId"`
	Token string `json:"token"`
}

type RefreshToken struct {
	gorm.Model
	AccountId string `json:"accountId"`
	Token string `json:"token"`
}

type SiteToken struct {
	gorm.Model
	AccountId string `json:"accountId"`
	Token string `json:"token"`
}

type SiteRefreshToken struct {
	gorm.Model
	AccountId string `json:"accountId"`
	Token string `json:"token"`
}