package models

type Message struct {
	Title map[string]string `json:"title"`
	Body  map[string]string `json:"body"`
}

type MessageS struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type LoginMessage struct {
	Title        string   `json:"_title"`
	LoginMessage MessageS `json:"loginmessage"`
	ActiveDate   string   `json:"_activeDate"`
	LastModified string   `json:"lastModified"`
	Locale       string   `json:"_locale"`
}

type SurvivalMessage struct {
	Title        string  `json:"_title"`
	Overrideable Message `json:"overrideablemessage"`
	ActiveDate   string  `json:"_activeDate"`
	LastModified string  `json:"lastModified"`
	Locale       string  `json:"_locale"`
}

type AthenaMessage struct {
	Title        string  `json:"_title"`
	Overrideable Message `json:"overrideablemessage"`
	ActiveDate   string  `json:"_activeDate"`
	LastModified string  `json:"lastModified"`
	Locale       string  `json:"_locale"`
}

type SubgameSelectData struct {
	SaveTheWorldUnowned struct {
		Type         string  `json:"_type"`
		Message      Message `json:"message"`
		Title        string  `json:"title"`
		ActiveDate   string  `json:"_activeDate"`
		LastModified string  `json:"lastModified"`
		Locale       string  `json:"_locale"`
	} `json:"saveTheWorldUnowned"`
	Title        string `json:"_title"`
	BattleRoyale struct {
		Type         string  `json:"_type"`
		Message      Message `json:"message"`
		Title        string  `json:"title"`
		ActiveDate   string  `json:"_activeDate"`
		LastModified string  `json:"lastModified"`
		Locale       string  `json:"_locale"`
	} `json:"battleRoyale"`
	Creative struct {
		Type         string  `json:"_type"`
		Message      Message `json:"message"`
		Title        string  `json:"title"`
		ActiveDate   string  `json:"_activeDate"`
		LastModified string  `json:"lastModified"`
		Locale       string  `json:"_locale"`
	} `json:"creative"`
	SaveTheWorld struct {
		Type         string  `json:"_type"`
		Message      Message `json:"message"`
		Title        string  `json:"title"`
		ActiveDate   string  `json:"_activeDate"`
		LastModified string  `json:"lastModified"`
		Locale       string  `json:"_locale"`
	} `json:"saveTheWorld"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type SaveTheWorldNews struct {
	News struct {
		Type     string    `json:"_type"`
		Messages []Message `json:"messages"`
	} `json:"news"`
	Title        string `json:"_title"`
	NoIndex      bool   `json:"_noIndex"`
	AlwaysShow   bool   `json:"alwaysShow"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type BattlePassAboutMessages struct {
	News struct {
		Type     string    `json:"_type"`
		Messages []Message `json:"messages"`
	} `json:"news"`
	Title        string `json:"_title"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type PlaylistInformation struct {
	FrontendMatchmakingHeaderStyle string `json:"frontend_matchmaking_header_style"`
	Title                          string `json:"_title"`
	FrontendMatchmakingHeaderText  string `json:"frontend_matchmaking_header_text"`
	PlaylistInfo                   struct {
		Type      string `json:"_type"`
		Playlists []struct {
			Image        string `json:"image"`
			PlaylistName string `json:"playlist_name"`
			Violator     string `json:"violator"`
			Type         string `json:"_type"`
			Description  string `json:"description"`
		} `json:"playlists"`
	} `json:"playlist_info"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type TournamentInformation struct {
	TournamentInfo struct {
		Type        string `json:"_type"`
		Tournaments []struct {
			TitleColor         string `json:"title_color"`
			LoadingScreenImage string `json:"loading_screen_image"`
			Background         string `json:"background"`
			HasCountdown       bool   `json:"has_countdown"`
			Title              string `json:"title"`
			Type               string `json:"_type"`
			Description        string `json:"description"`
		} `json:"tournaments"`
	} `json:"tournament_info"`
	Title        string `json:"_title"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type EmergencyNotice struct {
	EmergencyNotices []struct {
		Type     string `json:"_type"`
		Messages []struct {
			Type   string `json:"_type"`
			Region string `json:"region"`
			Title  string `json:"title"`
			Body   string `json:"body"`
		} `json:"messages"`
		Region   string `json:"region"`
		Platform string `json:"platform"`
	} `json:"emergencyNotices"`
	Title        string `json:"_title"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type EmergencyNoticeV2 struct {
	Notices []struct {
		Type         string    `json:"_type"`
		Messages     []Message `json:"messages"`
		Title        string    `json:"title"`
		Showonce     bool      `json:"showonce"`
		Platform     string    `json:"platform"`
		ActiveDate   string    `json:"_activeDate"`
		LastModified string    `json:"lastModified"`
		Locale       string    `json:"_locale"`
	} `json:"notices"`
	Title        string `json:"_title"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type CreativeAd struct {
	LinkURL      string `json:"linkURL"`
	ImageURL     string `json:"imageURL"`
	Title        string `json:"title"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type CreativeFeatures struct {
	CreativeFeatures []struct {
		Title        string `json:"title"`
		ImageURL     string `json:"imageURL"`
		LinkURL      string `json:"linkURL"`
		NoIndex      bool   `json:"_noIndex"`
		ActiveDate   string `json:"_activeDate"`
		LastModified string `json:"lastModified"`
		Locale       string `json:"_locale"`
	} `json:"creativeFeatures"`
	Title        string `json:"_title"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type SpecialOfferVideo struct {
	ShowVideo    bool   `json:"showVideo"`
	LinkURL      string `json:"linkURL"`
	Title        string `json:"title"`
	ImageURL     string `json:"imageURL"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type SubgameInfo struct {
	BattleRoyale struct {
		Overrideable Message           `json:"overrideablemessage"`
		Type         string            `json:"_type"`
		Title        map[string]string `json:"title"`
		ActiveDate   string            `json:"_activeDate"`
		LastModified string            `json:"lastModified"`
		Locale       string            `json:"_locale"`
	} `json:"battleRoyale"`
	SaveTheWorld struct {
		Overrideable Message           `json:"overrideablemessage"`
		Type         string            `json:"_type"`
		Title        map[string]string `json:"title"`
		ActiveDate   string            `json:"_activeDate"`
		LastModified string            `json:"lastModified"`
		Locale       string            `json:"_locale"`
	} `json:"saveTheWorld"`
	Creative struct {
		Overrideable Message           `json:"overrideablemessage"`
		Type         string            `json:"_type"`
		Title        map[string]string `json:"title"`
		ActiveDate   string            `json:"_activeDate"`
		LastModified string            `json:"lastModified"`
		Locale       string            `json:"_locale"`
	} `json:"creative"`
	Title        string `json:"_title"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type Lobby struct {
	Images []struct {
		ImageURL string `json:"image"`
		Type     string `json:"_type"`
	} `json:"images"`
	Title        string `json:"_title"`
	NoIndex      bool   `json:"_noIndex"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
}

type DynamicBackground struct {
	Stage string `json:"stage"`
	Type  string `json:"_type"`
	Key   string `json:"key"`
}

type DynamicBackgroundList struct {
	Backgrounds []DynamicBackground `json:"backgrounds"`
	Type        string              `json:"_type"`
}

type ShopSection struct {
	BSortOffersByOwnership          bool              `json:"bSortOffersByOwnership"`
	BShowIneligibleOffersIfGiftable bool              `json:"bShowIneligibleOffersIfGiftable"`
	BEnableToastNotification        bool              `json:"bEnableToastNotification"`
	Background                      DynamicBackground `json:"background"`
	Type                            string            `json:"_type"`
	LandingPriority                 int               `json:"landingPriority"`
	BHidden                         bool              `json:"bHidden"`
	SectionID                       string            `json:"sectionId"`
	BShowTimer                      bool              `json:"bShowTimer"`
	SectionDisplayName              string            `json:"sectionDisplayName"`
	BShowIneligibleOffers           bool              `json:"bShowIneligibleOffers"`
}

type ShopSectionList struct {
	Type     string        `json:"_type"`
	Sections []ShopSection `json:"sections"`
}

type Creativenews struct {
	News         BattleRoyaleNews `json:"news"`
	Title        string           `json:"_title"`
	Header       string           `json:"header"`
	Style        string           `json:"style"`
	NoIndex      bool             `json:"_noIndex"`
	AlwaysShow   bool             `json:"alwaysShow"`
	ActiveDate   string           `json:"_activeDate"`
	LastModified string           `json:"lastModified"`
	Locale       string           `json:"_locale"`
}

type BattleRoyaleNews struct {
	Type     string   `json:"_type"`
	Messages []string `json:"messages"`
}

type ContentPage struct {
	DynamicBackground     DynamicBackgroundList   `json:"dynamicbackgrounds"`
	ShopSectionList       ShopSectionList         `json:"shopSections"`
	Creativenews          Creativenews            `json:"creativenews"`
	BattleRoyaleNews      Creativenews            `json:"battleroyalenews"`
	Lobby                 Lobby                   `json:"lobby"`
	SubgameInfo           SubgameInfo             `json:"subgameinfo"`
	SpecialOfferVideo     SpecialOfferVideo       `json:"specialoffer"`
	CreativeFeatures      CreativeFeatures        `json:"creativeFeatures"`
	CreativeAd            CreativeAd              `json:"creativeAd"`
	EmergencyNoticeV2     EmergencyNoticeV2       `json:"emergencynotices"`
	EmergencyNotice       EmergencyNotice         `json:"emergencyNotices"`
	TournamentInformation TournamentInformation   `json:"tournament_info"`
	PlaylistInformation   PlaylistInformation     `json:"playlist_info"`
	BattlePassAbout       BattlePassAboutMessages `json:"battlepassabout"`
	SaveTheWorldNews      SaveTheWorldNews        `json:"news"`
	SubgameSelectData     SubgameSelectData       `json:"subgameselectdata"`
	AthenaMessage         AthenaMessage           `json:"athenamessage"`
	SurvivalMessage       SurvivalMessage         `json:"survivalmessage"`
	LoginMessage          LoginMessage            `json:"loginmessage"`
}