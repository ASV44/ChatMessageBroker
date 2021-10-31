package entity

import (
	"github.com/ASV44/chat-message-broker/common"
)

// User represents entity of broker user
type User struct {
	ID           int
	NickName     string
	PasswordHash string
	Connection   common.Connection
	Channels     []string
}

func (user User) IsSubscribedToChannel(name string) bool {
	for _, channel := range user.Channels {
		if channel == name {
			return true
		}
	}

	return false
}
