package broker

import "github.com/ASV44/ChatMessageBroker/broker/entity"

type Channel struct {
	Id          int
	Name        string
	Subscribers []entity.User
}

func (channel *Channel) Contains(user entity.User) bool {
	for _, subscriber := range channel.Subscribers {
		if subscriber == user && subscriber.ID == user.ID {
			return true
		}
	}
	return false
}

func (channel *Channel) ContainsSubscriber(user entity.User) (bool, int) {
	for index, subscriber := range channel.Subscribers {
		if subscriber == user && subscriber.ID == user.ID {
			return true, index
		}
	}
	return false, -1
}
