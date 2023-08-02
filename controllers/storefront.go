package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/common"
	"github.com/zombman/server/helpers"
	"github.com/zombman/server/models"
)

var randomPrices = []int{0, 200, 500, 600, 800, 1200, 1500, 2000}
var ItemShop models.StorePage
var ShouldRefresh int64

func StorefrontCatalog(c *gin.Context) {
	timeNow := time.Now().Unix()

	if timeNow > ShouldRefresh {
		RefreshItemShop()	
	}

	if len(ItemShop.Storefronts) == 0 {
		GenerateRandomItemShop()
	}

	c.JSON(http.StatusOK, ItemShop)
}

func StorefrontFromJSONFile(c *gin.Context) {
	file, err := os.Open("default/shop.json")
	if err != nil {
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return
	}
	str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

	var shop models.StorePage
	err = json.Unmarshal([]byte(str), &shop)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, shop)
}

func RefreshItemShop() {
	ShouldRefresh = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Now().Location()).Unix()
	GenerateRandomItemShop()
}

func GenerateRandomItemShop() {
	ItemShop = models.StorePage{
		RefreshIntervalHrs: 24,
		DailyPurchaseHrs:  24,
		Expiration: time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		Storefronts: []models.Storefront{
			{
				Name: "BRDailyStorefront",
				CatalogEntries: []models.CatalogEntry{},
			},
			{
				Name: "BRWeeklyStorefront",
				CatalogEntries: []models.CatalogEntry{},
			},
			{
				Name: "BRSeasonStorefront",
				CatalogEntries: []models.CatalogEntry{},
			},
		},
	}

	for i := 0; i < 6; i++ {
		entry := GenerateRandomCatalogEntry(1)
		ItemShop.Storefronts[0].CatalogEntries = append(ItemShop.Storefronts[0].CatalogEntries, entry)
	}

	for i := 0; i < 2; i++ {
		entry := GenerateRandomCatalogEntry(-1)
		ItemShop.Storefronts[1].CatalogEntries = append(ItemShop.Storefronts[1].CatalogEntries, entry)
	}
}

func GenerateRandomCatalogEntry(f int) models.CatalogEntry {
	price := randomPrices[rand.Intn(len(randomPrices) - 1) + 1]
	randomId := "AthenaCharacter:" + common.SkinList[rand.Intn(len(common.SkinList) -1 ) + 1]

	return models.CatalogEntry{
		DevName: helpers.HashStringSHA1(randomId),
		OfferID: helpers.HashStringSHA1(randomId),
		FulfillmentIDs: []string{},
		DailyLimit: -1,
		WeeklyLimit: -1,
		MonthlyLimit: -1,
		Categories: []string{},
		Prices: []models.Price{{
			CurrencyType: "MtxCurrency",
			CurrencySubType: "CurrencySource",
			RegularPrice: price,
			FinalPrice: price,
			SaleExpiration: time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
			BasePrice: price,
		}},
		MatchFilter: "",
		AppStoreID: []string{},
		FilterWeight: f,
		SortPriority: f,
		CatalogGroupPriority: 0,
		Refundable: false,
		DisplayAssetPath: "",
		OfferType: "StaticPrice",
		GiftInfo: map[string]any {
			"bIsEnabled": false,
			"forcedGiftBoxTemplateId": "",
			"purchaseRequirements": []any{},
			"giftRecordIds": []any{},
		},
		Meta: map[string]any {
			"SectionId": "Featured",
			"TileSize": "Small",
		},
		MetaInfo: []gin.H{{
			"key": "SectionId",
			"value": "Featured",
		}, {
			"key": "TileSize",
			"value": "Small",
		}},
		Requirements: []models.Requirement{{
			RequirementType: "DenyOnItemOwnership",
			RequiredID: randomId,
			MinQuantity: 1,
		}},
		ItemGrants: []models.ItemGrant{{
			TemplateID: randomId,
			Quantity: 1,
		}},
	}
}