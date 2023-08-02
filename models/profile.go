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
  ProfileId string `json:"profileId"`
  Version string `json:"version"`
  Items interface{} `json:"items"`
  Stats struct {
    Attributes StatsAttributes `json:"attributes"`
  } `json:"stats"`
  CommandRevision int `json:"commandRevision"`
  Created string `json:"created"`
  Updated string `json:"updated"`
}

type Item struct {
  ItemName string `json:"itemName"`
  TemplateId string `json:"templateId"`
  Attributes ItemAttributes `json:"attributes"`
  Quantity int `json:"quantity"`
}

type ItemAttributes struct {
  ItemSeen bool `json:"item_seen"`
  InventoryLimitBonus int `json:"inventory_limit_bonus"`
  MaxLevelBonus int `json:"max_level_bonus"`
  Level int `json:"level"`
  Xp int `json:"xp"`
  Variants []int `json:"variants"`
  Favorite bool `json:"favorite"`
}


type StatsAttributes struct {
  SeasonMatchBoost          int                 `json:"season_match_boost"`
  Loadouts                  []string            `json:"loadouts"`
  RestedXpOverflow          int                 `json:"rested_xp_overflow"`
  MfaRewardClaimed          bool                `json:"mfa_reward_claimed"`
  QuestManager              interface{}         `json:"quest_manager"`
  BookLevel                 int                 `json:"book_level"`
  SeasonNum                 int                 `json:"season_num"`
  SeasonUpdate              int                 `json:"season_update"`
  BookXp                    int                 `json:"book_xp"`
  Permissions               []string            `json:"permissions"`
  BookPurchased             bool                `json:"book_purchased"`
  LifetimeWins              int                 `json:"lifetime_wins"`
  PartyAssistQuest          string              `json:"party_assist_quest"`
  PurchasedBattlePassTierOffers []string        `json:"purchased_battle_pass_tier_offers"`
  RestedXpExchange          float64             `json:"rested_xp_exchange" gorm:"default:0.0"`
  Level                     int                 `json:"level"`
  XpOverflow                int                 `json:"xp_overflow"`
  RestedXp                  int                 `json:"rested_xp"`
  RestedXpMult              int                 `json:"rested_xp_mult"`
  AccountLevel              int                 `json:"accountLevel"`
  CompetitiveIdentity       interface{}         `json:"competitive_identity"`
  LastAppliedLoadout        string              `json:"last_applied_loadout"`
  DailyRewards              interface{}         `json:"daily_rewards"`
  Xp                        int                 `json:"xp"`
  SeasonFriendMatchBoost    int                 `json:"season_friend_match_boost"`
  ActiveLoadoutIndex        int                 `json:"active_loadout_index"`
  FavoriteMusicpack         string              `json:"favorite_musicpack"`
  FavoriteGlider            string              `json:"favorite_glider"`
  FavoritePickaxe           string              `json:"favorite_pickaxe"`
  FavoriteSkydivecontrail   string              `json:"favorite_skydivecontrail"`
  FavoriteBackpack          string              `json:"favorite_backpack"`
  FavoriteDance             []string            `json:"favorite_dance"`
  FavoriteItemwraps         []string            `json:"favorite_itemwraps"`
  FavoriteCharacter         string              `json:"favorite_character"`
  FavoriteLoadingscreen     string              `json:"favorite_loadingscreen"`
  FavoriteVictoryPose       string              `json:"favorite_victorypose"`
  FavoriteConsumableEmote   string              `json:"favorite_consumableemote"`
  BannerColor               string              `json:"banner_color"`
  FavoriteCallingCard       string              `json:"favorite_callingcard"`
  FavoriteSpray             []string            `json:"favorite_spray"`
  FavoriteHat               string              `json:"favorite_hat"`
  FavoriteBattleBus         string              `json:"favorite_battlebus"`
  FavoriteMapMarker         string              `json:"favorite_mapmarker"`
  FavoriteVehicleDeco       string              `json:"favorite_vehicledeco"`
  BannerIcon                string              `json:"banner_icon"`
  SurveyData                interface{}         `json:"survey_data"`
  PersonalOffers            interface{}         `json:"personal_offers"`
  IntroGamePlayed           bool                `json:"intro_game_played"`
  ImportFriendsClaimed      interface{}         `json:"import_friends_claimed"`
  MTXPurchaseHistory        struct {
    RefundsUsed   int      `json:"refundsUsed"`
    RefundCredits int      `json:"refundCredits"`
    Purchases     []string `json:"purchases"`
  }                                             `json:"mtx_purchase_history"`
  UndoCooldowns             []string            `json:"undo_cooldowns"`
  MTXAffiliateSetTime       string              `json:"mtx_affiliate_set_time"`
  CurrentMTXPlatform        string              `json:"current_mtx_platform"`
  MTXAffiliate              string              `json:"mtx_affiliate"`
  ForcedIntroPlayed         string              `json:"forced_intro_played"`
  WeeklyPurchases           interface{}         `json:"weekly_purchases"`
  DailyPurchases            interface{}         `json:"daily_purchases"`
  BanHistory                interface{}         `json:"ban_history"`
  InAppPurchases            interface{}         `json:"in_app_purchases"`
  UndoTimeout               string              `json:"undo_timeout"`
  MonthlyPurchases          interface{}         `json:"monthly_purchases"`
  AllowedToSendGifts        bool                `json:"allowed_to_send_gifts"`
  MfaEnabled                bool                `json:"mfa_enabled"`
  AllowedToReceiveGifts     bool                `json:"allowed_to_receive_gifts"`
  GiftHistory               interface{}         `json:"gift_history"`
}

type DefaultLoadout struct {
	TemplateID string `json:"templateId"`
	Attributes struct {
		LockerSlotsData struct {
			Slots struct {
				Pickaxe struct {
					Items        []string     `json:"items"`
					ActiveVariants interface{} `json:"activeVariants"`
				} `json:"Pickaxe"`
				Dance struct {
					Items []string `json:"items"`
				} `json:"Dance"`
				Glider struct {
					Items []string `json:"items"`
				} `json:"Glider"`
				Character struct {
					Items         []string `json:"items"`
					ActiveVariants []struct {
						Variants interface{} `json:"variants"`
					} `json:"activeVariants"`
				} `json:"Character"`
				Backpack struct {
					Items         []string `json:"items"`
					ActiveVariants []struct {
						Variants interface{} `json:"variants"`
					} `json:"activeVariants"`
				} `json:"Backpack"`
				ItemWrap struct {
					Items         []string `json:"items"`
					ActiveVariants []interface{} `json:"activeVariants"`
				} `json:"ItemWrap"`
				LoadingScreen struct {
					Items         []string `json:"items"`
					ActiveVariants []interface{} `json:"activeVariants"`
				} `json:"LoadingScreen"`
				MusicPack struct {
					Items         []string `json:"items"`
					ActiveVariants []interface{} `json:"activeVariants"`
				} `json:"MusicPack"`
				SkyDiveContrail struct {
					Items         []string `json:"items"`
					ActiveVariants []interface{} `json:"activeVariants"`
				} `json:"SkyDiveContrail"`
			} `json:"slots"`
		} `json:"locker_slots_data"`
		UseCount            int  `json:"use_count"`
		BannerIconTemplate  string `json:"banner_icon_template"`
		LockerName          string `json:"locker_name"`
		BannerColorTemplate string `json:"banner_color_template"`
		ItemSeen            bool   `json:"item_seen"`
		Favorite            bool   `json:"favorite"`
	} `json:"attributes"`
	Quantity int `json:"quantity"`
}