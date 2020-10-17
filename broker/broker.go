package broker

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ASV44/ChatMessageBroker/broker/components"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
)

// Broker represents main structure which contains all related fields for routing message
type Broker struct {
	workspace string
	server    broker.Server
	incoming  chan models.IncomingMessage
	users     map[string]entity.User
	channels  map[string]entity.Channel
}

// Init creates and initialize Broker instance
func Init() Broker {
	brokerInstance := Broker{workspace: "Matrix"}
	brokerInstance.server = broker.InitServer(broker.DefaultHost, broker.DefaultPort, broker.DefaultType)
	brokerInstance.incoming = make(chan models.IncomingMessage)
	brokerInstance.users = make(map[string]entity.User)
	brokerInstance.channels = make(map[string]entity.Channel)
	brokerInstance.channels["random"] = entity.Channel{Id: 0, Name: "random"}

	return brokerInstance
}

// Start init broker server, creates channels and start receiving and routing of connections
func (broker Broker) Start() error {
	err := broker.server.Start()
	if err != nil {
		return err
	}

	broker.run()

	return nil
}

func (broker Broker) listenIncomingMessages(connection broker.Connection) {
	var message models.IncomingMessage
	for {
		if err := connection.GetMessage(&message); err != io.EOF {
			broker.incoming <- message
		} else {
			return
		}
	}
}

func (broker Broker) run() {
	for {
		select {
		case connection := <-broker.server.Connection:
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

func (broker Broker) register(connection broker.Connection) {
	text := fmt.Sprintf("Welcome to %s workspace!\nEnter nickname:", broker.workspace)
	message := models.Register{UserId: len(broker.users), Text: text, Time: time.Now()}
	err := connection.SendMessage(message)
	if err != nil {
		fmt.Println("Could not send register data ", err)
		broker.close(connection)
	}

	var user models.User
	err = connection.GetMessage(&user)
	newUser := user.ToUserEntity(connection)

	broker.users[newUser.NickName] = newUser
	randomChannel := broker.channels["random"]
	randomChannel.Subscribers = append(randomChannel.Subscribers, newUser)
	broker.channels["random"] = randomChannel
	broker.show(user, "all")

	go broker.listenIncomingMessages(connection)
	fmt.Printf("Connected user: %s Id: %d addrr: %v\n", newUser.NickName, newUser.ID, connection.RemoteAddr())
}

func (broker Broker) close(connection broker.Connection) {
	err := connection.Close()
	if err != nil {
		fmt.Println("Error during closing connection ", err)
	}
}

func (broker Broker) sendMessageToUser(user entity.User, message models.OutgoingMessage) {
	err := user.Connection.SendMessage(message)
	if err != nil {
		fmt.Printf("Failed to send message to user ID:%s Nickname: %s. Error: %s\n", user.ID, user.NickName, err)
	}
}

func (broker Broker) handleDirectMessage(message models.IncomingMessage) {
	user, isPresent := broker.users[message.Target]
	if isPresent {
		_ = user.Connection.SendMessage(message.ToOutgoingMessage())
	}
}

func (broker Broker) handleChannelMessage(message models.IncomingMessage) {
	if channel, isPresent := broker.channels[message.Target]; isPresent {
		sender := broker.users[message.Sender.NickName]
		if channel.Contains(sender) {
			draft := message.ToOutgoingMessage()
			for _, user := range channel.Subscribers {
				if user.ID != message.Sender.ID {
					go broker.sendMessageToUser(user, draft)
				}
			}
		} else {
			broker.sendMessageToUser(sender, models.OutgoingMessage{Channel: channel.Name, Text: "not joined yet!"})
		}
	}
}

func (broker Broker) createChannel(sender models.User, name string) {
	user := broker.users[sender.NickName]
	_, exist := broker.channels[name]
	if !exist {
		channel := entity.Channel{Id: len(broker.channels), Name: name}
		channel.Subscribers = append(channel.Subscribers, user)
		broker.channels[name] = channel
		broker.sendMessageToUser(user, broker.getWorkspaceChannelsMessage())
	} else {
		broker.sendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: "already exist!"})
	}
}

func (broker Broker) joinChannel(sender models.User, name string) {
	user := broker.users[sender.NickName]
	if channel, exist := broker.channels[name]; exist {
		if !channel.Contains(user) {
			channel.Subscribers = append(channel.Subscribers, user)
			broker.channels[name] = channel
			broker.sendMessageToUser(user, broker.getChannelSubscribersMessage(channel))
		} else {
			broker.sendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: "already joined!"})
		}
	} else {
		broker.sendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: "does not exist!"})
	}
}

func (broker Broker) leaveChannel(sender models.User, name string) {
	user := broker.users[sender.NickName]
	channel, exist := broker.channels[name]
	if exist {
		if isPresent, index := channel.ContainsSubscriber(user); isPresent {
			channel.Subscribers = append(channel.Subscribers[:index], channel.Subscribers[index+1:]...)
			broker.sendMessageToUser(user, broker.getChannelSubscribersMessage(channel))
		} else {
			broker.sendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: "not subscribed to!"})
		}
	} else {
		broker.sendMessageToUser(user, models.OutgoingMessage{Channel: name, Text: "does not exist!"})
	}
}

func (broker Broker) getWorkspaceChannelsMessage() models.OutgoingMessage {
	var channels []string
	for _, channel := range broker.channels {
		channels = append(channels, channel.Name)
	}

	return models.OutgoingMessage{Channel: "Channels",
		Text: strings.Join(channels, " | "),
		Time: time.Now()}
}

func (broker Broker) getWorkspaceUsersMessage() models.OutgoingMessage {
	var users []string
	for _, user := range broker.users {
		users = append(users, user.NickName)
	}
	return models.OutgoingMessage{Sender: "Users",
		Text: strings.Join(users, " | "),
		Time: time.Now()}

}

func (broker Broker) getChannelSubscribersMessage(channel entity.Channel) models.OutgoingMessage {
	var users []string
	for _, user := range channel.Subscribers {
		users = append(users, user.NickName)
	}
	return models.OutgoingMessage{Channel: channel.Name,
		Sender: "Users",
		Text:   strings.Join(users, " | "),
		Time:   time.Now()}
}

func (broker Broker) show(sender models.User, param string) {
	user := broker.users[sender.NickName]
	switch param {
	case "users":
		broker.sendMessageToUser(user, broker.getWorkspaceUsersMessage())
	case "channels":
		broker.sendMessageToUser(user, broker.getWorkspaceChannelsMessage())
	case "all":
		broker.sendMessageToUser(user, broker.getWorkspaceChannelsMessage())
		broker.sendMessageToUser(user, broker.getWorkspaceUsersMessage())
	default: // Get all users of specific channel if it exist
		if channel, exist := broker.channels[param]; exist {
			broker.sendMessageToUser(user, broker.getChannelSubscribersMessage(channel))
		}
	}
}
