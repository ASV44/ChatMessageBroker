package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/ASV44/ChatMessageBroker/broker/components"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
)

// Broker represents main structure which contains all related fields for routing message
type Broker struct {
	workspace string
	server    *broker.Server
	incoming  chan models.IncomingMessage
	users     map[string]entity.User
	channels  map[string]*entity.Channel
}

func Init() Broker {
	brokerInstance := Broker{workspace: "Matrix"}
	brokerInstance.server = &broker.Server{
		Host:           broker.DefaultHost,
		Port:           broker.DefaultPort,
		ConnectionType: broker.DefaultType,
	}
	brokerInstance.incoming = make(chan models.IncomingMessage)
	brokerInstance.users = make(map[string]entity.User)
	brokerInstance.channels = make(map[string]*entity.Channel)
	brokerInstance.channels["random"] = &entity.Channel{Id: 0, Name: "random"}

	return brokerInstance
}

// Start init broker server, creates channels and start receiving and routing of connections
func (broker Broker) Start() {
	broker.server.Start()
	defer broker.server.Close()

	broker.run()
}

func (broker Broker) listen(connection net.Conn) {
	decoder := json.NewDecoder(connection)
	var message models.IncomingMessage
	for {
		if err := decoder.Decode(&message); err != io.EOF {
			broker.incoming <- message
		} else {
			return
		}
	}
}

func (broker Broker) run() {
	for {
		select {
		case connection := <-broker.server.Connections:
			go broker.register(connection)
		case message := <-broker.incoming:
			go broker.dispatchMessage(message)
		}
	}
}

func (broker Broker) dispatchMessage(message models.IncomingMessage) {
	switch message.Type {
	case models.CMD:
		broker.dispatchCommand(message)
	case models.DIRECT:
		broker.handleDirectMessage(message)
	case models.CHANNEL:
		broker.handleChannelMessage(message)
	}
}

func (broker Broker) dispatchCommand(message models.IncomingMessage) {
	switch message.Target {
	case "create":
		broker.createChannel(message.Sender, message.Text)
	case "join":
		broker.joinChannel(message.Sender, message.Text)
	case "leave":
		broker.leaveChannel(message.Sender, message.Text)
	case "show":
		broker.show(message.Sender, message.Text)
	}
}

func (broker Broker) register(connection net.Conn) {
	text := fmt.Sprintf("Welcome to %s workspace!\nEnter nickname:", broker.workspace)
	message := models.Register{UserId: len(broker.users), Text: text, Time: time.Now()}
	data, _ := json.Marshal(message)
	_, err := connection.Write(data)
	if err != nil {
		fmt.Println("Could not write data at register ", err)
	}

	decoder := json.NewDecoder(connection)
	var user models.User
	err = decoder.Decode(&user)
	if err != nil {
		fmt.Println("Could not decode register response ", err)
	}
	newUser := user.ToUserEntity(connection)

	broker.users[newUser.NickName] = newUser
	subscribers := broker.channels["random"].Subscribers
	broker.channels["random"].Subscribers = append(subscribers, newUser)
	broker.show(user, "all")

	go broker.listen(connection)
	fmt.Printf("Connected user: %s Id: %d addrr: %v\n", newUser.NickName, newUser.ID, connection.RemoteAddr())
}

func (broker Broker) sendMessage(connection net.Conn, message models.OutcomingMessage) {
	data, _ := json.Marshal(message)
	_, err := connection.Write(data)
	if err != nil {
		fmt.Println("Could not write message data ", err)
	}
}

func (broker Broker) handleDirectMessage(message models.IncomingMessage) {
	user, isPresent := broker.users[message.Target]
	if isPresent {
		broker.sendMessage(user.Connection, message.ToOutcomingMessage())
	}
}

func (broker Broker) handleChannelMessage(message models.IncomingMessage) {
	channel, isPresent := broker.channels[message.Target]
	if isPresent {
		sender := broker.users[message.Sender.NickName]
		if channel.Contains(sender) {
			draft := message.ToOutcomingMessage()
			for _, user := range channel.Subscribers {
				if user.ID != message.Sender.ID {
					go broker.sendMessage(user.Connection, draft)
				}
			}
		} else {
			broker.sendMessage(sender.Connection, models.OutcomingMessage{Channel: channel.Name, Text: "not joined yet!"})
		}
	}
}

func (broker Broker) createChannel(sender models.User, name string) {
	user := broker.users[sender.NickName]
	_, exist := broker.channels[name]
	if !exist {
		channel := &entity.Channel{Id: len(broker.channels), Name: name}
		broker.channels[name] = channel
		channel.Subscribers = append(channel.Subscribers, user)
		broker.sendMessage(user.Connection, broker.getWorkspaceChannelsMessage())
	} else {
		broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "already exist!"})
	}
}

func (broker Broker) joinChannel(sender models.User, name string) {
	user := broker.users[sender.NickName]
	channel, exist := broker.channels[name]
	if exist {
		if !channel.Contains(user) {
			channel.Subscribers = append(channel.Subscribers, user)
			broker.sendMessage(user.Connection, broker.getChannelSubscribersMessage(channel))
		} else {
			broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "already joined!"})
		}
	} else {
		broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "does not exist!"})
	}
}

func (broker Broker) leaveChannel(sender models.User, name string) {
	user := broker.users[sender.NickName]
	channel, exist := broker.channels[name]
	if exist {
		if isPresent, index := channel.ContainsSubscriber(user); isPresent {
			channel.Subscribers = append(channel.Subscribers[:index], channel.Subscribers[index+1:]...)
			broker.sendMessage(user.Connection, broker.getChannelSubscribersMessage(channel))
		} else {
			broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "not subscribed to!"})
		}
	} else {
		broker.sendMessage(user.Connection, models.OutcomingMessage{Channel: name, Text: "does not exist!"})
	}
}

func (broker Broker) getWorkspaceChannelsMessage() models.OutcomingMessage {
	var channels []string
	for _, channel := range broker.channels {
		channels = append(channels, channel.Name)
	}

	return models.OutcomingMessage{Channel: "Channels",
		Text: strings.Join(channels, " | "),
		Time: time.Now()}
}

func (broker Broker) getWorkspaceUsersMessage() models.OutcomingMessage {
	var users []string
	for _, user := range broker.users {
		users = append(users, user.NickName)
	}
	return models.OutcomingMessage{Sender: "Users",
		Text: strings.Join(users, " | "),
		Time: time.Now()}

}

func (broker Broker) getChannelSubscribersMessage(channel *entity.Channel) models.OutcomingMessage {
	var users []string
	for _, user := range channel.Subscribers {
		users = append(users, user.NickName)
	}
	return models.OutcomingMessage{Channel: channel.Name,
		Sender: "Users",
		Text:   strings.Join(users, " | "),
		Time:   time.Now()}
}

func (broker Broker) show(sender models.User, param string) {
	connection := broker.users[sender.NickName].Connection
	switch param {
	case "users":
		broker.sendMessage(connection, broker.getWorkspaceUsersMessage())
	case "channels":
		broker.sendMessage(connection, broker.getWorkspaceChannelsMessage())
	case "all":
		broker.sendMessage(connection, broker.getWorkspaceChannelsMessage())
		broker.sendMessage(connection, broker.getWorkspaceUsersMessage())
	default: // Get all users of specific channel if it exist
		if channel, exist := broker.channels[param]; exist {
			broker.sendMessage(connection, broker.getChannelSubscribersMessage(channel))
		}
	}
}
