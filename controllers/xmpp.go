package controllers

import "github.com/gin-gonic/gin"

func KillSession(c *gin.Context) {
	// when xmpp is implemented, this will be used to kill other sessions
	c.Status(204)
  c.Abort()
}

func KillSessionWithToken(c *gin.Context) {
	// when xmpp is implemented, this will be used to kill other sessions
	c.Status(204)
  c.Abort()
}