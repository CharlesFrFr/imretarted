package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func DefaultEpicError(c *gin.Context, code string, message string, numberCode int, err string) {
	c.Header("X-Epic-Error-Code", fmt.Sprint(numberCode))
	c.Header("X-Epic-Error-Name", code)

	c.JSON(400, gin.H{
		"error": err,
		"errorCode": code,
		"errorMessage": message,
		"errorDescription": message,
		"numericErrorCode": numberCode,
		"originatingService": "com.epicgames.account.public",
		"intent": "prod",
		"messageVars": []string{},
	})
	c.Abort()
}

func ErrorInvalidCredentials(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.account.invalid_account_credentials", "Your username and/or password are incorrect. Please check them and try again.", 18031, "invalid_grant")
}

func ErrorInvalidOAuthRequest(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.oauth.invalid_request", "Invalid Request", 1013, "invalid_request")
}

func ErrorAuthFailed(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.authorization.authorization_failed", "auth Failed", 1032, "")
}

func ErrorBadRequest(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.bad_request", "Bad Request", 1000, "")
}

func ErrorItemNotFound(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.item_not_found", "Item not found", 1004, "")
}