package common

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/zombman/server/models"
)

func GetAllFortniteItems() ([]models.BeforeStoreItem, error) {
	pathToProfile := "default/items.json"

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

	var itemsData []models.BeforeStoreItem
	err = json.Unmarshal([]byte(str), &itemsData)
	if err != nil {
		return []models.BeforeStoreItem{}, err
	}

	ExcludeType(&itemsData, "AthenaSpray")
	ExcludeType(&itemsData, "AthenaEmoji")
	ExcludeType(&itemsData, "AthenaToy")
	ExcludeType(&itemsData, "AthenaBackpack")
	ExcludeType(&itemsData, "AthenaPetCarrier")

	return itemsData, nil
}

func GetItemsFromSeason(season int) ([]models.BeforeStoreItem, error) {
	items, err := GetAllFortniteItems()
	if err != nil {
		return []models.BeforeStoreItem{}, err
	}

	var itemsFromSeason []models.BeforeStoreItem
	for _, item := range items {
		if item.IntroductionSeason > season {
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