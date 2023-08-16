package controllers

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
)

func KillSession(c *gin.Context) {
	// when xmpp is implemented, this will be used to kill other sessions
	c.Status(204)
  c.Abort()
}

func KillSessionWithToken(c *gin.Context) {
	tokenString := c.Param("token")
	if tokenString == "" {
		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}
	tokenString = tokenString[4:]
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		all.MarshPrintJSON(tokenString)
		fmt.Println("jwt parse error:", err)
		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}

	accountId := token.Claims.(jwt.MapClaims)["iai"].(string)
	dbToken, err := common.GetAccessToken(accountId)

	if err != nil {
		fmt.Println("db fail to get token:", err)
		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}
	
	if dbToken.Token != strings.Join([]string{"eg1~", tokenString}, "") {
		fmt.Println("token not match")

		all.PrintRed([]any{"dbToken", dbToken.Token})
		all.PrintGreen([]any{"tokenString", strings.Join([]string{"eg1~", tokenString}, "")})

		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}

	partyId, ok := AccountIdToPartyId[accountId]
	all.PrintGreen([]any{"partyId", partyId, accountId})
	if ok {
		delete(ActiveParties, partyId)
		delete(AccountIdToPartyId, accountId)
	}

	c.Status(204)
  c.Abort()
}