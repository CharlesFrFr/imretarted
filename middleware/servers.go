package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ServerSecret(c *gin.Context) {
	if c.GetHeader("X-Server-Secret") != os.Getenv("SECRET") {
		c.AbortWithStatus(401)
		return
	}
	c.Next()
}