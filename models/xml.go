package models

import "encoding/xml"

type OpenXML struct {
	XMLName xml.Name `xml:"open"`
	To      string   `xml:"to,attr"`
	Version string   `xml:"version,attr"`
}

type StreamFeaturesNeedAuth struct {
	XMLName     xml.Name      `xml:"stream:features"`
	Stream      string        `xml:"xmlns:stream,attr"`
	Mechanisms  Mechanisms    `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Ver         Ver           `xml:"urn:xmpp:features:rosterver ver"`
	StartTLS    StartTLS      `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Compression Compression   `xml:"http://jabber.org/features/compress compression"`
	Auth        Auth          `xml:"http://jabber.org/features/iq-auth auth"`
}

type Mechanisms struct {
	XMLName   xml.Name   `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanism string     `xml:"mechanism"`
}

type Ver struct {
	XMLName xml.Name `xml:"urn:xmpp:features:rosterver ver"`
}

type StartTLS struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
}

type Compression struct {
	XMLName xml.Name `xml:"http://jabber.org/features/compress compression"`
	Method  string   `xml:"method"`
}

type Auth struct {
	XMLName xml.Name `xml:"http://jabber.org/features/iq-auth auth"`
}

type StreamFeatures struct {
	XMLName     xml.Name     `xml:"stream:features"`
	Stream      string       `xml:"xmlns:stream,attr"`
	Ver         Ver          `xml:"ver"`
	StartTLS    StartTLS     `xml:"starttls"`
	Bind        Bind         `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Compression Compression  `xml:"compression"`
	Session     Session      `xml:"urn:ietf:params:xml:ns:xmpp-session session"`
}

type Bind struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
}

type Session struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-session session"`
}

type AuthXML struct {
	XMLName   xml.Name `xml:"auth"`
	Mechanism string   `xml:"mechanism,attr"`
	XMLNS     string   `xml:"xmlns,attr"`
	Value     string   `xml:",chardata"`
}

type CloseXML struct {
	XMLName xml.Name `xml:"close"`
	XMLNS   string   `xml:"xmlns,attr"`
}

type BindIQXML struct {
	XMLName xml.Name `xml:"iq"`
	ID      string   `xml:"id,attr"`
	Type    string   `xml:"type,attr"`
	Bind    IQBind     `xml:"bind"`
}

type IQBind struct {
	XMLName  xml.Name `xml:"bind"`
	XMLNS    string   `xml:"xmlns,attr"`
	Resource string   `xml:"resource"`
	Value    string   `xml:",chardata"`
}

type SessionIQXML struct {
	XMLName xml.Name `xml:"iq"`
	ID      string   `xml:"id,attr"`
	Type    string   `xml:"type,attr"`
	Session IQSession  `xml:"session"`
}

type IQSession struct {
	XMLName xml.Name `xml:"session"`
	XMLNS   string   `xml:"xmlns,attr"`
}

type MessageWithBodyXML struct {
	XMLName xml.Name `xml:"message"`
	ID      string   `xml:"id,attr"`
	To      string   `xml:"to,attr"`
	Type    string   `xml:"type,attr"`
	Body    MessageBody     `xml:"body"`
}

type MessageBody struct {
	XMLName xml.Name `xml:"body"`
	Value   string   `xml:",chardata"`
}


type MessageBodyPayload struct {
	Type      string `json:"type"`
	Payload   Party  `json:"payload"`
	Timestamp string `json:"timestamp"`
}

type Party struct {
	PartyID string `json:"partyId"`
}

type PresenceXML struct {
	XMLName xml.Name `xml:"presence"`
	Status  Status   `xml:"status"`
	Delay   Delay    `xml:"delay"`
}

type Status struct {
	Value string `xml:",chardata"`
}

type Delay struct {
	Stamp string `xml:"stamp,attr"`
	XMLNS string `xml:"xmlns,attr"`
}

type FriendProperties struct {
	FortBasicInfo  FortBasicInfo    `json:"FortBasicInfo_j"`
	FortGameplayStats FortGameplayStats `json:"FortGameplayStats_j"`
	FortLFG        string           `json:"FortLFG_I"`
	FortPartySize  int              `json:"FortPartySize_i"`
	FortSubGame    int              `json:"FortSubGame_i"`
	InUnjoinableMatch bool           `json:"InUnjoinableMatch_b"`
	JoinInfoData   JoinInfoData     `json:"party.joininfodata.286331153_j"`
}

type FortBasicInfo struct {
	HomeBaseRating int `json:"homeBaseRating"`
}

type FortGameplayStats struct {
	BFellToDeath bool   `json:"bFellToDeath"`
	NumKills     int    `json:"numKills"`
	Playlist     string `json:"playlist"`
	State        string `json:"state"`
}

type JoinInfoData struct {
	AppId            string `json:"appId"`
	BuildId          string `json:"buildId"`
	Key              string `json:"key"`
	NotAcceptingReason int   `json:"notAcceptingReason"`
	PartyFlags       int    `json:"partyFlags"`
	PartyId          string `json:"partyId"`
	PartyTypeId      int    `json:"partyTypeId"`
	Pc               int    `json:"pc"`
	SourceDisplayName string `json:"sourceDisplayName"`
	SourceId         string `json:"sourceId"`
	SourcePlatform   string `json:"sourcePlatform"`
}

type StatusJSON struct {
	Properties FriendProperties `json:"Properties"`
	SessionId  string          `json:"SessionId"`
	Status     string          `json:"Status"`
	BHasVoiceSupport bool      `json:"bHasVoiceSupport"`
	BIsJoinable     bool      `json:"bIsJoinable"`
	BIsPlaying      bool      `json:"bIsPlaying"`
}

/*<presence to="Party-2ad8b220-4f13-4456-a0a3-cbde2bcbfcfd@muc.prod.ol.epicgames.com/admin:571f16e7-c6aa-41f5-b24c-edc70fc88406:V2:Fortnite:WIN::E0EB415645D78EC5C252798418B1548A"><x xmlns="http://jabber.org/protocol/muc"><history maxstanzas="50"/></x></presence>*/

type PartyPresenceXML struct {
	XMLName xml.Name `xml:"presence"`
	Type    string   `xml:"type,attr"`
	To      string   `xml:"to,attr"`
	X       X        `xml:"x"`
}

type X struct {
	XMLName xml.Name `xml:"x"`
	XMLNS   string   `xml:"xmlns,attr"`
	History History  `xml:"history"`
}

type History struct {
	XMLName     xml.Name `xml:"history"`
	Maxstanzas  string   `xml:"maxstanzas,attr"`
}