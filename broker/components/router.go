package broker

import (
	"github.com/ASV44/ChatMessageBroker/broker/controllers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

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
	router.Path("/connect/{id}").
		Methods(http.MethodGet, http.MethodOptions).
		Handler(controllers.ServeWebSocket(upgrader, websocketService)).
		Name("connect")
}
