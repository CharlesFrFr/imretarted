package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
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

	user, err := common.GetUserByAccountId(accountId)

	if err != nil {
		fmt.Println("db fail to get user:", err)
		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}

	all.PrintYellow([]any{"token verified for account", accountId})

	c.Set("user", user)
	c.Next()
}

func VerifyAccessTokenXMPP(tokenString string) (models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		all.MarshPrintJSON(tokenString)
		fmt.Println("jwt parse error:", err)
		return models.User{}, err
	}

	accountId := token.Claims.(jwt.MapClaims)["iai"].(string)
	dbToken, err := common.GetAccessToken(accountId)

	if err != nil {
		fmt.Println("db fail to get token:", err)
		return models.User{}, err
	}
	
	if dbToken.Token != strings.Join([]string{"eg1~", tokenString}, "") {
		fmt.Println("token not match")

		all.PrintRed([]any{"dbToken", dbToken.Token})
		all.PrintGreen([]any{"tokenString", strings.Join([]string{"eg1~", tokenString}, "")})

		return models.User{}, err
	}

	user, err := common.GetUserByAccountId(accountId)

	if err != nil {
		fmt.Println("db fail to get user:", err)
		return models.User{}, err
	}

	all.PrintYellow([]any{"token verified for account", accountId})

	return user, nil
}

func VerifySiteToken(c *gin.Context) {
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
		all.MarshPrintJSON(tokenString)
		fmt.Println("jwt parse error:", err)
		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}

	accountId := token.Claims.(jwt.MapClaims)["iai"].(string)
	dbToken, err := common.GetSiteToken(accountId)

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

	user, err := common.GetUserByAccountId(accountId)

	if err != nil {
		fmt.Println("db fail to get user:", err)
		common.ErrorAuthFailed(c)
		c.Abort()
		return
	}

	all.PrintYellow([]any{"token verified for account", accountId})

	c.Set("user", user)
	c.Next()
}