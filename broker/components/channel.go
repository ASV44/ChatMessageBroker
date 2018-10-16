package broker

import "ChatMessageBroker/broker/entity"

type Channel struct {
	Id          int
	Name        string
	Subscribers []entity.User
}

func (channel *Channel) Contains(user entity.User) bool {
	for _, subscriber := range channel.Subscribers {
		if subscriber == user && subscriber.Id == user.Id {
			return true
		}
	}
	return false
}
