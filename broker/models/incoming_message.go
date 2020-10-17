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

func (incoming IncomingMessage) ToOutgoingMessage() OutgoingMessage {
	var channel string
	if incoming.Type == CHANNEL {
		channel = incoming.Target
	}
	return OutgoingMessage{Channel: channel, Sender: incoming.Sender.NickName,
		Text: incoming.Text, Time: incoming.Time}
}
