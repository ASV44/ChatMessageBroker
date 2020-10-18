package receiver

import (
	"time"
)

// Message represents model of message received from broker
type Message struct {
	Channel string    `json:"channel"`
	Sender  string    `json:"user"`
	Text    string    `json:"text"`
	Time    time.Time `json:"time"`
}
