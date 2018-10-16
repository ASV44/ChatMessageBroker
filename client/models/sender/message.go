package sender

import (
	"ChatMessageBroker/client/models"
	"time"
)

const (
	CMD     = "cmd"
	DIRECT  = "direct"
	CHANNEL = "channel"
)

type Message struct {
	Type   string      `json:"type"`
	Target string      `json:"target"`
	Sender models.User `json:"user"`
	Text   string      `json:"text"`
	Time   time.Time   `json:"time"`
}
