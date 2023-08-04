package models

import (
	"gorm.io/gorm"
)

type UserProfile struct {
  gorm.Model
  AccountId string `gorm:"default:null" json:"accountId,omitempty"`
  ProfileId string `gorm:"default:null" json:"profileId,omitempty"`
  Profile string `gorm:"type:text" json:"profile,omitempty"`
}

type UserLoadout struct {
	gorm.Model
	AccountId		string	`gorm:"default:null" json:"accountId,omitempty"`
	LoadoutName string	`gorm:"default:null" json:"loadoutName,omitempty"`
	Loadout 		string	`gorm:"type:text" json:"loadout,omitempty"`
}

type Loadout struct {
	TemplateId string              `json:"templateId,omitempty"`
	Attributes LoadoutAttributes   `json:"attributes,omitempty"`
	Quantity   int                 `json:"quantity,omitempty"`
}

type LoadoutAttributes struct {
	LockerSlotsData LockerSlotsData `json:"locker_slots_data,omitempty"`
	UseCount        int             `json:"use_count,omitempty"`
	BannerIconTemplate string      `json:"banner_icon_template,omitempty"`
	LockerName string             `json:"locker_name,omitempty"`
	BannerColorTemplate string    `json:"banner_color_template,omitempty"`
	ItemSeen bool                `json:"item_seen,omitempty"`
	Favorite bool                `json:"favorite,omitempty"`
}

type LockerSlotsData struct {
	Slots map[string]LockerSlotItem `json:"slots,omitempty"`
}

type LockerSlotItem struct {
	Items         []string     `json:"items,omitempty"`
	ActiveVariants []Variant   `json:"activeVariants,omitempty"`
}

type Variant struct {
	Variants []string `json:"variants,omitempty"`
}

type Profile struct {
	Created         string             `json:"created,omitempty"`
	Updated         string             `json:"updated,omitempty"`
	Rvn             int                `json:"rvn,omitempty"`
	WipeNumber      int                `json:"wipeNumber,omitempty"`
	AccountId       string             `json:"accountId,omitempty"`
	ProfileId       string             `json:"profileId,omitempty"`
	Version         string             `json:"version,omitempty"`
	Items           map[string]any    `json:"items,omitempty"`
	Stats           Stats              `json:"stats,omitempty"`
	CommandRevision int                `json:"commandRevision,omitempty"`
}

type Item struct {
	Attributes ItemAttributes 	`json:"attributes,omitempty"`
	TemplateId string			 	`json:"templateId,omitempty"`
	Quantity   int				 	`json:"quantity,omitempty"`
}

type CommonCoreItem struct {
	Attributes map[string]any 	`json:"attributes,omitempty"`
	TemplateId string			 	`json:"templateId,omitempty"`
	Quantity   int				 	`json:"quantity,omitempty"`
}

type ItemAttributes struct {
	Favorite                    bool        `json:"favorite,omitempty"`
	ItemSeen                    bool        `json:"item_seen,omitempty"`
	Level                       int         `json:"level,omitempty"`
	MaxLevelBonus               int         `json:"max_level_bonus,omitempty"`
	RndSelCnt                   int         `json:"rnd_sel_cnt,omitempty"`
	Variants                    []any    `json:"variants,omitempty"`
	Xp                          int         `json:"xp,omitempty"`
	Platform 										*string      `json:"platform,omitempty"`
}

type Stats struct {
	Attributes StatsAttributes `json:"attributes,omitempty"`
}

type StatsAttributes struct {
	SeasonMatchBoost                 int               `json:"season_match_boost,omitempty"`
	Loadouts                         []string          `json:"loadouts,omitempty"`
	RestedXpOverflow                 int               `json:"rested_xp_overflow,omitempty"`
	MfaRewardClaimed                 bool              `json:"mfa_reward_claimed,omitempty"`
	QuestManager                     map[string]string `json:"quest_manager,omitempty"`
	BookLevel                        int               `json:"book_level,omitempty"`
	SeasonNum                        int               `json:"season_num,omitempty"`
	SeasonUpdate                     int               `json:"season_update,omitempty"`
	BookXp                           int               `json:"book_xp,omitempty"`
	Permissions                      []string          `json:"permissions,omitempty"`
	BookPurchased                    bool              `json:"book_purchased,omitempty"`
	LifetimeWins                     int               `json:"lifetime_wins,omitempty"`
	PartyAssistQuest                 string            `json:"party_assist_quest,omitempty"`
	PurchasedBattlePassTierOffers    []string          `json:"purchased_battle_pass_tier_offers,omitempty"`
	RestedXpExchange                 int               `json:"rested_xp_exchange,omitempty"`
	Level                            int               `json:"level,omitempty"`
	XpOverflow                       int               `json:"xp_overflow,omitempty"`
	RestedXp                         int               `json:"rested_xp,omitempty"`
	RestedXpMult                     int               `json:"rested_xp_mult,omitempty"`
	AccountLevel                     int               `json:"accountLevel,omitempty"`
	CompetitiveIdentity              map[string]string `json:"competitive_identity,omitempty"`
	InventoryLimitBonus              int               `json:"inventory_limit_bonus,omitempty"`
	LastAppliedLoadout               string            `json:"last_applied_loadout,omitempty"`
	DailyRewards                     map[string]string `json:"daily_rewards,omitempty"`
	Xp                               int               `json:"xp,omitempty"`
	SeasonFriendMatchBoost           int               `json:"season_friend_match_boost,omitempty"`
	ActiveLoadoutIndex               int               `json:"active_loadout_index,omitempty"`
	FavoriteMusicPack                string            `json:"favorite_musicpack,omitempty"`
	FavoriteGlider                   string            `json:"favorite_glider,omitempty"`
	FavoritePickaxe                  string            `json:"favorite_pickaxe,omitempty"`
	FavoriteSkyDiveContrail          string            `json:"favorite_skydivecontrail,omitempty"`
	FavoriteBackpack                 string            `json:"favorite_backpack,omitempty"`
	FavoriteDance                    []string          `json:"favorite_dance,omitempty"`
	FavoriteItemWraps                []string          `json:"favorite_itemwraps,omitempty"`
	FavoriteCharacter                string            `json:"favorite_character,omitempty"`
	FavoriteLoadingScreen            string            `json:"favorite_loadingscreen,omitempty"`
}

type ProfileResponse struct {
	ProfileRevision            int           `json:"profileRevision,omitempty"`
	ProfileID                  string        `json:"profileId,omitempty"`
	ProfileChangesBaseRevision int           `json:"profileChangesBaseRevision,omitempty"`
	ProfileCommandRevision     int           `json:"profileCommandRevision,omitempty"`
	ServerTime                 string     `json:"serverTime,omitempty"`
	ResponseVersion            int           `json:"responseVersion,omitempty"`
	ProfileChanges             []ProfileChange `json:"profileChanges,omitempty"`
	MultiUpdate                []MultiUpdate `json:"multiUpdate,omitempty"`
	Notifications              []Notification `json:"notifications,omitempty"`
}

type ProfileChange struct {
	ChangeType  string 	`json:"changeType,omitempty"`
	ItemID      string 	`json:"itemId,omitempty"`
	Quantity    int    	`json:"quantity,omitempty"`
	Item   Item   	`json:"item,omitempty"`
	Profile     Profile `json:"profile,omitempty"`
	Name        string 	`json:"name,omitempty"`
	Value       any    	`json:"value,omitempty"`
}

type Notification struct {
	Type      string `json:"type,omitempty"`
	Primary   bool   `json:"primary,omitempty"`
	LootResult LootResult `json:"lootResult,omitempty"`
}

type LootResult struct {
	Items []LootResultItem `json:"items,omitempty"`
}

type LootResultItem struct {
	ItemType string `json:"itemType,omitempty"`
	ItemGuid string `json:"itemGuid,omitempty"`
	ItemProfile string `json:"itemProfile,omitempty"`
	Quantity int `json:"quantity,omitempty"`
}

type MultiUpdate struct {
	ProfileRevision            int           `json:"profileRevision,omitempty"`
	ProfileID                  string        `json:"profileId,omitempty"`
	ProfileChangesBaseRevision int           `json:"profileChangesBaseRevision,omitempty"`
	ProfileChanges             []ProfileChange `json:"profileChanges,omitempty"`
	ProfileCommandRevision     int           `json:"profileCommandRevision,omitempty"`
}