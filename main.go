package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/controllers"
	"github.com/zombman/server/helpers"
	"github.com/zombman/server/middleware"
)

func init() {
  helpers.LoadEnviroment()
  helpers.ConnectToDatabase()
  helpers.AutoMigrate()
}

func main() {
  helpers.PrintGreen([]string{"development mode"})

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

  r.GET("/account/api/public/account", controllers.UserAccountPublic)
  r.GET("/account/api/public/account/:accountId", middleware.VerifyAccessToken, controllers.UserAccountPrivate)
  r.GET("/account/api/public/account/:accountId/externalAuths", func(c *gin.Context) {
    helpers.PrintGreen([]string{"externalAuths"})

    c.JSON(http.StatusOK, []string{})
  })
  r.DELETE("/account/api/oauth/sessions/kill/:token", controllers.KillSessionWithToken)
  r.DELETE("/account/api/oauth/sessions/kill", controllers.KillSession)
  
  // does not need to be implemented

  r.GET("/fortnite/api/game/v2/enabled_features", func(c *gin.Context) {
    c.JSON(http.StatusOK, []gin.H{})
  })
  r.POST("/fortnite/api/game/v2/tryPlayOnPlatform/account/*accountId", func(c *gin.Context) {
    c.String(http.StatusOK, "true")
  })
  r.GET("/waitingroom/api/waitingroom", func(c *gin.Context) {
    c.Status(204)
    c.Abort()
  })
  r.POST("/datarouter/*api", func(c *gin.Context) {
    c.Status(204)
    c.Abort()
  })

  r.GET("/lightswitch/api/service/bulk/status", func(c *gin.Context) {
    c.JSON(http.StatusOK, []gin.H{{
        "serviceInstanceId": "fortnite",
        "status": "UP",
        "message": "fortnite is up.",
        "maintenanceUri": nil,
        "overrideCatalogIds": []string{"a7f138b2e51945ffbfdacc1af0541053"},
        "allowedActions": []string{"PLAY", "DOWNLOAD"},
        "banned": false,
        "launcherInfoDTO": gin.H{
          "appName": "Fortnite",
          "catalogItemId": "4fe75bbc5a674f4f9b356b5c90567da5",
          "namespace": "fn",
        },
    }})
  })
  r.GET("/lightswitch/api/service/Fortnite/status", func(c *gin.Context) {
    c.JSON(http.StatusOK, []gin.H{{
        "serviceInstanceId": "fortnite",
        "status": "UP",
        "message": "fortnite is up.",
        "maintenanceUri": nil,
        "overrideCatalogIds": []string{"a7f138b2e51945ffbfdacc1af0541053"},
        "allowedActions": []string{"PLAY", "DOWNLOAD"},
        "banned": false,
        "launcherInfoDTO": gin.H{
          "appName": "Fortnite",
          "catalogItemId": "4fe75bbc5a674f4f9b356b5c90567da5",
          "namespace": "fn",
        },
    }})
  })

  r.Run()
}