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

  account := r.Group("/account/api")
  {
    account.GET("/public/account", controllers.UserAccountPublic)
    account.GET("/public/account/:accountId", middleware.VerifyAccessToken, controllers.UserAccountPrivate)
    account.GET("/public/account/:accountId/externalAuths", func(c *gin.Context) { c.JSON(http.StatusOK, []string{}) })
    account.DELETE("/oauth/sessions/kill/:token", controllers.KillSessionWithToken)
    account.DELETE("/oauth/sessions/kill", controllers.KillSession)
  }

  fortnite := r.Group("/fortnite/api")
  {
    fortnite.POST("/game/v2/profile/:accountId/client/:action", middleware.VerifyAccessToken, controllers.ProfileQuery)
    fortnite.GET("/game/v2/enabled_features", controllers.EmptyArray)
    fortnite.POST("/game/v2/tryPlayOnPlatform/account/*accountId", controllers.True)
    fortnite.GET("/storefront/v2/keychain", controllers.StorefrontKeychain)
  }

  blank := r.Group("/")
  {
    blank.GET("/waitingroom/api/waitingroom", controllers.NoResponse)
    blank.POST("/datarouter/*api", controllers.NoResponse)
    blank.GET("/lightswitch/api/service/bulk/status", controllers.Lightswitch)
    blank.GET("/lightswitch/api/service/Fortnite/status", controllers.Lightswitch)
  }

  r.Run()
}