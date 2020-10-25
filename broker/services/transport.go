package services

import "github.com/gorilla/websocket"

// WebsocketJSONConnIO represents websocket abstraction of JSON connection communication
type WebsocketJSONConnIO struct {
	connection *websocket.Conn
}

// NewWebsocketJSONConnIO creates new instance of WebsocketJSONConnIO
func NewWebsocketJSONConnIO(connection *websocket.Conn) WebsocketJSONConnIO {
	return WebsocketJSONConnIO{connection: connection}
}

// SendMessage send JSON message to websocket connection
func (conn WebsocketJSONConnIO) SendMessage(message interface{}) error {
	return conn.connection.WriteJSON(&message)
}

// GetMessage get JSON message from websocket connection
func (conn WebsocketJSONConnIO) GetMessage(message interface{}) error {
	return conn.connection.ReadJSON(&message)
}
