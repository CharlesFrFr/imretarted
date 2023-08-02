package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

func UserCreate(c *gin.Context) {
	var body struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := common.CreateUser(body.Username, body.Password)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UserLogin(c *gin.Context) {
	var body struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := common.GetUserByUsernameAndPlainPassword(body.Username, body.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UserAccountPrivate(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	c.JSON(http.StatusOK, gin.H{
		"id": user.AccountId,
		"email": strings.Join([]string{user.Username, "@."}, ""),
		"emailVerified": true,
		"minorVerified": false,
		"minorStatus": "NOT_MINOR",
		"cabinedMode": false,
		"hasHashedEmail": false,
		"displayName": user.Username,
		"canUpdateDisplayName": false,
		"numberOfDisplayNameChanges": 0,
		"name": user.Username,
		"lastName": user.Username,
		"country": "US",
		"preferredLanguage": "en",
		"failedLoginAttempts": 0,
		"lastLogin": time.Now().Format("2006-01-02T15:04:05.999Z"),
		"ageGroup": "UNKNOWN",
		"headless": false,
	})
}

type UserAccountPublicResponse struct {
	Id string `json:"id"`
	DisplayName string `json:"displayName"`
	ExternalAuths interface{} `json:"externalAuths"`
}

func UserAccountPublic(c *gin.Context) {
	response := [](UserAccountPublicResponse){}
	
	accountId := c.Query("accountId")
	if accountId != "" {
		user, err := common.GetUserByAccountId(accountId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		response = append(response, UserAccountPublicResponse{
			Id: user.AccountId,
			DisplayName: user.Username,
			ExternalAuths: []string{},
		})
	}

	c.JSON(http.StatusOK, response)
}