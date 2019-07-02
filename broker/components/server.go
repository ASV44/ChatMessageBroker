package broker

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/common"
	"io"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// Server represents instance of running server
type Server struct {
	Address        string
	ConnectionType string
	Connection     chan common.Connection
	WebSocketConn  chan *websocket.Conn
	upgrader 		   websocket.Upgrader
}

// InitServer creates and initialize instance of Server
func InitServer(address string, connectionType string) Server {
	server := Server{
		Address:        address,
		ConnectionType: connectionType,
		Connection:     make(chan common.Connection),
	}

	return server
}

// Start init and start tcp server and start accepting connections
func (server Server) Start() error {
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

func (server Server) run(listener net.Listener) {
	defer server.close(listener)
	server.acceptConnections(listener)
}

func (server Server) acceptConnections(listener net.Listener) {
	// connection, err := server.upgrader.Upgrade(w, r, nil)
	for {
		rawConnection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept new connection ", err)
		} else {
			server.Connection <- common.NewConnection(rawConnection, common.NewJSONConnIO(rawConnection))
		}
	}
	server.Connections <- connection
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

// serveWebSocket handles websocket requests from the peer.
func (server *erver) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	//client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// close end server listening of new connection
func (server Server) close(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		fmt.Println("Could not close server listener ", err)
	}
}
