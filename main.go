package main

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/controllers"
	"github.com/zombman/server/middleware"
)

func init() {
  cmd := exec.Command("cmd", "/c", "cls")
  cmd.Stdout = os.Stdout
  cmd.Run()

  all.LoadEnviroment()
  all.ConnectToDatabase()
  all.AutoMigrate()
}

func main() {
  all.PrintGreen([]any{"development mode"})
  r := gin.Default()
  r.Use(func(c *gin.Context) {
    if all.Postgres == nil {
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
    account.GET("/public/account/:accountId/externalAuths", controllers.EmptyArray)
    account.DELETE("/oauth/sessions/kill/:token", controllers.KillSessionWithToken)
    account.DELETE("/oauth/sessions/kill", controllers.KillSession)
  }

  fortnite := r.Group("/fortnite/api")
  {
    fortnite.POST("/game/v2/profile/:accountId/client/:action", middleware.VerifyAccessToken, controllers.ProfileActionHandler)
    fortnite.POST("/game/v2/tryPlayOnPlatform/account/*accountId", middleware.VerifyAccessToken, controllers.True)
    fortnite.GET("/game/v2/enabled_features", middleware.VerifyAccessToken, controllers.EmptyArray)
    fortnite.GET("/receipts/v1/account/:accountId/receipts", middleware.VerifyAccessToken, controllers.EmptyArray)
    fortnite.GET("/storefront/v2/keychain", middleware.VerifyAccessToken, controllers.StorefrontKeychain)
    fortnite.GET("/calendar/v1/timeline", controllers.CalendarTimeline)

    fortnite.GET("/cloudstorage/system", controllers.SystemCloudFilesList)
    fortnite.GET("/cloudstorage/system/:fileName", controllers.SystemCloudFile)
    fortnite.GET("/cloudstorage/user/:accountId", middleware.VerifyAccessToken, controllers.UserCloudFilesList)
    fortnite.GET("/cloudstorage/user/:accountId/:fileName", middleware.VerifyAccessToken, controllers.UserCloudFile)
    fortnite.PUT("/cloudstorage/user/:accountId/ClientSettings.Sav", middleware.VerifyAccessToken, controllers.SaveUserCloudFile)

    store := fortnite.Group("/storefront")
    {
      store.GET("/v2/catalog", controllers.StorefrontCatalog)
    }
  }

  blank := r.Group("/")
  {
    blank.GET("/content/api/pages/*contentPageName", middleware.VerifyAccessToken, controllers.ContentPage)
    blank.GET("/waitingroom/api/waitingroom", controllers.NoResponse)
    blank.POST("/datarouter/*api", controllers.NoResponse)
    blank.GET("/lightswitch/api/service/bulk/status", controllers.Lightswitch)
    blank.GET("/lightswitch/api/service/Fortnite/status", controllers.Lightswitch)
  }

  r.Run()
}