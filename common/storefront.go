package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/all"
	"github.com/zombman/server/models"
)
var (
	AllItemsKeys map[string]models.BeforeStoreItem
	AllItems []models.BeforeStoreItem
	AllItemsItemShop []models.BeforeStoreItem
	AllSets models.BeforeStoreSetMap
	Prices = models.Prices{
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
)

func GetAllFortniteItems() ([]models.BeforeStoreItem, error) {
	if len(AllItems) == 0 {
		pathToProfile := "data/shop/all.json"
	
		file, err := os.Open(pathToProfile)
		if err != nil {
			return []models.BeforeStoreItem{}, err
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			return []models.BeforeStoreItem{}, err
		}
		str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

		err = json.Unmarshal([]byte(str), &AllItems)
		if err != nil {
			return []models.BeforeStoreItem{}, err
		}

		AllItemsItemShop = AllItems

		tempAllItemsKey := make(map[string]models.BeforeStoreItem)
		for _, item := range AllItems {
			tempAllItemsKey[item.BackendType + ":" + item.ID] = item
		}
		AllItemsKeys = tempAllItemsKey
	}

	ExcludeType(&AllItemsItemShop, "AthenaSpray")
	ExcludeType(&AllItemsItemShop, "AthenaEmoji")
	ExcludeType(&AllItemsItemShop, "AthenaToy")
	ExcludeType(&AllItemsItemShop, "AthenaBackpack")
	ExcludeType(&AllItemsItemShop, "AthenaPetCarrier")
	ExcludeType(&AllItemsItemShop, "AthenaPet")
	ExcludeType(&AllItemsItemShop, "AthenaBackpack")

	return AllItemsItemShop, nil
}

func GetAllSets() (models.BeforeStoreSetMap, error) {
	if len(AllSets) == 0 {
		pathToSets := "data/shop/sets.json"

		file, err := os.Open(pathToSets)
		if err != nil {
			return models.BeforeStoreSetMap{}, err
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			return models.BeforeStoreSetMap{}, err
		}
		str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

		err = json.Unmarshal([]byte(str), &AllSets)
		if err != nil {
			return models.BeforeStoreSetMap{}, err
		}
	}

	return AllSets, nil
}

func GetItemsFromSeason(season int) ([]models.BeforeStoreItem, error) {
	items, err := GetAllFortniteItems()
	if err != nil {
		return []models.BeforeStoreItem{}, err
	}

	var itemsFromSeason []models.BeforeStoreItem
	for _, item := range items {
		if item.IntroductionSeason > season + 10 {
			continue
		}
		
		itemsFromSeason = append(itemsFromSeason, item)
	}

	return itemsFromSeason, nil
}

func FilterRarity(items []models.BeforeStoreItem, rarity string) []models.BeforeStoreItem {
	var filteredItems []models.BeforeStoreItem
	for _, item := range items {
		if item.Rarity == rarity {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func ExcludeRarity(items []models.BeforeStoreItem, rarity string) []models.BeforeStoreItem{
	var filteredItems []models.BeforeStoreItem
	for _, item := range items {
		if item.Rarity != rarity {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func ExcludeType(items *[]models.BeforeStoreItem, backendType string) {
	var filteredItems []models.BeforeStoreItem
	for _, item := range *items {
		if item.BackendType != backendType {
			filteredItems = append(filteredItems, item)
		}
	}

	*items = filteredItems
}

func ExcludeId(items *[]models.BeforeStoreItem, id string) {
	var filteredItems []models.BeforeStoreItem
	for _, item := range *items {
		if item.ID != id {
			filteredItems = append(filteredItems, item)
		}
	}

	*items = filteredItems
}

func GetCatalogEntry(offerId string) (models.CatalogEntry, error) {
	itemshop := GetItemShop()

	var entry models.CatalogEntry
	for _, storefront := range itemshop.Storefronts {
		for _, catalogEntry := range storefront.CatalogEntries {
			if catalogEntry.OfferID == offerId {
				entry = catalogEntry
			}
		}
	}

	if entry.OfferID == "" {
		return models.CatalogEntry{}, fmt.Errorf("could not find catalog entry with offerId %s", offerId)
	}

	return entry, nil
}

func GetItemShop() models.StorePage {
	pathToProfile := "data/shop/shop.json"

	file, err := os.Open(pathToProfile)
	if err != nil {
		return models.StorePage{}
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return models.StorePage{}
	}
	str := string(bytes.ReplaceAll(bytes.ReplaceAll(fileData, []byte("\n"), []byte("")), []byte("\t"), []byte("")))

	var itemshop models.StorePage
	err = json.Unmarshal([]byte(str), &itemshop)
	if err != nil {
		return models.StorePage{}
	}

	return itemshop
}

func GetBackpackItemGrant(btype string, mainItemGrant string) (models.ItemGrant, error) {
	allSets, err := GetAllSets()
	if err != nil {
		all.PrintRed([]any{"error getting all sets", err})
		return models.ItemGrant{}, err
	}
	
	item, ok := AllItemsKeys[btype + ":" + mainItemGrant]
	if !ok {
		return models.ItemGrant{}, fmt.Errorf("could not find item with id %s", mainItemGrant)
	}
	
	setItems, ok := allSets[item.Set]
	if !ok {
		return models.ItemGrant{}, fmt.Errorf("could not find set with id %s", item.Set)
	}

	for _, setItem := range setItems {
		if !strings.Contains(setItem, "AthenaBackpack") {
			continue
		}

		if strings.Contains(item.ID, "Commando_M") {
			if strings.Contains(setItem, "Male") {
				return models.ItemGrant{
					TemplateID: setItem,
					Quantity: 1,
				}, nil
			}
		}

		if strings.Contains(item.ID, "Commando_F") {
			if strings.Contains(setItem, "Female") {
				return models.ItemGrant{
					TemplateID: setItem,
					Quantity: 1,
				}, nil
			}
		}
		
		return models.ItemGrant{
			TemplateID: setItem,
			Quantity: 1,
		}, nil
	}

	return models.ItemGrant{}, fmt.Errorf("could not find item with id %s", mainItemGrant)
}

func GenerateRandomCatalogEntry(f int, items *[]models.BeforeStoreItem, size string) models.CatalogEntry {
	endOfDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Now().Location()).Format("2006-01-02T15:04:05.999Z")
	randomItem := (*items)[rand.Intn(len(*items) - 1)]
	price := Prices[randomItem.BackendType][randomItem.Rarity]
	id := all.HashString(randomItem.ID)

	itemGrants := []models.ItemGrant{
		{
			TemplateID: randomItem.BackendType + ":" + randomItem.ID,
			Quantity: 1,
		},
	}

	if randomItem.BackendType == "AthenaCharacter" {
		item, ok := GetBackpackItemGrant(randomItem.BackendType, randomItem.ID)

		all.PrintRed([]any{ok})

		if ok == nil {
			itemGrants = append(itemGrants, item)
		}
	}

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
			"bIsEnabled": true,
			"forcedGiftBoxTemplateId": "",
			"purchaseRequirements": []any{},
			"giftRecordIds": []any{},
		},
		Meta: map[string]any {
			"SectionId": "Featured",
			"TileSize": size,
		},
		MetaInfo: []gin.H{
			{
				"key": "SectionId",
				"value": "Featured",
			}, {
				"key": "TileSize",
				"value": size,
			},
		},
		Requirements: []models.Requirement{{
			RequirementType: "DenyOnItemOwnership",
			RequiredID: randomItem.BackendType + ":" + randomItem.ID,
			MinQuantity: 1,
		}},
		ItemGrants: itemGrants,
	}
}