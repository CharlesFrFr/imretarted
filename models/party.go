package models

type V2Party struct {
	ID         string                 `json:"id"`
	CreatedAt  string                 `json:"created_at"`
	UpdatedAt  string                 `json:"updated_at"`
	Config     V2PartyConfig          `json:"config"`
	Meta       map[string]interface{} `json:"meta"`
	Members    []V2PartyMember        `json:"members"`
	Invites    []any                  `json:"invites"`
	Intentions []any                  `json:"intentions"`
	Revision   int                    `json:"revision"`
}

type V2PartyConfig struct {
	Type             string `json:"type"`
	Joinability      string `json:"joinability"`
	SubType          string `json:"sub_type"`
	MaxSize          int    `json:"max_size"`
	Discoverability  string `json:"discoverability"`
	InviteTtl        int    `json:"invite_ttl"`
	JoinConfirmation bool   `json:"join_confirmation"`
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

/*captain.Meta["urn:epic:member:joinrequestusers_j"] = gin.H{}*./*/

type V2CaptainJoinRequestUsers struct {
	Users []V2CaptainJoinRequestUser `json:"users"`
}

type V2CaptainJoinRequestUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"dn"`
	Platform    string `json:"plat"`
	Data        string `json:"data"`
}

// 	'{"RawSquadAssignments":[{"memberId":"39764163-6510-4915-8ec1-2aedb3ce6c94","absoluteMemberIdx":0}]}'

type V2RawSquadAssignments struct {
	RawSquadAssignments []V2RawSquadAssignment `json:"RawSquadAssignments"`
}

type V2RawSquadAssignment struct {
	MemberId          string `json:"memberId"`
	AbsoluteMemberIdx int    `json:"absoluteMemberIdx"`
}