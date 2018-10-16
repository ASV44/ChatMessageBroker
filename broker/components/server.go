package broker

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	DEFAULT_HOST = "localhost"
	DEFAULT_PORT = "8888"
	DEFAULT_TYPE = "tcp"
)

type Server struct {
	Host           string
	Port           string
	ConnectionType string
	Listener       net.Listener
	Connections    chan net.Conn
}

func (server *Server) Start() {
	server.Connections = make(chan net.Conn)
	var err error
	server.Listener, err = net.Listen(server.ConnectionType, server.Host+":"+server.Port)
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
		connection, err := server.Listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		server.Connections <- connection
	}
}

func (server *Server) IsConnectionActive(connection *net.Conn) bool {
	(*connection).SetReadDeadline(time.Now())
	var isConnected bool
	var one []byte
	if _, err := (*connection).Read(one); err == io.EOF {
		isConnected = false
	} else {
		var zero time.Time
		(*connection).SetReadDeadline(zero)
		isConnected = true
	}

	return isConnected
}
