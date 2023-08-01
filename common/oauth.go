package common

import (
	"encoding/base64"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/zombman/server/models"
)

func GenerateClientToken(clientId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"clsvc": "fortnite",
		"t": "s",
		"mver": false,
		"clid": clientId,
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

func GenerateAccessToken(user models.User, clientId string) string {
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
		"dvid": strings.ReplaceAll(uuid.New().String(), "-", ""),
		"jti": strings.ReplaceAll(uuid.New().String(), "-", ""),
		"creation_date": time.Now().Unix(),
		"hours_expire":  24,
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	
	return strings.Join([]string{"eg1~", tokenString}, "")
}

func GenerateRefreshToken(user models.User, clientId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.AccountId,
		"t": "r",
		"clid": clientId,
		"am": "refresh_token",
		"dvid": strings.ReplaceAll(uuid.New().String(), "-", ""),
		"jti": strings.ReplaceAll(uuid.New().String(), "-", ""),
		"creation_date": time.Now().Unix(),
		"hours_expire":  24,
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	
	return strings.Join([]string{"eg1~", tokenString}, "")
}