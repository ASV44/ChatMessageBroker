package models

import (
	broker "github.com/ASV44/ChatMessageBroker/broker/components"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
)

// User represents model with user data received from client in JSON format
type User struct {
	ID       int    `json:"id"`
	NickName string `json:"nickName"`
}

// ToUserEntity map user model to user entity
func (user User) ToUserEntity(connection broker.Connection) entity.User {
	return entity.User{ID: user.ID, NickName: user.NickName, Connection: connection}
}
