package services

import (
	"fmt"
	"github.com/ASV44/ChatMessageBroker/common"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type WebsocketProcessor struct {
	WebSocketConn chan common.Connection
}

func NewWebsocketProcessor() WebsocketProcessor {
	return WebsocketProcessor{
		WebSocketConn: make(chan common.Connection),
	}
}

func (service WebsocketProcessor) HandleNewConnection(websocketConn *websocket.Conn) {
	connection := common.NewConnection(websocketConn, NewWebsocketJSONConnIO(websocketConn))

	websocketConn.SetReadLimit(maxMessageSize)
	websocketConn.SetPongHandler(service.pongHandler(websocketConn))
	err := websocketConn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Println("Websocket handle new connection error at setting read deadline", err)
	}

	go service.ping(websocketConn)

	service.WebSocketConn <- connection
}

func (service WebsocketProcessor) pongHandler(conn *websocket.Conn) func(string) error {
	return func(appData string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	}
}

func (service WebsocketProcessor) ping(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer service.dispose(conn, ticker)
	for {
		select {
		case <-ticker.C:
			err := conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				fmt.Println("Websocket ping error at setting write deadline", err)
			}

			if err = conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("Websocket ping error at writing ping message", err)
			}
		}
	}
}

func (service WebsocketProcessor) dispose(conn *websocket.Conn, ticker *time.Ticker) {
	ticker.Stop()
	err := conn.Close()
	if err != nil {
		fmt.Println("Websocket ping error at closing connection", err)
	}
}
