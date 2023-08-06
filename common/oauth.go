package common

import (
	"encoding/base64"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/zombman/server/all"
	"github.com/zombman/server/models"
)

func GenerateClientToken(client string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"clsvc": "fortnite",
		"t": "s",
		"mver": false,
		"clid": client,
		"am": "client_credentials",
		"ic": true,
		"p": base64.StdEncoding.EncodeToString([]byte(uuid.New().String())),
		"jti": strings.ReplaceAll(uuid.New().String(), "-", ""),
		"creation_date": time.Now().Unix(),
		"hours_expire":  24,
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	
	return strings.Join([]string{"eg1~", tokenString}, "")
}

func GetClientToken(ip string) (models.ClientToken, error) {
	var clientToken models.ClientToken

	result := all.Postgres.Where("ip = ?", ip).First(&clientToken)

	if result.Error != nil {
		return models.ClientToken{}, result.Error
	}

	return clientToken, nil
}

func GenerateAccessToken(user models.User, clientId string, device string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"clsvc": "fortnite",
		"app": "fortnite",
		"iai": user.AccountId,
		"dn": user.Username,
		"sec": 1,
		"t": "s",
		"mver": false,
		"clid": clientId,
		"am": "password",
		"ic": true,
		"p": base64.StdEncoding.EncodeToString([]byte(uuid.New().String())),
		"dvid": device,
		"jti": strings.ReplaceAll(uuid.New().String(), "-", ""),
		"creation_date": time.Now().Unix(),
		"hours_expire":  24,
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	
	return strings.Join([]string{"eg1~", tokenString}, "")
}

func GetAccessToken(accountId string) (models.AccessToken, error) {
	var accessToken models.AccessToken

	result := all.Postgres.Where("account_id = ?", accountId).First(&accessToken)

	if result.Error != nil {
		return models.AccessToken{}, result.Error
	}

	return accessToken, nil
}

func GetSiteToken(accountId string) (models.SiteToken, error) {
	var siteToken models.SiteToken

	result := all.Postgres.Where("account_id = ?", accountId).First(&siteToken)

	if result.Error != nil {
		return models.SiteToken{}, result.Error
	}

	return siteToken, nil
}

func GenerateRefreshToken(user models.User, clientId string, device string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.AccountId,
		"t": "r",
		"clid": clientId,
		"am": "refresh_token",
		"dvid": device,
		"jti": strings.ReplaceAll(uuid.New().String(), "-", ""),
		"creation_date": time.Now().Unix(),
		"hours_expire": 24 * 30,
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	
	return strings.Join([]string{"eg1~", tokenString}, "")
}

func GetRefreshToken(accountId string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	result := all.Postgres.Where("account_id = ?", accountId).First(&refreshToken)

	if result.Error != nil {
		return models.RefreshToken{}, result.Error
	}

	return refreshToken, nil
}

func GetRefreshTokenWithToken(token string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	result := all.Postgres.Where("token = ?", token).First(&refreshToken)

	if result.Error != nil {
		return models.RefreshToken{}, result.Error
	}

	return refreshToken, nil
}