package models

import "time"

const (
	CMD     = "cmd"
	DIRECT  = "direct"
	CHANNEL = "channel"
)

type IncomingMessage struct {
	Type   string    `json:"type"`
	Target string    `json:"target"`
	Sender User      `json:"user"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}
