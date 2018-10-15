package broker

import "ChatMessageBroker/broker/entity"

type Channel struct {
	Id          int
	Name        string
	Subscribers []entity.User
}
