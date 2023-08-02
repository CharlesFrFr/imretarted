package models

type ActiveEvent struct {
	EventType   string `json:"eventType"`
	ActiveUntil string `json:"activeUntil"`
	ActiveSince string `json:"activeSince"`
}

type State struct {
	ValidFrom    string        `json:"validFrom"`
	ActiveEvents []ActiveEvent `json:"activeEvents"`
	State        EventState    `json:"state"`
}

type EventState struct {
	ActiveStorefronts  []interface{}          `json:"activeStorefronts"`
	EventNamedWeights  map[string]interface{} `json:"eventNamedWeights"`
	SeasonNumber       int                    `json:"seasonNumber"`
	SeasonTemplateID   string                 `json:"seasonTemplateId"`
	MatchXPBonusPoints int                    `json:"matchXpBonusPoints"`
	SeasonBegin        string                 `json:"seasonBegin"`
	SeasonEnd          string                 `json:"seasonEnd"`
	SeasonDisplayedEnd string                 `json:"seasonDisplayedEnd"`
	WeeklyStoreEnd     string                 `json:"weeklyStoreEnd"`
	STWEventStoreEnd   string                 `json:"stwEventStoreEnd"`
	STWWeeklyStoreEnd  string                 `json:"stwWeeklyStoreEnd"`
	SectionStoreEnds   map[string]string      `json:"sectionStoreEnds"`
	DailyStoreEnd      string                 `json:"dailyStoreEnd"`
}

type Channels struct {
	ClientMatchmaking Channel `json:"client-matchmaking"`
	ClientEvents      Channel `json:"client-events"`
}

type Channel struct {
	States      []State `json:"states"`
	CacheExpire string  `json:"cacheExpire"`
}

type FortniteTimeCalendar struct {
	Channels            Channels `json:"channels"`
	EventsTimeOffsetHrs int      `json:"eventsTimeOffsetHrs"`
	CacheIntervalMins   int      `json:"cacheIntervalMins"`
	CurrentTime         string   `json:"currentTime"`
}