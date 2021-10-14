package broker

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/ASV44/chat-message-broker/common"
)

// TCPServer represents instance of running server
type TCPServer struct {
	Address        string
	ConnectionType string
	Connection     chan common.Connection
}

// InitServer creates and initialize instance of TCPServer
func InitServer(address string, connectionType string) TCPServer {
	server := TCPServer{
		Address:        address,
		ConnectionType: connectionType,
		Connection:     make(chan common.Connection),
	}

	return server
}

// Start init and start tcp server and start accepting connections
func (server TCPServer) Start() error {
	listener, err := net.Listen(server.ConnectionType, server.Address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("broker is running on :", server.Address)

	// Start goroutine for handling new incoming connections
	go server.run(listener)

	return nil
}

func (server TCPServer) run(listener net.Listener) {
	defer server.close(listener)
	server.acceptConnections(listener)
}

func (server TCPServer) acceptConnections(listener net.Listener) {
	// connection, err := server.upgrader.Upgrade(w, r, nil)
	for {
		rawConnection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept new connection ", err)
		} else {
			server.Connection <- common.NewConnection(rawConnection, common.NewJSONConnIO(rawConnection))
		}
	}
}

// IsConnectionActive checks if provided connection is still active
func (server TCPServer) IsConnectionActive(connection net.Conn) bool {
	err := connection.SetReadDeadline(time.Now())
	if err != nil {
		fmt.Println("Could not set read deadline ", err)
	}

	var isConnected bool
	var one []byte
	if _, err := connection.Read(one); errors.Is(err, io.EOF) {
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
func (server TCPServer) close(listener net.Listener) {
	if err := listener.Close(); err != nil {
		fmt.Println("Could not close server listener ", err)
	}
}
