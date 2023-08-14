package models

type V2Party struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Config    struct {
		Type             string `json:"type"`
		Joinability      string `json:"joinability"`
		SubType          string `json:"sub_type"`
		MaxSize          int    `json:"max_size"`
		Discoverability  string `json:"discoverability"`
		InviteTtl        int    `json:"invite_ttl"`
		JoinConfirmation string `json:"join_confirmation"`
	} `json:"config"`
	Meta       map[string]interface{} `json:"meta"`
	Members    []V2PartyMember        `json:"members"`
	Invites    []any                  `json:"invites"`
	Intentions []any                  `json:"intentions"`
	Revision   int                    `json:"revision"`
}

type V2PartyMember struct {
	AccountId   string                 `json:"account_id"`
	Meta        map[string]interface{} `json:"meta"`
	Connections []V2PartyConnection    `json:"connections"`
	Role        string                 `json:"role"`
	Revision    int                    `json:"revision"`
	JoinedAt    string                 `json:"joined_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

type V2PartyConnection struct {
	ID              string                 `json:"id"`
	ConnectedAt     string                 `json:"connected_at"`
	UpdatedAt       string                 `json:"updated_at"`
	YieldLeadership bool                   `json:"yield_leadership"`
	Meta            map[string]interface{} `json:"meta"`
}