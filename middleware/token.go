package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
)

func VerifyAccessToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}
	tokenString = tokenString[11:]
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
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
	
	all.PrintGreen([]any{"token verified for account", accountId})

	user, err := common.GetUserByAccountId(accountId)

	if err != nil {
		fmt.Println("db fail to get user:", err)
		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}