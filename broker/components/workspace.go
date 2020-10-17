package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"time"
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

func (workspace Workspace) RegisterNewUser(connection Connection) error {
	text := fmt.Sprintf("Welcome to %s workspace!\nEnter nickname:", workspace.name)
	message := models.Register{UserId: len(workspace.users), Text: text, Time: time.Now()}
	err := connection.SendMessage(message)
	if err != nil {
		fmt.Println("Could not send register data ", err)
		return err
	}

	var user models.User
	err = connection.GetMessage(&user)
	if err != nil {
		fmt.Println("Could not receive register data from user ", err)
		return err
	}
	newUser := user.ToUserEntity(connection)

	workspace.users[newUser.NickName] = newUser
	randomChannel := workspace.channels["random"]
	randomChannel.Subscribers = append(randomChannel.Subscribers, newUser)
	workspace.channels["random"] = randomChannel
	//TODO: Add sending of all users and channels at registration

	fmt.Printf("Connected user: %s Id: %d addrr: %v\n", newUser.NickName, newUser.ID, connection.RemoteAddr())

	return nil
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
