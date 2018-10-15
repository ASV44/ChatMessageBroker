package receiver

import (
	"ChatMessageBroker/client/models"
	"time"
)

const (
	SYSTEM        = "system"
	COMMUNICATION = "communication"
)

type Message struct {
	Type   string      `json:"type"`
	Sender models.User `json:"user"`
	Text   string      `json:"text"`
	Time   time.Time   `json:"time"`
}
