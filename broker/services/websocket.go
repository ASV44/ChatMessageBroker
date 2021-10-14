package services

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/ASV44/chat-message-broker/common"

	"github.com/ASV44/chat-message-broker/broker/config"
)

// WebsocketProcessor represents implementation of websocket connection handling logic
type WebsocketProcessor struct {
	websocketSettings config.WebsocketConnectionSettings
	WebSocketConn     chan common.Connection
}

// NewWebsocketProcessor creates new instance of WebsocketProcessor
func NewWebsocketProcessor(websocketSettings config.WebsocketConnectionSettings) WebsocketProcessor {
	return WebsocketProcessor{
		websocketSettings: websocketSettings,
		WebSocketConn:     make(chan common.Connection),
	}
}

// HandleNewConnection process new websocket connection by wrapping it to broker connection abstraction
func (service WebsocketProcessor) HandleNewConnection(websocketConn *websocket.Conn) {
	connIO := NewWebsocketJSONConnIO(websocketConn, service.websocketSettings.WriteWait)
	connection := common.NewConnection(websocketConn, connIO)

	websocketConn.SetReadLimit(service.websocketSettings.MaxMessageSize)
	websocketConn.SetPongHandler(service.pongHandler(websocketConn))
	err := websocketConn.SetReadDeadline(time.Now().Add(service.websocketSettings.PongWait))
	if err != nil {
		fmt.Println("Websocket handle new connection error at setting read deadline", err)
	}

	go service.ping(websocketConn)

	service.WebSocketConn <- connection
}

func (service WebsocketProcessor) pongHandler(conn *websocket.Conn) func(string) error {
	return func(appData string) error {
		return conn.SetReadDeadline(time.Now().Add(service.websocketSettings.PongWait))
	}
}

func (service WebsocketProcessor) ping(conn *websocket.Conn) {
	ticker := time.NewTicker(service.websocketSettings.PingPeriod)
	defer service.dispose(conn, ticker)
	for range ticker.C {
		err := conn.SetWriteDeadline(time.Now().Add(service.websocketSettings.WriteWait))
		if err != nil {
			fmt.Println("Websocket ping error at setting write deadline", err)
		}

		if err = conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			fmt.Println("Websocket ping error at writing ping message", err)
		}
	}
}

func (service WebsocketProcessor) dispose(conn *websocket.Conn, ticker *time.Ticker) {
	ticker.Stop()
	if err := conn.Close(); err != nil {
		fmt.Println("Websocket ping error at closing connection", err)
	}
}
