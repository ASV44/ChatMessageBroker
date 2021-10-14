package broker

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/gorilla/websocket"

	"github.com/ASV44/chat-message-broker/common"

	"github.com/ASV44/chat-message-broker/broker/components"
	"github.com/ASV44/chat-message-broker/broker/config"
	"github.com/ASV44/chat-message-broker/broker/entity"
	"github.com/ASV44/chat-message-broker/broker/models"
	"github.com/ASV44/chat-message-broker/broker/services"
)

// Broker socket connection supported types
const (
	tcpSocket = "tcp-socket"
	webSocket = "websocket"
	all       = "all"
)

// Broker represents main structure which contains all related fields for routing message
type Broker struct {
	workspace  broker.Workspace
	tcpServer  broker.TCPServer
	httpServer broker.HTTPServer
	incoming   chan models.IncomingMessage
	dispatcher broker.Dispatcher
	websocket  services.WebsocketProcessor
}

// Init creates and initialize Broker instance
func Init(configFilePath string) (Broker, error) {
	configManager, err := config.NewManager(configFilePath)
	if err != nil {
		return Broker{}, entity.ConfigInitFailed{Message: err.Error()}
	}

	websocketConfig, err := config.NewWebsocketConnectionSettings(configManager)
	if err != nil {
		return Broker{}, entity.WebsocketConfigDecodingFailed{Message: err.Error()}
	}

	workspace := broker.NewWorkspace(configManager.Workspace())
	transmitter := services.NewCommunicationManager()
	cmdDispatcher := broker.NewCommandDispatcher(&workspace, transmitter)
	connDispatcher := broker.NewConnectionDispatcher(&workspace, cmdDispatcher)

	websocketService := services.NewWebsocketProcessor(websocketConfig)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  websocketConfig.ReadBufferSize,
		WriteBufferSize: websocketConfig.WriteBufferSize,
	}

	return Broker{
		workspace:  workspace,
		tcpServer:  broker.InitServer(configManager.TCPAddress(), configManager.TCPServerConnectionType()),
		httpServer: broker.InitHTTPServer(configManager, broker.NewRouter(upgrader, websocketService)),
		incoming:   make(chan models.IncomingMessage),
		dispatcher: broker.NewDispatcher(&workspace, connDispatcher, cmdDispatcher, transmitter),
		websocket:  websocketService,
	}, nil
}

// Start init broker tcpServer, creates channels and start receiving and routing of connections
func (broker Broker) Start(connection string) error {
	err := broker.startByConnectionType(connection)
	if err != nil {
		return err
	}

	broker.run()

	return nil
}

func (broker Broker) startByConnectionType(connection string) error {
	switch connection {
	case tcpSocket:
		return broker.tcpServer.Start()
	case webSocket:
		return broker.httpServer.Start()
	case all:
		err := broker.tcpServer.Start()
		if err != nil {
			return err
		}
		err = broker.httpServer.Start()
		if err != nil {
			return err
		}
	default:
		return entity.NotSupportedConnectionType{ConnectionType: connection}
	}

	return nil
}

func (broker Broker) listenIncomingMessages(user entity.User) {
	var message models.IncomingMessage
	for {
		err := user.Connection.GetMessage(&message)
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			fmt.Println("Disconnected ", user.NickName, user.ID, err)
			return
		}

		broker.incoming <- message

		switch err.(type) {
		case net.Error:
			fmt.Println("Lost connection with", user.NickName, user.ID, err)
			return
		case *websocket.CloseError:
			fmt.Println("Closed connection", user.NickName, user.ID, err)
			return
		default:
			fmt.Println("Error at decoding message ", user.NickName, user.ID, err)
		}
	}
}

func (broker Broker) run() {
	for {
		select {
		case connection := <-broker.tcpServer.Connection:
			go broker.register(connection)
		case websocketConnection := <-broker.websocket.WebSocketConn:
			go broker.register(websocketConnection)
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
	if err := connection.Close(); err != nil {
		fmt.Println("Error during closing connection ", err)
	}
}
