package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func DefaultEpicError(c *gin.Context, code string, message string, numberCode int, err string, statusCode int) {
	c.Header("X-Epic-Error-Code", fmt.Sprint(numberCode))
	c.Header("X-Epic-Error-Name", code)

	c.JSON(statusCode, gin.H{
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
	DefaultEpicError(c, "errors.com.epicgames.account.invalid_account_credentials", "Your username and/or password are incorrect. Please check them and try again.", 18031, "invalid_grant", 401)
}

func ErrorInvalidOAuthRequest(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.oauth.invalid_request", "Invalid Request", 1013, "invalid_request", 400)
}

func ErrorAuthFailed(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.authorization.authorization_failed", "auth Failed", 1032, "", 401)
}

func ErrorBadRequest(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.bad_request", "Bad Request", 1000, "", 400)
}

func ErrorItemNotFound(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.item_not_found", "Item not found", 1004, "", 404)
}

func ErrorNameTaken(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.account.account_name_taken", "Sorry, that display name is already taken.", 18006, "", 400)
}

func ErrorInternalServer(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.server_error", "Internal Server Error", 10000, "", 500)
}

func ErrorUnauthorized(c *gin.Context) {
	DefaultEpicError(c, "errors.com.epicgames.common.oauth.unauthorized", "Unauthorized", 1002, "", 403)
}