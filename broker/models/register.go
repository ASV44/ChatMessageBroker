package models

import (
	"time"

	"github.com/ASV44/chat-message-broker/common"

	"github.com/ASV44/chat-message-broker/broker/entity"
)

// Register represents model of message which is sent to user at connection for register
type Register struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

// AccountData represents model of message which is sent from user with all account data required at sign up
type AccountData struct {
	NickName string `json:"nickName"`
	Password string `json:"password"`
}

// ToRegistrationData map AccountData model to entity.RegistrationData entity
func (accountData AccountData) ToRegistrationData(conn common.Connection) entity.RegistrationData {
	return entity.RegistrationData{
		NickName:   accountData.NickName,
		Connection: conn,
	}
}
