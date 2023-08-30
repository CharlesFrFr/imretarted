package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
	"github.com/zombman/server/socket"
)

func UserCreate(c *gin.Context) {
	var body struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Captcha  	string `json:"captcha"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok := common.VerifyGoogleRecaptcha(body.Captcha); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You are not a human. Please refresh the page and try again."})
		return
	}

	user, err := common.CreateUser(body.Username, body.Password, 0)
	if err != nil {
		common.ErrorNameTaken(c)
		return
	}

	token := GenerateSiteToken(user, "site")
	c.JSON(http.StatusOK, gin.H{"data": user, "token": token})
}

func UserLogin(c *gin.Context) {
	var body struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := common.GetUserByUsernameAndPlainPassword(body.Username, body.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	token := GenerateSiteToken(user, "site")
	c.JSON(http.StatusOK, gin.H{"data": user, "token": token})
}

func UserAccountPrivate(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	c.JSON(http.StatusOK, gin.H{
		"id": user.AccountId,
		"email": strings.Join([]string{user.Username, "@."}, ""),
		"emailVerified": true,
		"minorVerified": false,
		"minorStatus": "NOT_MINOR",
		"cabinedMode": false,
		"hasHashedEmail": false,
		"displayName": user.Username,
		"canUpdateDisplayName": false,
		"numberOfDisplayNameChanges": 0,
		"name": user.Username,
		"lastName": user.Username,
		"country": "US",
		"preferredLanguage": "en",
		"failedLoginAttempts": 0,
		"lastLogin": time.Now().Format("2006-01-02T15:04:05.999Z"),
		"ageGroup": "UNKNOWN",
		"headless": false,
	})
}

type UserAccountPublicResponse struct {
	Id string `json:"id"`
	DisplayName string `json:"displayName"`
	ExternalAuths interface{} `json:"externalAuths"`
}

func GetGoogleRecaptcha(c *gin.Context) {
	c.JSON(http.StatusOK, os.Getenv("GOOGLE_RECAPTCHA_SITE_KEY"))
}

// for some reason i think it is party v2 related
func UserAccountPublic(c *gin.Context) {
	response := [](UserAccountPublicResponse){}

	accountIds, _ := c.GetQueryArray("accountId")
	for _, accountId := range accountIds {
		user, err := common.GetUserByAccountId(accountId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		response = append(response, UserAccountPublicResponse{
			Id: user.AccountId,
			DisplayName: user.Username,
			ExternalAuths: []string{},
		})
	}

	c.JSON(http.StatusOK, response)
}

func UserAccountPublicFromDisplayName(c *gin.Context) {
	username := c.Param("displayName")
	all.PrintMagenta([]any{"username", username})
	if username != "" {
		user, err := common.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": user.AccountId,
			"displayName": user.Username,
			"externalAuths": []string{},
		})
		return
	}

	common.ErrorBadRequest(c)
}

type LockerItem struct {
	ItemId string `json:"itemId"`
	Rarity string `json:"rarity"`
	Season int `json:"season"`
}

func UserGetLocker(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	profile, err := common.ReadProfileFromUser(user.AccountId, "athena")
	if err != nil {
		
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	marshalledItems, err := json.Marshal(profile.Items)
	if err != nil {
		all.PrintGreen([]any{"serre", err})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	var items map[string]models.Item
	err = json.Unmarshal(marshalledItems, &items)
	if err != nil {
		all.PrintGreen([]any{"serre2222", err})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	if common.AllItems == nil {
		common.GetAllFortniteItems()
	}

	locker := []LockerItem{}
	for _, item := range items {
		if item.TemplateId == "Currency:MtxPurchased" || item.TemplateId == "CosmeticLocker:cosmeticlocker_athena" {
			continue
		}

		if len(common.AllItemsKeys) == 0 {
			common.GetAllFortniteItems()
		}

		itemInfo, ok := common.AllItemsKeys[item.TemplateId]
		if !ok {
			continue
		}

		if itemInfo.IntroductionSeason > common.Season + 10 && itemInfo.Rarity != "Mythic" {
			continue
		}

		locker = append(locker, LockerItem{
			ItemId: item.TemplateId,
			Rarity: itemInfo.Rarity,
			Season: itemInfo.IntroductionSeason,
		})
	}

	c.JSON(http.StatusOK, locker)
}

func SiteRefresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBind(&body); err != nil {
		common.ErrorInvalidOAuthRequest(c)
		return
	}

	dbRefreshToken := models.SiteRefreshToken{}
	result := all.Postgres.Where("token = ?", body.RefreshToken).First(&dbRefreshToken)
	if result.Error != nil {
		common.ErrorInvalidOAuthRequest(c)
		return
	}

	if dbRefreshToken.ID == 0 {
		common.ErrorInvalidOAuthRequest(c)
		return
	}

	user, err := common.GetUserByAccountId(dbRefreshToken.AccountId)
	if err != nil {
		common.ErrorInvalidOAuthRequest(c)
		return
	}

	user.Password = ""
	token := GenerateSiteToken(user, "site")
	c.JSON(http.StatusOK, gin.H{"data": user, "token": token})
}

func UserUpdate(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBind(&body); err != nil {
		common.ErrorBadRequest(c)
		return
	}

	if body.Username != "" {
		user.Username = OnlyAllowCharacters(body.Username)

		var dbUser models.User
		usernameCheckResult := all.Postgres.Where("username = ?", user.Username).First(&dbUser)
		if usernameCheckResult.Error == nil && user.ID != dbUser.ID {
			common.ErrorNameTaken(c)
			return
		}
	}
	if body.Password != "" && body.Password != user.Password {
		user.Password = all.HashString(body.Password)
	}

	result := all.Postgres.Save(&user)
	if result.Error != nil {
		common.ErrorInternalServer(c)
		return
	}

	user.Password = ""
	token := GenerateSiteToken(user, "site")
	c.JSON(http.StatusOK, gin.H{"data": user, "token": token})
}

func OnlyAllowCharacters(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		all.PrintGreen([]any{"serre", err})
		return ""
	}
	return reg.ReplaceAllString(s, "")
}

func AdminGetProfile(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 1 {
		common.ErrorUnauthorized(c)
		return
	} 

	accountId := c.Param("accountId")
	profileId := c.Param("profileId")

	user := models.User{}
	result := all.Postgres.Where("account_id = ?", accountId).First(&user)
	if result.Error != nil {
		common.ErrorBadRequest(c)
		return
	}

	profile, err := common.ReadProfileFromUser(accountId, profileId)
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	athenaProfile, err := common.ReadProfileFromUser(accountId, "athena")
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	commonCoreProfile, err := common.ReadProfileFromUser(accountId, "common_core")
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"profile": profile,
		"athenaProfile": athenaProfile,
		"commonCoreProfile": commonCoreProfile,
	})
}

func AdminSaveProfile(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 1 {
		common.ErrorUnauthorized(c)
		return
	}
	
	accountId := c.Param("accountId")

	var body struct {
		Profile models.Profile `json:"profile" binding:"required"`
		AthenaProfile models.AthenaProfile `json:"athenaProfile" binding:"required"`
		CommonCoreProfile models.CommonCoreProfile `json:"commonCoreProfile" binding:"required"`
		User models.User `json:"user" binding:"required"`
	}

	if err := c.ShouldBind(&body); err != nil {
		all.PrintRed([]any{"serre", err})
		common.ErrorBadRequest(c)
		return
	}

	var user models.User
	result := all.Postgres.Where("account_id = ?", accountId).First(&user)
	if result.Error != nil {
		all.PrintRed([]any{"cannot find user ", accountId})
		common.ErrorBadRequest(c)
		return
	}

	athenaProfile, err := common.ReadProfileFromUser(accountId, "athena")
	if err != nil {
		all.PrintRed([]any{"cannot find athena profile", accountId})
		common.ErrorBadRequest(c)
		return
	}

	commonCoreProfile, err := common.ReadProfileFromUser(accountId, "common_core")
	if err != nil {
		all.PrintRed([]any{"cannot find common core profile", accountId})
		common.ErrorBadRequest(c)
		return
	}

	athenaProfileConverted, err := common.ConvertProfileToAthena(athenaProfile)
	if err != nil {
		common.ErrorInternalServer(c)
		return
	}

	commonCoreProfileConverted, err := common.ConvertProfileToCommonCore(commonCoreProfile)
	if err != nil {
		common.ErrorInternalServer(c)
		return
	}

	athenaProfileConverted.Items = body.AthenaProfile.Items
	commonCoreProfileConverted.Items = body.CommonCoreProfile.Items

	athenaProfileConverted.Stats.Attributes = body.AthenaProfile.Stats.Attributes
	commonCoreProfileConverted.Stats.Attributes = body.CommonCoreProfile.Stats.Attributes

	defaultAthenaProfile, err := common.ConvertAthenaToDefault(athenaProfileConverted)
	if err != nil {
		common.ErrorInternalServer(c)
		return
	}

	defaultCommonCoreProfile, err := common.ConvertCommonCoreToDefault(commonCoreProfileConverted)
	if err != nil {
		common.ErrorInternalServer(c)
		return
	}

	common.AppendLoadoutsToProfile(&defaultAthenaProfile, user.AccountId)
	common.AppendLoadoutsToProfile(&defaultCommonCoreProfile, user.AccountId)

	commonCore, err := common.ReadProfileFromUser(accountId, "common_core")
	if err != nil {
		return
	}

	gift := models.CommonCoreItem{
		TemplateId: "GiftBox:gb_default",
		Attributes: gin.H{
			"fromAccountId": "Server",
			"lootList": []gin.H{
				
			},
			"params": gin.H{
				"userMessage": "Server has updated your account. Enjoy!",
			},
			"level": 1,
			"giftedOn": time.Now().Format("2006-01-02T15:04:05.999Z"),
		},
		Quantity: 1,
	}
	commonCore.Items["GiftBox:gb_default"] = gift
	common.SaveProfileToUser(accountId, commonCore)

	if body.User.VBucks != 0 {
		user.VBucks = body.User.VBucks
		gift.Attributes["lootList"] = append(gift.Attributes["lootList"].([]gin.H), gin.H{
			"itemType": "MtxCurrency:MTXCurrency",
			"itemGuid": "MtxCurrency:MTXCurrency",
			"itemProfile": "athena",
			"quantity": user.VBucks - body.User.VBucks,
		})
	}
	if user.Username != "admin" && body.User.AccessLevel != 0 {
		user.AccessLevel = body.User.AccessLevel
	}
	user.Banned = body.User.Banned
	all.Postgres.Save(&user)

	socket.XMPPSendBodyToAccountId(gin.H{
		"payload": gin.H{},
		"type": "com.epicgames.gift.received",
		"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
	}, accountId)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"athenaProfile": defaultAthenaProfile,
		"commonCoreProfile": defaultCommonCoreProfile,
	})
}

func AdminGetAllUsers(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 1 {
		common.ErrorUnauthorized(c)
		return
	}

	var users []models.User
	result := all.Postgres.Find(&users)
	if result.Error != nil {
		common.ErrorInternalServer(c)
		return
	}

	c.JSON(http.StatusOK, users)
}

func AdminGiveAllSkins(c * gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 1 {
		common.ErrorUnauthorized(c)
		return
	}
	
	accountId := c.Param("accountId")

	profile, err := common.ReadProfileFromUser(accountId, "athena")
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	common.AddEverythingToProfile(&profile, accountId)
	common.AppendLoadoutsToProfile(&profile, accountId)
	
	socket.XMPPSendBodyToAccountId(gin.H{
		"payload": gin.H{},
		"type": "com.epicgames.gift.received",
		"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
	}, accountId)

	c.JSON(http.StatusOK, profile)
}

func AdminGiveItem(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 1 {
		common.ErrorUnauthorized(c)
		return
	}

	accountId := c.Param("accountId")
	itemId := c.Param("itemId")

	profile, err := common.ReadProfileFromUser(accountId, "athena")
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	commonCore, err := common.ReadProfileFromUser(accountId, "common_core")
	if err != nil {
		return
	}

	gift := models.CommonCoreItem{
		TemplateId: "GiftBox:gb_default",
		Attributes: gin.H{
			"fromAccountId": "Server",
			"lootList": []gin.H{
				{
					"itemType": itemId,
					"itemGuid": itemId,
					"itemProfile": "athena",
					"quantity": 1,
				},
			},
			"params": gin.H{
				"userMessage": "Enjoy this gift from the server!",
			},
			"level": 1,
			"giftedOn": time.Now().Format("2006-01-02T15:04:05.999Z"),
		},
		Quantity: 1,
	}
	commonCore.Items["GiftBox:gb_default"] = gift
	common.SaveProfileToUser(accountId, commonCore)

	common.AddItemToProfile(&profile, itemId, accountId)
	common.AppendLoadoutsToProfile(&profile, accountId)

	socket.XMPPSendBodyToAccountId(gin.H{
		"payload": gin.H{},
		"type": "com.epicgames.gift.received",
		"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
	}, accountId)

	c.JSON(http.StatusOK, profile)
}

func AdminTakeAllSkins(c * gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 1 {
		common.ErrorUnauthorized(c)
		return
	}
	
	accountId := c.Param("accountId")

	profile, err := common.ReadProfileFromUser(accountId, "athena")
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	common.RemoveEverythingFromProfile(&profile, accountId)
	common.AddItemsToProfile(&profile, []string{
		"AthenaCharacter:CID_001_Athena_Commando_F_Default",
		"AthenaPickaxe:DefaultPickaxe",
		"AthenaGlider:DefaultGlider",
		"AthenaDance:EID_DanceMoves",
	}, accountId)
	common.SaveProfileToUser(accountId, profile)

	socket.XMPPSendBodyToAccountId(gin.H{
		"payload": gin.H{},
		"type": "com.epicgames.gift.received",
		"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
	}, accountId)

	c.JSON(http.StatusOK, profile)
}

func AdminTakeItem(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 1 {
		common.ErrorUnauthorized(c)
		return
	}

	accountId := c.Param("accountId")
	itemId := c.Param("itemId")

	profile, err := common.ReadProfileFromUser(accountId, "athena")
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	common.RemoveItemFromProfile(&profile, itemId, accountId)
	common.AppendLoadoutsToProfile(&profile, accountId)

	socket.XMPPSendBodyToAccountId(gin.H{
		"payload": gin.H{},
		"type": "com.epicgames.gift.received",
		"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
	}, accountId)

	c.JSON(http.StatusOK, profile)
}

func AdminGetLocker(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 1 {
		common.ErrorUnauthorized(c)
		return
	}

	accountId := c.Param("accountId")

	profile, err := common.ReadProfileFromUser(accountId, "athena")
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	marshalledItems, err := json.Marshal(profile.Items)
	if err != nil {
		all.PrintGreen([]any{"serre", err})
		common.ErrorBadRequest(c)
		return
	}

	var items map[string]models.Item
	err = json.Unmarshal(marshalledItems, &items)
	if err != nil {
		all.PrintGreen([]any{"serre2222", err})
		common.ErrorBadRequest(c)
		return
	}

	if common.AllItems == nil {
		common.GetAllFortniteItems()
	}

	locker := []LockerItem{}
	for _, item := range items {
		if item.TemplateId == "Currency:MtxPurchased" || item.TemplateId == "CosmeticLocker:cosmeticlocker_athena" {
			continue
		}

		if common.AllItemsKeys[item.TemplateId].IntroductionSeason > 20 {
			continue
		}

		locker = append(locker, LockerItem{
			ItemId: item.TemplateId,
			Rarity: common.AllItemsKeys[item.TemplateId].Rarity,
			Season: common.AllItemsKeys[item.TemplateId].IntroductionSeason,
		})
	}

	c.JSON(http.StatusOK, locker)
}

func AdminGiveUserAdmin(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 2 {
		common.ErrorUnauthorized(c)
		return
	}

	accountId := c.Param("accountId")

	var user models.User
	result := all.Postgres.Where("account_id = ?", accountId).First(&user)
	if result.Error != nil {
		common.ErrorBadRequest(c)
		return
	}

	user.AccessLevel = 1
	all.Postgres.Save(&user)

	c.JSON(http.StatusOK, user)
}

func AdminTakeUserAdmin(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 2 {
		common.ErrorUnauthorized(c)
		return
	}

	accountId := c.Param("accountId")

	var user models.User
	result := all.Postgres.Where("account_id = ?", accountId).First(&user)
	if result.Error != nil {
		common.ErrorBadRequest(c)
		return
	}

	user.AccessLevel = 0
	all.Postgres.Save(&user)

	c.JSON(http.StatusOK, user)
}

func GetFriendlyShop(c *gin.Context) {
	realShop := ItemShop

	fmt.Println(len(realShop.Storefronts))

	if len(realShop.Storefronts) == 0 {
		GenerateRandomItemShop()
	}

	daily := []models.SiteShopItem{}
	for _, item := range realShop.Storefronts[0].CatalogEntries {
		simpleItem, ok := common.AllItemsKeys[item.ItemGrants[0].TemplateID]
		if !ok {
			continue
		}

		daily = append(daily, models.SiteShopItem{
			ItemId: item.ItemGrants[0].TemplateID,
			Price: item.Prices[0].FinalPrice,
			Rarity: simpleItem.Rarity,
			Season: simpleItem.IntroductionSeason,
			Name: simpleItem.Name,
		})
	}

	featured := []models.SiteShopItem{}
	for _, item := range realShop.Storefronts[1].CatalogEntries {
		simpleItem, ok := common.AllItemsKeys[item.ItemGrants[0].TemplateID]
		if !ok {
			continue
		}

		featured = append(featured, models.SiteShopItem{
			ItemId: item.ItemGrants[0].TemplateID,
			Price: item.Prices[0].FinalPrice,
			Rarity: simpleItem.Rarity,
			Season: simpleItem.IntroductionSeason,
			Name: simpleItem.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"daily": daily,
		"featured": featured,
	})
}

func AdminChangeShop(c *gin.Context) {
	me := c.MustGet("user").(models.User)
	if me.AccessLevel < 2 {
		common.ErrorUnauthorized(c)
		return
	}
	
	GenerateRandomItemShop()
	GetFriendlyShop(c)
}