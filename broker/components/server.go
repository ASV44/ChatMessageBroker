package broker

import (
	"fmt"
	"io"
	"net"
	"os"
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
	Connections    chan net.Conn
	listener       net.Listener
}

// Start init and start tcp server and start accepting connections
func (server *Server) Start() {
	server.Connections = make(chan net.Conn)
	var err error
	server.listener, err = net.Listen(server.ConnectionType, server.Host+":"+server.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("broker is running on port :", server.Port)
	}

	go server.acceptConnections()
}

func (server *Server) acceptConnections() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		server.Connections <- connection
	}
}

// IsConnectionActive checks if provided connection is still active
func (server *Server) IsConnectionActive(connection net.Conn) bool {
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

// Close end server listening of new connection
func (server *Server) Close() {
	err := server.listener.Close()
	if err != nil {
		fmt.Println("Could not close server listener ", err)
	}
}
