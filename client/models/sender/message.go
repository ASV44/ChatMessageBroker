package sender

import (
	"time"

	"github.com/ASV44/ChatMessageBroker/client/models"
)

// Constant values of message types
const (
	CMD     = "cmd"
	DIRECT  = "direct"
	CHANNEL = "channel"
)

// Message represents model with data of user message sent to broker
type Message struct {
	Type   string      `json:"type"`
	Target string      `json:"target"`
	Sender models.User `json:"user"`
	Text   string      `json:"text"`
	Time   time.Time   `json:"time"`
}
