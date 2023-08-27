package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

type Body struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Username     string	`form:"username"`
	Password     string	`form:"password"`
	ExchangeCode string	`form:"exchange_code"`
	RefreshToken string	`form:"refresh_token"`
}

func OAuthMain(c *gin.Context) {
	var body Body

	client := c.GetHeader("Authorization")
	if client == "" {
		common.ErrorInvalidOAuthRequest(c)
		return
	}
	if len(strings.Split(client, " ")) <= 1 {
		common.ErrorInvalidOAuthRequest(c)
		return
	}
	client = strings.Split(all.DecodeBase64(strings.Split(client, " ")[1]), ":")[0]
	
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	all.PrintRed([]any{"grant_type: ", body.GrantType})
	switch body.GrantType {
		case "client_credentials": 
			ClientCredentials(c, client)
		case "password":
			Password(c, body, client)
		case "refresh_token":
			RefreshToken(c, body, client)
		default: 
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid grant_type"})
	}
}

func Generate(user models.User, client string) gin.H {
	device := strings.ReplaceAll(uuid.New().String(), "-", "")
	accessToken := common.GenerateAccessToken(user, client, device)
	refreshToken := common.GenerateRefreshToken(user, client, device)

	all.Postgres.
		Where(models.AccessToken{AccountId: user.AccountId}).
		Assign(models.AccessToken{Token: accessToken}).
		FirstOrCreate(&models.AccessToken{})

	all.Postgres.
		Where(models.RefreshToken{AccountId: user.AccountId}).
		Assign(models.RefreshToken{Token: refreshToken}).
		FirstOrCreate(&models.RefreshToken{})

	return gin.H{
		"app": "fortnite",
		"account_id": user.AccountId,
		"device_id": device,
		"client_id": client,
		"client_service": "fortnite",
		"internal_client": true,
		"displayName": user.Username,
		"access_token": accessToken,
		"token_type": "bearer",
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": time.Hour.Milliseconds() * 24,
		"refresh_token": refreshToken,
		"refresh_expires": time.Hour.Milliseconds() * 24,
		"refresh_expires_at": time.Now().Add(time.Hour * 24 * 30).Format("2006-01-02T15:04:05.999Z"),
	}
}

func GenerateSiteToken(user models.User, client string) gin.H {
	device := strings.ReplaceAll(uuid.New().String(), "-", "")
	accessToken := common.GenerateAccessToken(user, client, device)
	refreshToken := common.GenerateRefreshToken(user, client, device)

	all.Postgres.
		Where(models.SiteToken{AccountId: user.AccountId}).
		Assign(models.SiteToken{Token: accessToken}).
		FirstOrCreate(&models.SiteToken{})

	all.Postgres.
		Where(models.SiteRefreshToken{AccountId: user.AccountId}).
		Assign(models.SiteRefreshToken{Token: refreshToken}).
		FirstOrCreate(&models.SiteRefreshToken{})

	return gin.H{
		"app": "fortnite",
		"account_id": user.AccountId,
		"device_id": device,
		"client_id": client,
		"client_service": "fortnite",
		"internal_client": true,
		"displayName": user.Username,
		"access_token": accessToken,
		"token_type": "bearer",
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": time.Hour.Milliseconds() * 24,
		"refresh_token": refreshToken,
		"refresh_expires": time.Hour.Milliseconds() * 24,
		"refresh_expires_at": time.Now().Add(time.Hour * 24 * 30).Format("2006-01-02T15:04:05.999Z"),
	}
}

func RefreshToken(c *gin.Context, body Body, client string) {
	refreshToken, err := common.GetRefreshTokenWithToken(body.RefreshToken)
	if err != nil {
		common.ErrorInvalidCredentials(c)
		return
	}

	user, err := common.GetUserByAccountId(refreshToken.AccountId)
	if err != nil {
		common.ErrorInvalidCredentials(c)
		return
	}

	c.JSON(http.StatusOK, Generate(user, client))
}

func Password(c *gin.Context, body Body, client string) {
	user, err := common.GetUserByUsernameAndPlainPassword(strings.ReplaceAll(body.Username, "@.", ""), body.Password)
	if err != nil {
		user, err = common.GetUserByUsernameAndHashPassword(strings.ReplaceAll(body.Username, "@.", ""), body.Password)
		if err != nil {
			common.ErrorInvalidCredentials(c)
			return
		}
	}

	c.JSON(http.StatusOK, Generate(user, client))
}

func ClientCredentials(c *gin.Context, client string) {
	ip := c.ClientIP()
	existingClientToken, _ := common.GetClientToken(ip)

	if existingClientToken.ID != 0 {
		all.Postgres.Delete(&existingClientToken)
	}

	clientToken := common.GenerateClientToken(client)
	all.Postgres.Create(&models.ClientToken{
		IP: ip,
		Token: clientToken,
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token": clientToken,
		"token_type": "bearer",
		"client_id": client,
		"client_service": "fortnite",
		"internal_client": true,
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": time.Hour.Milliseconds() * 24,
	})
}

func OAuthVerify(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	all.MarshPrintJSON(user)
	c.AbortWithStatus(204)
}