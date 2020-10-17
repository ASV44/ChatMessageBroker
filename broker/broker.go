package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker/components"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"io"
)

// Broker represents main structure which contains all related fields for routing message
type Broker struct {
	workspace  broker.Workspace
	server     broker.Server
	incoming   chan models.IncomingMessage
	dispatcher broker.MessageDispatcher
}

// Init creates and initialize Broker instance
func Init() Broker {
	workspace := broker.NewWorkspace("Matrix")
	return Broker{
		workspace:  workspace,
		server:     broker.InitServer(broker.DefaultHost, broker.DefaultPort, broker.DefaultType),
		incoming:   make(chan models.IncomingMessage),
		dispatcher: broker.NewMessageDispatcher(&workspace),
	}
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
			go broker.dispatcher.DispatchMessage(message)
		}
	}
}

func (broker Broker) register(connection broker.Connection) {
	err := broker.workspace.RegisterNewUser(connection)
	switch err.(type) {
	case nil:
		broker.listenIncomingMessages(connection)
	default:
		fmt.Println("Register of new user failed ", err)
		broker.close(connection)
	}
}

func (broker Broker) close(connection broker.Connection) {
	err := connection.Close()
	if err != nil {
		fmt.Println("Error during closing connection ", err)
	}
}

func SendMessageToUser(user entity.User, message models.OutgoingMessage) {
	err := user.Connection.SendMessage(message)
	if err != nil {
		fmt.Printf("Failed to send message to user ID:%s Nickname: %s. Error: %s\n", user.ID, user.NickName, err)
	}
}
