package models

type Ping struct {
	SentBy    string                 `json:"sent_by"`
	SentTo    string                 `json:"sent_to"`
	SentAt    string                 `json:"sent_at"`
	ExpiresAt string                 `json:"expires_at"`
	Meta      map[string]interface{} `json:"meta"`
}