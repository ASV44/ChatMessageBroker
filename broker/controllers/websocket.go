package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebsocketService represents abstraction for managing web socket connections
type WebsocketService interface {
	HandleNewConnection(connection *websocket.Conn)
}

// ServeWebSocket handles websocket requests from the peer.
func ServeWebSocket(upgrader websocket.Upgrader, websocketService WebsocketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error at websocket connection", err)
			return
		}

		websocketService.HandleNewConnection(conn)
	}
}
