package models

import (
	"ChatMessageBroker/broker/entity"
	"net"
)

type User struct {
	Id       int    `json:"id"`
	NickName string `json:"nickName"`
}

func (user User) ToUserEntity(connection net.Conn) entity.User {
	return entity.User{Id: user.Id, NickName: user.NickName, Connection: connection}
}
