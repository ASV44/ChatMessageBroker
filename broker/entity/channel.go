package entity

type Channel struct {
	Id          int
	Name        string
	Subscribers []User
}

func (channel *Channel) Contains(user User) bool {
	for _, subscriber := range channel.Subscribers {
		if subscriber == user && subscriber.ID == user.ID {
			return true
		}
	}
	return false
}

func (channel *Channel) ContainsSubscriber(user User) (bool, int) {
	for index, subscriber := range channel.Subscribers {
		if subscriber == user && subscriber.ID == user.ID {
			return true, index
		}
	}
	return false, -1
}
