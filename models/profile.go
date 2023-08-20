package models

import (
	"gorm.io/gorm"
)

type UserProfile struct {
  gorm.Model
  AccountId string `gorm:"default:null" json:"accountId"`
  ProfileId string `gorm:"default:null" json:"profileId"`
  Profile string `gorm:"type:text" json:"profile"`
}

type UserLoadout struct {
	gorm.Model
	AccountId		string	`gorm:"default:null" json:"accountId"`
	LoadoutName string	`gorm:"default:null" json:"loadoutName"`
	Loadout 		string	`gorm:"type:text" json:"loadout"`
}

type Loadout struct {
	TemplateId string              `json:"templateId"`
	Attributes LoadoutAttributes   `json:"attributes"`
	Quantity   int                 `json:"quantity"`
}

type LoadoutAttributes struct {
	LockerSlotsData LockerSlotsData `json:"locker_slots_data"`
	UseCount        int             `json:"use_count"`
	BannerIconTemplate string      `json:"banner_icon_template"`
	LockerName string             `json:"locker_name"`
	BannerColorTemplate string    `json:"banner_color_template"`
	ItemSeen bool                `json:"item_seen"`
	Favorite bool                `json:"favorite"`
}

type LockerSlotsData struct {
	Slots map[string]LockerSlotItem `json:"slots"`
}

type LockerSlotItem struct {
	Items         []string     `json:"items"`
	ActiveVariants []Variant   `json:"activeVariants"`
}

type Variant struct {
	Variants []string `json:"variants"`
}

type Profile struct {
	Created         string             `json:"created"`
	Updated         string             `json:"updated"`
	Rvn             int                `json:"rvn"`
	WipeNumber      int                `json:"wipeNumber"`
	AccountId       string             `json:"accountId"`
	ProfileId       string             `json:"profileId"`
	Version         string             `json:"version"`
	Items           map[string]any    `json:"items"`
	Stats           struct {
		Attributes map[string]any `json:"attributes"`
	}              `json:"stats"`
	CommandRevision int                `json:"commandRevision"`
}

type Item struct {
	Attributes ItemAttributes 	`json:"attributes,"`
	TemplateId string			 	`json:"templateId"`
	Quantity   int				 	`json:"quantity"`
}

type CommonCoreItem struct {
	Attributes map[string]any 	`json:"attributes"`
	TemplateId string			 	`json:"templateId"`
	Quantity   int				 	`json:"quantity"`
}

type ItemVariant struct {
	Channel string `json:"channel"`
	Active 	string `json:"active"`
	Owned []string `json:"owned"`
}

type ItemAttributes struct {
	Favorite                    bool        `json:"favorite"`
	ItemSeen                    bool        `json:"item_seen"`
	Level                       int         `json:"level"`
	MaxLevelBonus               int         `json:"max_level_bonus"`
	RndSelCnt                   int         `json:"rnd_sel_cnt"`
	Variants                    []ItemVariant    `json:"variants"`
	Xp                          int         `json:"xp"`
	Platform 										*string      `json:"platform"`
}

type AthenaProfile struct {
	Created string `json:"created"`
	Updated string `json:"updated"`
	Rvn int `json:"rvn"`
	WipeNumber int `json:"wipeNumber"`
	AccountId string `json:"accountId"`
	ProfileId string `json:"profileId"`
	Version string `json:"version"`
	Items map[string]any    `json:"items"`
	Stats Stats `json:"stats"`
	CommandRevision int `json:"commandRevision"`
}

type CommonCoreProfile struct {
	Created string `json:"created"`
	Updated string `json:"updated"`
	Rvn int `json:"rvn"`
	WipeNumber int `json:"wipeNumber"`
	AccountId string `json:"accountId"`
	ProfileId string `json:"profileId"`
	Version string `json:"version"`
	Items map[string]CommonCoreItem `json:"items"`
	Stats CommonCoreStats `json:"stats"`
	CommandRevision int `json:"commandRevision"`
}


type CommonCoreStats struct {
	Attributes CommonCoreStatsAttributes `json:"attributes"`
}

type CommonCoreStatsAttributes struct {
	MtxPurchaseHistory struct {
		RefundsUsed int `json:"refundsUsed"`
		RefundCredits int `json:"refundCredits"`
		Purchases []any `json:"purchases"`
	} `json:"mtx_purchase_history"`
	CurrentMtxPlatform string `json:"current_mtx_platform"`
	MtxAffiliate string `json:"mtx_affiliate"`
}

type Stats struct {
	Attributes StatsAttributes `json:"attributes"`
}

type StatsAttributes struct {
	SeasonMatchBoost                 float32               `json:"season_match_boost"`
	Loadouts                         []string          `json:"loadouts"`
	RestedXpOverflow                 int               `json:"rested_xp_overflow"`
	MfaRewardClaimed                 bool              `json:"mfa_reward_claimed"`
	QuestManager                     map[string]string `json:"quest_manager"`
	BookLevel                        int               `json:"book_level"`
	SeasonNum                        int               `json:"season_num"`
	SeasonUpdate                     int               `json:"season_update"`
	BookXp                           int               `json:"book_xp"`
	Permissions                      []string          `json:"permissions"`
	BookPurchased                    bool              `json:"book_purchased"`
	LifetimeWins                     int               `json:"lifetime_wins"`
	PartyAssistQuest                 string            `json:"party_assist_quest"`
	PurchasedBattlePassTierOffers    []string          `json:"purchased_battle_pass_tier_offers"`
	RestedXpExchange                 int               `json:"rested_xp_exchange"`
	Level                            int               `json:"level"`
	XpOverflow                       int               `json:"xp_overflow"`
	RestedXp                         int               `json:"rested_xp"`
	RestedXpMult                     float32               `json:"rested_xp_mult"`
	AccountLevel                     int               `json:"accountLevel"`
	CompetitiveIdentity              map[string]string `json:"competitive_identity"`
	InventoryLimitBonus              int               `json:"inventory_limit_bonus"`
	LastAppliedLoadout               string            `json:"last_applied_loadout"`
	DailyRewards                     map[string]string `json:"daily_rewards"`
	Xp                               int               `json:"xp"`
	SeasonFriendMatchBoost           int               `json:"season_friend_match_boost"`
	ActiveLoadoutIndex               int               `json:"active_loadout_index"`
	FavoriteMusicPack                string            `json:"favorite_musicpack"`
	FavoriteGlider                   string            `json:"favorite_glider"`
	FavoritePickaxe                  string            `json:"favorite_pickaxe"`
	FavoriteSkyDiveContrail          string            `json:"favorite_skydivecontrail"`
	FavoriteBackpack                 string            `json:"favorite_backpack"`
	FavoriteDance                    []string          `json:"favorite_dance"`
	FavoriteItemWraps                []string          `json:"favorite_itemwraps"`
	FavoriteCharacter                string            `json:"favorite_character"`
	FavoriteLoadingScreen            string            `json:"favorite_loadingscreen"`
	BannerIcon                       string            `json:"banner_icon"`
	BannerColor                      string            `json:"banner_color"`
}

type ProfileResponse struct {
	ProfileRevision            int           `json:"profileRevision"`
	ProfileID                  string        `json:"profileId"`
	ProfileChangesBaseRevision int           `json:"profileChangesBaseRevision"`
	ProfileCommandRevision     int           `json:"profileCommandRevision"`
	ServerTime                 string     `json:"serverTime"`
	ResponseVersion            int           `json:"responseVersion"`
	ProfileChanges             []ProfileChange `json:"profileChanges"`
	MultiUpdate                []MultiUpdate `json:"multiUpdate"`
	Notifications              []Notification `json:"notifications"`
}

type ProfileChange struct {
	ChangeType  string 	`json:"changeType"`
	ItemID      string 	`json:"itemId"`
	Quantity    int    	`json:"quantity"`
	Item   Item   	`json:"item"`
	Profile     Profile `json:"profile"`
	Name        string 	`json:"name"`
	Value       any    	`json:"value"`
	AttributeName string `json:"attributeName"`
	AttributeValue any `json:"attributeValue"`
	ItemId 		string 	`json:"itemId"`
}

type Notification struct {
	Type      string `json:"type"`
	Primary   bool   `json:"primary"`
	LootResult LootResult `json:"lootResult"`
}

type LootResult struct {
	Items []LootResultItem `json:"items"`
}

type LootResultItem struct {
	ItemType string `json:"itemType"`
	ItemGuid string `json:"itemGuid"`
	ItemProfile string `json:"itemProfile"`
	Quantity int `json:"quantity"`
}

type MultiUpdate struct {
	ProfileRevision            int           `json:"profileRevision"`
	ProfileID                  string        `json:"profileId"`
	ProfileChangesBaseRevision int           `json:"profileChangesBaseRevision"`
	ProfileChanges             []ProfileChange `json:"profileChanges"`
	ProfileCommandRevision     int           `json:"profileCommandRevision"`
}