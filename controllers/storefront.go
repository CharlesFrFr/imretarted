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
	"github.com/zombman/server/all"
	"github.com/zombman/server/common"
	"github.com/zombman/server/models"
)

var Prices = models.Prices{
	"AthenaCharacter": {
		"Legendary": 2000,
		"Epic":      1500,
		"Rare":      1200,
		"Uncommon":  800,
		"Common":    500,
	},
	"AthenaPickaxe": {
		"Legendary": 1500,
		"Epic":      1200,
		"Rare":      800,
		"Uncommon":  500,
		"Common":    300,
	},
	"AthenaGlider": {
		"Legendary": 1500,
		"Epic":      1200,
		"Rare":      800,
		"Uncommon":  500,
		"Common":    300,
	},
	"AthenaDance": {
		"Legendary": 800,
		"Epic":      800,
		"Rare":      500,
		"Uncommon":  200,
		"Common":    200,
	},
	"AthenaSkyDiveContrail": {
		"Legendary": 500,
		"Epic":      400,
		"Rare":      300,
		"Uncommon":  200,
		"Common":    100,
	},
}

var ItemShop models.StorePage
var ShouldRefresh int64 = 0

func StorefrontCatalog(c *gin.Context) {
	timeNow := time.Now().Unix()

	if timeNow > ShouldRefresh {
		RefreshItemShop()	
	}
	
	if len(ItemShop.Storefronts) == 0 {
		GenerateRandomItemShop()
	}

	if (common.LoadShopFromJson) {
		pathToProfile := "data/shop.json"

		file, err := os.Open(pathToProfile)
		if err != nil {
			all.PrintRed([]any{"error opening shop", err})
			return
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			all.PrintRed([]any{"error reading shop", err})
			return
		}
		str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

		err = json.Unmarshal([]byte(str), &ItemShop)
		if err != nil {
			all.PrintRed([]any{"error unmarshalling shop", err})
			return
		}

		ItemShop.Expiration = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Now().Location()).Format("2006-01-02T15:04:05.999Z")
	}

	c.JSON(http.StatusOK, ItemShop)
}

func RefreshItemShop() {
	ShouldRefresh = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Now().Location()).Unix()
	GenerateRandomItemShop()
}

func GenerateRandomItemShop() {
	endOfDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Now().Location()).Format("2006-01-02T15:04:05.999Z")
	allItems, err := common.GetItemsFromSeason(common.Season + (common.Chapter * 10))
	if err != nil {
		return
	}

	common.ExcludeId(&allItems, "EID_AirHornRaisin")

	legendaryItems := common.FilterRarity(allItems, "Legendary")
	dailyItems := common.ExcludeRarity(allItems, "Legendary")

	ItemShop = models.StorePage{
		RefreshIntervalHrs: 24,
		DailyPurchaseHrs:  24,
		Expiration: endOfDay,
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
		entry := GenerateRandomCatalogEntry(1, &dailyItems)
		ItemShop.Storefronts[0].CatalogEntries = append(ItemShop.Storefronts[0].CatalogEntries, entry)
	}

	for i := 0; i < 2; i++ {
		entry := GenerateRandomCatalogEntry(-1, &legendaryItems)
		ItemShop.Storefronts[1].CatalogEntries = append(ItemShop.Storefronts[1].CatalogEntries, entry)
	}

	all.PrintGreen([]any{"generated new random item shop"})
	if !common.LoadShopFromJson {
		SaveItemShop()
	}
}

func GenerateRandomCatalogEntry(f int, items *[]models.BeforeStoreItem) models.CatalogEntry {
	endOfDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Now().Location()).Format("2006-01-02T15:04:05.999Z")
	randomItem := (*items)[rand.Intn(len(*items) - 1)]
	price := Prices[randomItem.BackendType][randomItem.Rarity]
	id := all.HashString(randomItem.ID)
	id = id[:40]

	return models.CatalogEntry{
		DevName: id,
		OfferID: id,
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
			SaleExpiration: endOfDay,
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
			RequiredID: randomItem.BackendType + ":" + randomItem.ID,
			MinQuantity: 1,
		}},
		ItemGrants: []models.ItemGrant{{
			TemplateID: randomItem.BackendType + ":" + randomItem.ID,
			Quantity: 1,
		}},
	}
}

func SaveItemShop() {
	file, err := os.OpenFile("data/shop.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := json.MarshalIndent(ItemShop, "", "\t")
	if err != nil {
		return
	}

	file.Write(data)
}