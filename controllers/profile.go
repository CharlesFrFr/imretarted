package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

func ProfileActionHandler(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	profileId, _ := c.GetQuery("profileId")
	action := c.Param("action")

	response := models.ProfileResponse{}

	profile, err := common.ReadProfileFromUser(user.AccountId, profileId)
	if err != nil {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	switch action {
		case "QueryProfile":
		case "SetMtxPlatform":
		case "ClientQuestLogin":
		case "BulkEquipBattleRoyaleCustomization":
		case "MarkItemSeen":
		case "SetItemFavoriteStatusBatch":
			break
		case "EquipBattleRoyaleCustomization":
			EquipBattleRoyaleCustomization(c, user, &profile, &response)
		case "PurchaseCatalogEntry":
			PurchaseCatalogEntry(c, user, &profile, &response)
		case "SetBattleRoyaleBanner":
			SetBattleRoyaleBanner(c, user, &profile, &response)
		default:
			all.PrintRed([]any{action, "Not Handled"})
			return
	}

	if profile.ProfileId == "athena" {
		profile.Stats.Attributes["season_num"] = common.Season
	}
	
	if queryRevision, err := strconv.Atoi(c.Query("rvn")); err == nil && queryRevision != profile.Rvn {
		response.ProfileChanges = []models.ProfileChange{{
			ChangeType: "fullProfileUpdate",
			Profile: profile,
		}}
	}

	profile.Rvn += 1
	profile.CommandRevision = profile.Rvn
	profile.AccountId = user.AccountId
	profile.Updated = time.Now().Format("2006-01-02T15:04:05.999Z")

	common.SaveProfileToUser(user.AccountId, profile)

	response.ProfileRevision = profile.Rvn
	response.ProfileCommandRevision = profile.CommandRevision
	response.ProfileID = profileId
	response.ProfileChangesBaseRevision = profile.Rvn - 1
	response.ServerTime = time.Now().Format("2006-01-02T15:04:05.999Z")
	response.ResponseVersion = 1

	c.JSON(200, response)
}

func DedicatedServerProfileHandler(c *gin.Context) {
	userId := c.Param("accountId")
	profileId, _ := c.GetQuery("profileId")

	response := models.ProfileResponse{}

	profile, err := common.ReadProfileFromUser(userId, profileId)
	if err != nil {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	if profile.ProfileId == "athena" {
		profile.Stats.Attributes["season_num"] = common.Season
	}
	
	response.ProfileChanges = []models.ProfileChange{{
		ChangeType: "fullProfileUpdate",
		Profile: profile,
	}}

	profile.Rvn += 1
	profile.CommandRevision = profile.Rvn
	profile.AccountId = userId
	profile.Updated = time.Now().Format("2006-01-02T15:04:05.999Z")

	response.ProfileRevision = profile.Rvn
	response.ProfileCommandRevision = profile.CommandRevision
	response.ProfileID = profileId
	response.ProfileChangesBaseRevision = profile.Rvn - 1
	response.ServerTime = time.Now().Format("2006-01-02T15:04:05.999Z")
	response.ResponseVersion = 1

	c.JSON(200, response)
}

func PurchaseCatalogEntry(c *gin.Context, user models.User, profile *models.Profile, response *models.ProfileResponse) {
	athenaProfile, nerr := common.ReadProfileFromUser(user.AccountId, "athena")
	if nerr != nil {
		common.ErrorBadRequest(c)
		return
	}

	var body struct {
		OfferId string `json:"offerId"`
		PurchaseQuantity int `json:"purchaseQuantity"`
		Currency string `json:"currency"`
		CurrencySubType string `json:"currencySubType"`
		ExpectedTotalPrice int `json:"expectedTotalPrice"`
		GameContext string `json:"gameContext"`
	}

	if err := c.ShouldBind(&body); err != nil {
		all.PrintRed([]any{"could not bind body", err.Error()})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	if strings.Contains(body.OfferId, "battlepass") {
		PurchaseCatalogEntryBattlePass(c, user, profile, response, body)
		return
	}

	offer, err := common.GetCatalogEntry(body.OfferId)
	if err != nil {
		all.PrintRed([]any{"could not find offer", body.OfferId})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	playerHasItem := false
	for _, item := range profile.Items {
		var marshItem models.Item
		marshalItem, err := json.Marshal(item)
		if err != nil {
			all.PrintRed([]any{"could not marshal item", item})
			common.ErrorBadRequest(c)
			c.Abort()
			return
		}
		err = json.Unmarshal(marshalItem, &marshItem)
		if err != nil {
			all.PrintRed([]any{"could not unmarshal item", item})
			common.ErrorBadRequest(c)
			c.Abort()
			return
		}

		if marshItem.TemplateId == offer.ItemGrants[0].TemplateID {
			playerHasItem = true
			break
		}
	}

	if playerHasItem {
		all.PrintRed([]any{"player already has item", offer.ItemGrants[0].TemplateID})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	if offer.Prices[0].FinalPrice != body.ExpectedTotalPrice {
		all.PrintRed([]any{"expected price does not match", offer.Prices[0].FinalPrice, body.ExpectedTotalPrice})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	if user.VBucks < offer.Prices[0].FinalPrice {
		all.PrintRed([]any{"player does not have enough vbucks", user.VBucks, offer.Prices[0].FinalPrice})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	common.AddItemToProfile(&athenaProfile, offer.ItemGrants[0].TemplateID, user.AccountId)
	common.TakeUserVBucks(user.AccountId, profile, offer.Prices[0].FinalPrice)

	athenaProfile.Stats.Attributes["season_num"] = common.Season
	athenaProfile.Rvn += 1
	athenaProfile.CommandRevision = athenaProfile.Rvn
	athenaProfile.AccountId = user.AccountId
	athenaProfile.Updated = time.Now().Format("2006-01-02T15:04:05.999Z")
	common.SaveProfileToUser(user.AccountId, athenaProfile)

	response.MultiUpdate = append(response.MultiUpdate, models.MultiUpdate{
		ProfileRevision: athenaProfile.Rvn,
		ProfileCommandRevision: athenaProfile.CommandRevision,
		ProfileID: "athena",
		ProfileChangesBaseRevision: athenaProfile.Rvn - 1,
		ProfileChanges: []models.ProfileChange{{
			ChangeType: "itemAdded",
			ItemID: offer.ItemGrants[0].TemplateID,
			Item: models.Item{
				TemplateId: offer.ItemGrants[0].TemplateID,
				Attributes: models.ItemAttributes{
					ItemSeen: false,
					Variants: []models.ItemVariant{},
				},
				Quantity: 1,
			},
		}},
	})

	response.Notifications = append(response.Notifications, models.Notification{
		Type: "CatalogPurchase",
		Primary: true,
		LootResult: models.LootResult{ Items: []models.LootResultItem{{
			ItemType: offer.ItemGrants[0].TemplateID,
			ItemGuid: offer.ItemGrants[0].TemplateID,
			ItemProfile: "athena",
			Quantity: 1,
		}}},
	})

	response.ProfileChanges = append(response.ProfileChanges, models.ProfileChange{
		ChangeType: "itemQuantityChanged",
		ItemID: "Currency:MtxPurchased",
		Quantity: user.VBucks - offer.Prices[0].FinalPrice,
	})
}

func EquipBattleRoyaleCustomization(c *gin.Context, user models.User, profile *models.Profile, response *models.ProfileResponse) {
	if profile.ProfileId != "athena" {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	var body struct {
		SlotName string `json:"slotName"`
		ItemToSlot string `json:"itemToSlot"`
		IndexWithinSlot int `json:"indexWithinSlot"`
		VariantUpdates []models.ItemVariant `json:"variantUpdates"`
	}

	if err := c.ShouldBind(&body); err != nil {
		all.PrintRed([]any{"could not bind body", err.Error()})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	athenaProfile,err := common.ConvertProfileToAthena(*profile)
	if err != nil {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	activeLoadoutId := athenaProfile.Stats.Attributes.Loadouts[athenaProfile.Stats.Attributes.ActiveLoadoutIndex]
	activeLoadout, err := common.GetLoadout(activeLoadoutId, user.AccountId)
	if err != nil {
		common.ErrorItemNotFound(c)
		c.Abort()
		return
	}

	lowercaseItemType := strings.ToLower(body.SlotName)
	var valueChanged any

	switch lowercaseItemType {
		case "character":
			athenaProfile.Stats.Attributes.FavoriteCharacter = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Character"].Items[0] = body.ItemToSlot
			valueChanged = athenaProfile.Stats.Attributes.FavoriteCharacter
		case "backpack":
			athenaProfile.Stats.Attributes.FavoriteBackpack = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Backpack"].Items[0] = body.ItemToSlot
			valueChanged = athenaProfile.Stats.Attributes.FavoriteBackpack
		case "pickaxe":
			athenaProfile.Stats.Attributes.FavoritePickaxe = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Pickaxe"].Items[0] = body.ItemToSlot
			valueChanged = athenaProfile.Stats.Attributes.FavoritePickaxe
		case "glider":
			athenaProfile.Stats.Attributes.FavoriteGlider = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Glider"].Items[0] = body.ItemToSlot
			valueChanged = athenaProfile.Stats.Attributes.FavoriteGlider
		case "skydivecontrail":
			athenaProfile.Stats.Attributes.FavoriteSkyDiveContrail = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["SkyDiveContrail"].Items[0] = body.ItemToSlot
			valueChanged = athenaProfile.Stats.Attributes.FavoriteSkyDiveContrail
		case "loadingscreen":
			athenaProfile.Stats.Attributes.FavoriteLoadingScreen = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["LoadingScreen"].Items[0] = body.ItemToSlot
			valueChanged = athenaProfile.Stats.Attributes.FavoriteLoadingScreen
		case "musicpack":
			athenaProfile.Stats.Attributes.FavoriteMusicPack = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["MusicPack"].Items[0] = body.ItemToSlot
			valueChanged = athenaProfile.Stats.Attributes.FavoriteMusicPack
		case "dance":
			athenaProfile.Stats.Attributes.FavoriteDance[body.IndexWithinSlot] = body.ItemToSlot
			activeLoadout.Attributes.LockerSlotsData.Slots["Dance"].Items[body.IndexWithinSlot] = body.ItemToSlot
			valueChanged = athenaProfile.Stats.Attributes.FavoriteDance
		case "itemwrap":
			if body.IndexWithinSlot >= 0 {
				athenaProfile.Stats.Attributes.FavoriteItemWraps[body.IndexWithinSlot] = body.ItemToSlot
				activeLoadout.Attributes.LockerSlotsData.Slots["ItemWrap"].Items[body.IndexWithinSlot] = body.ItemToSlot
			} 
			if body.IndexWithinSlot == -1 {
				for i := range activeLoadout.Attributes.LockerSlotsData.Slots["ItemWrap"].Items {
					activeLoadout.Attributes.LockerSlotsData.Slots["ItemWrap"].Items[i] = body.ItemToSlot
				}
				for i := range athenaProfile.Stats.Attributes.FavoriteItemWraps {
					athenaProfile.Stats.Attributes.FavoriteItemWraps[i] = body.ItemToSlot
				}
			}

			valueChanged = athenaProfile.Stats.Attributes.FavoriteItemWraps
			lowercaseItemType = "itemwraps"
		default:
			all.PrintRed([]any{"unknown item type", athenaProfile.Stats.Attributes})
			common.ErrorBadRequest(c)
			c.Abort()
	}

	
	defaultProfile, err := common.ConvertAthenaToDefault(athenaProfile)
	if err != nil {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	profile.Items = defaultProfile.Items
	profile.Stats = defaultProfile.Stats

	response.ProfileChanges = append(response.ProfileChanges, models.ProfileChange{
		ChangeType: "statModified",
		Name: "favorite_" + lowercaseItemType,
		Value: valueChanged,
	})

	for _, variant := range body.VariantUpdates {
		itemWithVariant, err := common.GetItemFromProfile(profile, body.ItemToSlot)
		if err != nil {
			all.PrintRed([]any{"could not find item", body.ItemToSlot})
			common.ErrorBadRequest(c)
			return
		}

		variantFound, err := common.FindVariant(&itemWithVariant, variant.Channel)
		if err != nil {
			variantFound = models.ItemVariant{
				Channel: variant.Channel,
				Active: variant.Active,
				Owned: variant.Owned,
			}
		}

		variantFound.Active = variant.Active
		variantFound.Owned = []string{variant.Active}
		common.SetVariantInItem(&itemWithVariant, variantFound)
		athenaProfile.Items[body.ItemToSlot] = itemWithVariant

		response.ProfileChanges = append(response.ProfileChanges, models.ProfileChange{
			ChangeType: "itemAttrChanged",
			ItemID: body.ItemToSlot,
			AttributeName: "variants",
			AttributeValue: itemWithVariant.Attributes.Variants,
		})
	}

	athenaProfile.Stats.Attributes.LastAppliedLoadout = activeLoadoutId

	common.AppendLoadoutToProfileNoSave(profile, &activeLoadout, user.AccountId)
}

func SetBattleRoyaleBanner(c *gin.Context, user models.User, profile *models.Profile, response *models.ProfileResponse) {
	var body struct {
		HomebaseBannerIconId string `json:"homebaseBannerIconId"`
		HomebaseBannerColorId string `json:"homebaseBannerColorId"`
	}

	if err := c.ShouldBind(&body); err != nil {
		all.PrintRed([]any{"could not bind body", err.Error()})
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	athenaProfile,err := common.ConvertProfileToAthena(*profile)
	if err != nil {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	athenaProfile.Stats.Attributes.BannerIcon = body.HomebaseBannerIconId
	athenaProfile.Stats.Attributes.BannerColor = body.HomebaseBannerColorId

	activeLoadoutId := athenaProfile.Stats.Attributes.Loadouts[athenaProfile.Stats.Attributes.ActiveLoadoutIndex]
	activeLoadout, err := common.GetLoadout(activeLoadoutId, user.AccountId)
	if err != nil {
		common.ErrorItemNotFound(c)
		c.Abort()
		return
	}

	activeLoadout.Attributes.BannerIconTemplate = body.HomebaseBannerIconId
	activeLoadout.Attributes.BannerColorTemplate = body.HomebaseBannerColorId

	defaultProfile, err := common.ConvertAthenaToDefault(athenaProfile)
	if err != nil {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	profile.Items = defaultProfile.Items
	profile.Stats = defaultProfile.Stats

	common.AppendLoadoutToProfileNoSave(profile, &activeLoadout, user.AccountId)

	response.ProfileChanges = append(response.ProfileChanges, models.ProfileChange{
		ChangeType: "statModified",
		Name: "banner_icon",
		Value: body.HomebaseBannerIconId,
	})

	response.ProfileChanges = append(response.ProfileChanges, models.ProfileChange{
		ChangeType: "statModified",
		Name: "banner_color",
		Value: body.HomebaseBannerColorId,
	})
}

func PurchaseCatalogEntryBattlePass(c *gin.Context, user models.User, profile *models.Profile, response *models.ProfileResponse, body struct {
	OfferId string `json:"offerId"`
	PurchaseQuantity int `json:"purchaseQuantity"`
	Currency string `json:"currency"`
	CurrencySubType string `json:"currencySubType"`
	ExpectedTotalPrice int `json:"expectedTotalPrice"`
	GameContext string `json:"gameContext"`
}) {
	offer, err := common.GetCatalogEntry(body.OfferId)
	if err != nil {
		common.ErrorItemNotFound(c)
		c.Abort()
		return
	}

	all.PrintGreen([]any{"purchasing battle pass"})
	all.MarshPrintJSON(offer)

	if user.VBucks < offer.Prices[0].FinalPrice {
		common.ErrorBadRequest(c)
		c.Abort()
		return
	}

	response.ProfileChanges = append(response.ProfileChanges, models.ProfileChange{
		ChangeType: "itemQuantityChanged",
		ItemID: "Currency:MtxPurchased",
		Quantity: user.VBucks - offer.Prices[0].FinalPrice,
	})
}