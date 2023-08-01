package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/common"
	"github.com/zombman/server/controllers"
	"github.com/zombman/server/helpers"
)

func init() {
  helpers.LoadEnviroment()
  helpers.ConnectToDatabase()
  helpers.AutoMigrate()
}

func main() {
  r := gin.Default()

  r.Use(func(c *gin.Context) {
    if helpers.Postgres == nil {
		  c.JSON(http.StatusInternalServerError, gin.H{"error": "database not connected"})
		  c.Abort()
	  }

    c.Next()
  })

  r.POST("/api/user/login", controllers.UserLogin)
  r.POST("/api/user/create", controllers.UserCreate)

  r.POST("/account/api/oauth/token", controllers.OAuthMain)

  r.GET("/waitingroom/api/waitingroom", func(c *gin.Context) {
    c.Status(204)
    c.Abort()
  })

  fmt.Println(common.GenerateClientToken("hello world"))

  r.Run()
}