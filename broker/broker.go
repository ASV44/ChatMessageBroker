package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker/components"
	"github.com/ASV44/ChatMessageBroker/broker/config"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/models"
	"github.com/ASV44/ChatMessageBroker/broker/services"
	"github.com/ASV44/ChatMessageBroker/common"
	"io"
)

// Broker represents main structure which contains all related fields for routing message
type Broker struct {
	workspace  broker.Workspace
	server     broker.Server
	incoming   chan models.IncomingMessage
	dispatcher broker.Dispatcher
}

// Init creates and initialize Broker instance
func Init(configFilePath string) (Broker, error) {
	configManager, err := config.NewManager(configFilePath)
	if err != nil {
		return Broker{}, entity.ConfigInitFailed{Message: err.Error()}
	}

	workspace := broker.NewWorkspace(configManager.Workspace())
	transmitter := services.NewCommunicationManager()
	cmdDispatcher := broker.NewCommandDispatcher(&workspace, transmitter)
	connDispatcher := broker.NewConnectionDispatcher(&workspace, cmdDispatcher)

	return Broker{
		workspace:  workspace,
		server:     broker.InitServer(configManager.TCPAddress(), configManager.TCPServerConnectionType()),
		incoming:   make(chan models.IncomingMessage),
		dispatcher: broker.NewDispatcher(&workspace, connDispatcher, cmdDispatcher, transmitter),
	}, nil
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

func (broker Broker) listenIncomingMessages(user entity.User) {
	var message models.IncomingMessage
	for {
		err := user.Connection.GetMessage(&message)
		switch err {
		case io.EOF:
			fmt.Println("Disconnected ", user.NickName, user.ID)
			return
		case nil:
			broker.incoming <- message
		default:
			fmt.Println("Error at decoding message ", user.NickName, user.ID, err)
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

func (broker Broker) register(connection common.Connection) {
	user, err := broker.dispatcher.RegisterNewConnection(connection)
	switch err.(type) {
	case nil:
		broker.listenIncomingMessages(user)
	default:
		fmt.Println("Register of new user failed ", err)
		broker.close(connection)
	}
}

func (broker Broker) close(connection common.Connection) {
	err := connection.Close()
	if err != nil {
		fmt.Println("Error during closing connection ", err)
	}
}
