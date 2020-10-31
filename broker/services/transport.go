package services

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

// WebsocketJSONConnIO represents websocket abstraction of JSON connection communication
type WebsocketJSONConnIO struct {
	connection *websocket.Conn
	writeWait  time.Duration
}

// NewWebsocketJSONConnIO creates new instance of WebsocketJSONConnIO
func NewWebsocketJSONConnIO(connection *websocket.Conn, writeWait time.Duration) WebsocketJSONConnIO {
	return WebsocketJSONConnIO{connection: connection, writeWait: writeWait}
}

// SendMessage send JSON message to websocket connection
func (conn WebsocketJSONConnIO) SendMessage(message interface{}) error {
	err := conn.connection.SetWriteDeadline(time.Now().Add(conn.writeWait))
	if err != nil {
		fmt.Println("Websocket message send error at setting write deadline", err)
		return err
	}

	return conn.connection.WriteJSON(&message)
}

// GetMessage get JSON message from websocket connection
func (conn WebsocketJSONConnIO) GetMessage(message interface{}) error {
	return conn.connection.ReadJSON(&message)
}
