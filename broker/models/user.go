package models

import (
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/common"
)

// User represents model with user data received from client in JSON format
type User struct {
	ID       int    `json:"ID"`
	NickName string `json:"nickName"`
}

// ToUserEntity map user model to user entity
func (user User) ToUserEntity(connection common.Connection) entity.User {
	return entity.User{ID: user.ID, NickName: user.NickName, Connection: connection}
}
