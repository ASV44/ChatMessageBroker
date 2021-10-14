package entity

// Channel represents entity of workspace channel
type Channel struct {
	ID          int
	Name        string
	Subscribers []User
}

// Contains check if user is part of specific channel
func (channel Channel) Contains(user User) bool {
	for _, subscriber := range channel.Subscribers {
		if subscriber.NickName == user.NickName && subscriber.ID == user.ID {
			return true
		}
	}

	return false
}

// ContainsSubscriber check if user is part of specific channel and returns index of that user in channel
func (channel Channel) ContainsSubscriber(user User) (bool, int) {
	for index, subscriber := range channel.Subscribers {
		if subscriber.NickName == user.NickName && subscriber.ID == user.ID {
			return true, index
		}
	}

	return false, -1
}
