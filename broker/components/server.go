package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/broker/entity"
	"github.com/ASV44/ChatMessageBroker/broker/services"
	"io"
	"net"
	"time"
)

// Constant value of server config
const (
	DefaultHost = "localhost"
	DefaultPort = "8888"
	DefaultType = "tcp"
)

// Server represents instance of running server
type Server struct {
	Host           string
	Port           string
	ConnectionType string
	Connection     chan entity.Connection
}

// InitServer creates and initialize instance of Server
func InitServer(host string, port string, connectionType string) Server {
	server := Server{
		Host:           host,
		Port:           port,
		ConnectionType: connectionType,
		Connection:     make(chan entity.Connection),
	}

	return server
}

// Start init and start tcp server and start accepting connections
func (server Server) Start() error {
	listener, err := net.Listen(server.ConnectionType, server.Host+":"+server.Port)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("broker is running on port :", server.Port)

	// Start goroutine for handling new incoming connections
	go server.run(listener)

	return nil
}

func (server Server) run(listener net.Listener) {
	defer server.close(listener)
	server.acceptConnections(listener)
}

func (server Server) acceptConnections(listener net.Listener) {
	for {
		rawConnection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept new connection ", err)
		} else {
			server.Connection <- entity.NewConnection(rawConnection, services.NewJsonConnIO(rawConnection))
		}
	}
}

// IsConnectionActive checks if provided connection is still active
func (server Server) IsConnectionActive(connection net.Conn) bool {
	err := connection.SetReadDeadline(time.Now())
	if err != nil {
		fmt.Println("Could not set read deadline ", err)
	}

	var isConnected bool
	var one []byte
	if _, err := connection.Read(one); err == io.EOF {
		isConnected = false
	} else {
		var zero time.Time
		err = connection.SetReadDeadline(zero)
		if err != nil {
			fmt.Println("Could not set read deadline to zero value ", err)
		}

		isConnected = true
	}

	return isConnected
}

// close end server listening of new connection
func (server Server) close(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		fmt.Println("Could not close server listener ", err)
	}
}
