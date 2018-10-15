package sender

import (
	"ChatMessageBroker/client/models"
	"time"
)

type Message struct {
	ChannelId int         `json:"channelId"`
	Sender    models.User `json:"user"`
	Text      string      `json:"text"`
	Time      time.Time   `json:"time"`
}
