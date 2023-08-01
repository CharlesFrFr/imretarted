package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/controllers"
	"github.com/zombman/server/helpers"
)

func init() {
  helpers.LoadEnviroment()
  helpers.ConnectToDatabase()
}

func main() {
  r := gin.Default()

  r.Use(func(c *gin.Context) {
    if helpers.Postgres == nil {
		  c.JSON(http.StatusBadRequest, gin.H{"error": "database not connected"})
		  c.Abort()
	  }

    c.Next()
  })

  r.POST("/api/user/create", controllers.UserCreate)

  r.Run()
}