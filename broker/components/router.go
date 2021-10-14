package broker

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/ASV44/chat-message-broker/broker/controllers"
)

// NewRouter creates new instance of router for HTTP server
func NewRouter(upgrader websocket.Upgrader, websocketService controllers.WebsocketService) *mux.Router {
	router := mux.NewRouter()
	addWebSocketRoutes(router, upgrader, websocketService)

	return router
}

func addWebSocketRoutes(
	router *mux.Router,
	upgrader websocket.Upgrader,
	websocketService controllers.WebsocketService,
) {
	router.Path("/websocket-connect").
		Methods(http.MethodGet, http.MethodOptions).
		Handler(controllers.ServeWebSocket(upgrader, websocketService)).
		Name("connect")
}
