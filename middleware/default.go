package middleware

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
)

func CheckDatabase(c *gin.Context) {
	if all.Postgres == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not connected"})
		c.Abort()
	}

	c.Next()
}

func AllowFromAnywhere(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

var Limit = tollbooth.NewLimiter(30, &limiter.ExpirableOptions{ DefaultExpirationTTL: time.Minute })

func RateLimitMiddleware(c *gin.Context) {
	httpError := tollbooth.LimitByRequest(Limit, c.Writer, c.Request)
	if httpError != nil {
		c.JSON(httpError.StatusCode, gin.H{"error": httpError.Message})
		c.Abort()
		return
	}

	c.Next()
}