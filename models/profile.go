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

type Profile struct {
	Created         string             `json:"created"`
	Updated         string             `json:"updated"`
	Rvn             int                `json:"rvn"`
	WipeNumber      int                `json:"wipeNumber"`
	AccountId       string             `json:"accountId"`
	ProfileId       string             `json:"profileId"`
	Version         string             `json:"version"`
	Items           map[string]Item    `json:"items"`
	Stats           Stats              `json:"stats"`
	CommandRevision int                `json:"commandRevision"`
}

type Item struct {
	Attributes ItemAttributes 	`json:"attributes"`
	TemplateId string			 	`json:"templateId"`
	Quantity   int					`json:"quantity"`
}

type Loadout struct {
	Attributes interface{} 	`json:"attributes"`
	TemplateId string			 	`json:"templateId"`
	Quantity   int					`json:"quantity"`
}

type ItemAttributes struct {
	Favorite                    bool        `json:"favorite"`
	ItemSeen                    bool        `json:"item_seen"`
	Level                       int         `json:"level"`
	MaxLevelBonus               int         `json:"max_level_bonus"`
	RndSelCnt                   int         `json:"rnd_sel_cnt"`
	Variants                    []string    `json:"variants"`
	Xp                          int         `json:"xp"`
}

type LoadoutAttributes struct {
	LockerSlotsData 		LockerSlots `json:"locker_slots_data"`
	UseCount        		int         `json:"use_count"`
	BannerIconTemplate 	string 			`json:"banner_icon_template"`
	LockerName 					string 			`json:"locker_name"`
	BannerColorTemplate string 			`json:"banner_color_template"`
	ItemSeen					 	bool 				`json:"item_seen"`
	Favorite 						bool 				`json:"favorite"`
}

type LockerSlots struct {
	Slots map[string]LockerSlotItem `json:"slots"`
}

type LockerSlotItem struct {
	Items         	[]string 			`json:"items"`
	ActiveVariants 	[]interface{} `json:"activeVariants"`
}

type Stats struct {
	Attributes StatsAttributes `json:"attributes"`
}

type StatsAttributes struct {
	SeasonMatchBoost                 int               `json:"season_match_boost"`
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
	RestedXpMult                     int               `json:"rested_xp_mult"`
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
	FavoriteLoadingscreen            string            `json:"favorite_loadingscreen"`
}