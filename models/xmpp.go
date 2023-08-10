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