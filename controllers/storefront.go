package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		pathToShop := "data/shop/shop.json"

		file, err := os.Open(pathToShop)
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
	allItems, err := common.GetItemsFromSeason(common.Season)
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
			{
				Name: "BRSeason8",
				CatalogEntries: []models.CatalogEntry{
					{
						OfferID: "battlepassbundle" + fmt.Sprint(common.Season),
						DevName: "BR.Season"+ fmt.Sprint(common.Season) +".BattleBundle.01",
						OfferType: "StaticPrice",
						Prices: []models.Price{{
							CurrencyType: "MtxCurrency",
							CurrencySubType: "",
							RegularPrice: 4700,
							FinalPrice: 2800,
							SaleExpiration: "9999-12-31T23:59:59.999Z",
							BasePrice: 2800,
						}},
						Categories: []string{},
						DailyLimit: -1,
						WeeklyLimit: -1,
						MonthlyLimit: -1,
						AppStoreID: []string{"", "", "", "", "", "", "", "", "", ""},
						Requirements: []models.Requirement{{
							RequirementType: "DenyOnFulfillment",
							RequiredID: "battlepass"+ fmt.Sprint(common.Season),
							MinQuantity: 1,
						}},
						MetaInfo: []gin.H{},
						CatalogGroupPriority: 0,
						SortPriority: 0,
						DisplayAssetPath: "/Game/Catalog/DisplayAssets/DA_BR_Season"+ fmt.Sprint(common.Season) +"_BattlePassWithLevels.DA_BR_Season"+ fmt.Sprint(common.Season) +"_BattlePassWithLevels",
						ItemGrants: []models.ItemGrant{},
						Refundable: false,
						Title: "Battle Bundle",
						ShortDescription: "Battle Pass + 25 tiers!",
						Description: "Season 8 \n\nInstantly get these items <Bold>valued at over 10,000 V-Bucks</>.\n  • <ItemName>Blackheart</> Progressive Outfit\n  • <ItemName>Hybrid</> Progressive Outfit\n  • <ItemName>Sidewinder</> Outfit\n  • <ItemName>Tropical Camo</> Wrap\n  • <ItemName>Woodsy</> Pet\n  • <ItemName>Sky Serpents</> Glider\n  • <ItemName>Cobra</> Back Bling\n  • <ItemName>Flying Standard</> Contrail\n  • 300 V-Bucks\n  • 1 Music Track\n  • <Bold>70% Bonus</> Season Match XP\n  • <Bold>20% Bonus</> Season Friend Match XP\n  • <Bold>Extra Weekly Challenges</>\n  • and more!\n\nPlay to level up your Battle Pass, unlocking <Bold>over 75 rewards</> (typically takes 75 to 150 hours of play).\n  • <Bold>4 more Outfits</>\n  • <Bold>1,000 V-Bucks</>\n  • 6 Emotes\n  • 5 Wraps\n  • 3 Gliders\n  • 3 Back Blings\n  • 4 Harvesting Tools\n  • 4 Contrails\n  • 1 Pet\n  • 12 Sprays\n  • 2 Music Tracks\n  • and so much more!\nWant it all faster? You can use V-Bucks to buy tiers any time!",
					},
					{
						OfferID: "battlepass"+ fmt.Sprint(common.Season) +"",
						DevName: "BR.Season"+ fmt.Sprint(common.Season) +".BattlePass.01",
						OfferType: "StaticPrice",
						Prices: []models.Price{{
							CurrencyType: "MtxCurrency",
							CurrencySubType: "",
							RegularPrice: 950,
							FinalPrice: 950,
							SaleExpiration: "9999-12-31T23:59:59.999Z",
							BasePrice: 950,
						}},
						Categories: []string{},
						DailyLimit: -1,
						WeeklyLimit: -1,
						MonthlyLimit: -1,
						AppStoreID: []string{"", "", "", "", "", "", "", "", "", ""},
						Requirements: []models.Requirement{{
							RequirementType: "DenyOnFulfillment",
							RequiredID: "battlepass"+ fmt.Sprint(common.Season),
							MinQuantity: 1,
						}},
						MetaInfo: []gin.H{},
						CatalogGroupPriority: 0,
						SortPriority: 0,
						DisplayAssetPath: "/Game/Catalog/DisplayAssets/DA_BR_Season8_BattlePass.DA_BR_Season8_BattlePass",
						ItemGrants: []models.ItemGrant{},
						Refundable: false,
						Title: "Battle Pass",
						ShortDescription: "Season"+ fmt.Sprint(common.Season),
						Description: "Season 8 \n\nInstantly get these items <Bold>valued at over 3,500 V-Bucks</>.\n  • <ItemName>Blackheart</> Progressive Outfit\n  • <ItemName>Hybrid</> Progressive Outfit\n  • <Bold>50% Bonus</> Season Match XP\n  • <Bold>10% Bonus</> Season Friend Match XP\n  • <Bold>Extra Weekly Challenges</>\n\nPlay to level up your Battle Pass, unlocking <Bold>over 100 rewards</> (typically takes 75 to 150 hours of play).\n  • <ItemName>Sidewinder</> and <Bold>4 more Outfits</>\n  • <Bold>1,300 V-Bucks</>\n  • 7 Emotes\n  • 6 Wraps\n  • 2 Pets\n  • 5 Harvesting Tools\n  • 4 Gliders\n  • 4 Back Blings\n  • 5 Contrails\n  • 14 Sprays\n  • 3 Music Tracks\n  • 1 Toy\n  • 20 Loading Screens\n  • and so much more!\nWant it all faster? You can use V-Bucks to buy tiers any time!",
					},
				},
			},
		},
	}

	for i := 0; i < 6; i++ {
		entry := common.GenerateRandomCatalogEntry(-1, &dailyItems, "Small")
		ItemShop.Storefronts[0].CatalogEntries = append(ItemShop.Storefronts[0].CatalogEntries, entry)
	}

	for i := 0; i < 2; i++ {
		entry := common.GenerateRandomCatalogEntry(1, &legendaryItems, "Normal")
		ItemShop.Storefronts[1].CatalogEntries = append(ItemShop.Storefronts[1].CatalogEntries, entry)
	}
	
	all.PrintGreen([]any{"generated new random item shop"})
	if !common.LoadShopFromJson {
		SaveItemShop()
	}
}

func SaveItemShop() {
	file, err := os.OpenFile("data/shop/shop.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

func CheckGiftStatus(c *gin.Context) {
	offerId := c.Param("offerId")
	recipient := c.Param("recipientId")

	if offerId == "" || recipient == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
		return
	}

	itemShopEntry, err := common.GetCatalogEntry(offerId)
	if err != nil {
		common.ErrorBadRequest(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"price": itemShopEntry.Prices[0],
		"items": itemShopEntry.ItemGrants,
	})
}