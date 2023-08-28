package models

import (
	"github.com/gin-gonic/gin"
)

type StorePage struct {
	RefreshIntervalHrs int          `json:"refreshIntervalHrs"`
	DailyPurchaseHrs   int          `json:"dailyPurchaseHrs"`
	Expiration         string		    `json:"expiration"`
	Storefronts        []Storefront `json:"storefronts"`
}

type Storefront struct {
	Name           string         `json:"name"`
	CatalogEntries []CatalogEntry `json:"catalogEntries"`
}

type CatalogEntry struct {
	DevName          string            `json:"devName"`
	OfferID          string            `json:"offerId"`
	FulfillmentIDs   []string          `json:"fulfillmentIds"`
	DailyLimit       int               `json:"dailyLimit"`
	WeeklyLimit      int               `json:"weeklyLimit"`
	MonthlyLimit     int               `json:"monthlyLimit"`
	Categories       []string          `json:"categories"`
	Prices           []Price           `json:"prices"`
	Meta             map[string]any `json:"meta"`
	MatchFilter      string            `json:"matchFilter"`
	FilterWeight     int               `json:"filterWeight"`
	AppStoreID       []string          `json:"appStoreId"`
	Requirements     []Requirement     `json:"requirements"`
	OfferType        string            `json:"offerType"`
	GiftInfo         map[string]any `json:"giftInfo"`
	Refundable       bool              `json:"refundable"`
	MetaInfo         []gin.H          `json:"metaInfo"`
	DisplayAssetPath string            `json:"displayAssetPath"`
	ItemGrants       []ItemGrant       `json:"itemGrants"`
	SortPriority     int               `json:"sortPriority"`
	CatalogGroupPriority int               `json:"catalogGroupPriority"`
	Title						string            `json:"title"`
	ShortDescription	string            `json:"shortDescription"`
	Description				string            `json:"description"`
}

type Price struct {
	CurrencyType    string    `json:"currencyType"`
	CurrencySubType string    `json:"currencySubType"`
	RegularPrice    int       `json:"regularPrice"`
	FinalPrice      int       `json:"finalPrice"`
	SaleExpiration  string		`json:"saleExpiration"`
	BasePrice       int       `json:"basePrice"`
}

type Requirement struct {
	RequirementType string `json:"requirementType"`
	RequiredID      string `json:"requiredId"`
	MinQuantity     int    `json:"minQuantity"`
}

type ItemGrant struct {
	TemplateID string `json:"templateId"`
	Quantity   int    `json:"quantity"`
}

type BeforeStoreItem struct {
	ID                	string `json:"id"`
	BackendType       	string `json:"backendType"`
	IntroductionSeason 	int    `json:"introductionSeason"`
	Rarity            	string `json:"rarity"`
	Gender							string `json:"gender"`
	Set									string `json:"set"`
}

type BeforeStoreSetMap map[string][]string

type Prices map[string]map[string]int