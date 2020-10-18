package broker

import (
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
)

type Workspace struct {
	name     string
	users    map[string]entity.User
	channels map[string]entity.Channel
}

func NewWorkspace(name string) Workspace {
	workspace := Workspace{
		name:     name,
		users:    make(map[string]entity.User),
		channels: make(map[string]entity.Channel),
	}
	workspace.channels["random"] = entity.Channel{Id: 0, Name: "random"}

	return workspace
}

func (workspace Workspace) RegisterNewUser(user entity.User) {
	workspace.users[user.NickName] = user
	randomChannel := workspace.channels["random"]
	randomChannel.Subscribers = append(randomChannel.Subscribers, user)
	workspace.channels["random"] = randomChannel
}

func (workspace Workspace) CreateChannel(sender models.User, name string) error {
	user := workspace.users[sender.NickName]
	if channel, exist := workspace.channels[name]; exist {
		return entity.ChannelAlreadyExist{Name: channel.Name}
	}

	channel := entity.Channel{Id: len(workspace.channels), Name: name}
	channel.Subscribers = append(channel.Subscribers, user)
	workspace.channels[name] = channel

	return nil
}

func (workspace Workspace) WorkspaceChannels() []string {
	var channels []string
	for _, channel := range workspace.channels {
		channels = append(channels, channel.Name)
	}

	return channels
}

func (workspace Workspace) WorkspaceUsers() []string {
	var users []string
	for _, user := range workspace.users {
		users = append(users, user.NickName)
	}

	return users
}

func (workspace Workspace) ChannelSubscribers(channel entity.Channel) []string {
	var users []string
	for _, user := range channel.Subscribers {
		users = append(users, user.NickName)
	}

	return users
}

func (workspace Workspace) AddUserToChannel(sender models.User, name string) error {
	user := workspace.users[sender.NickName]
	if channel, exist := workspace.channels[name]; exist {
		if channel.Contains(user) {
			return entity.ChannelAlreadyJoined{Name: channel.Name}
		}

		channel.Subscribers = append(channel.Subscribers, user)
		workspace.channels[name] = channel
	} else {
		return entity.ChannelNotExist{Name: name}
	}

	return nil
}

func (workspace Workspace) RemoveUserFromChannel(sender models.User, name string) error {
	user := workspace.users[sender.NickName]
	if channel, exist := workspace.channels[name]; exist {
		if isPresent, index := channel.ContainsSubscriber(user); isPresent {
			channel.Subscribers = append(channel.Subscribers[:index], channel.Subscribers[index+1:]...)
		} else {
			return entity.ChannelNotJoined{Name: channel.Name}
		}
	} else {
		return entity.ChannelNotExist{Name: name}
	}

	return nil
}
