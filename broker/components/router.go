package broker

import (
	"github.com/ASV44/ChatMessageBroker/broker/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(websocketService controllers.WebsocketService) *mux.Router {
	router := mux.NewRouter()
	addWebSocketRoutes(router, websocketService)

	return router
}

func addWebSocketRoutes(router *mux.Router, websocketService controllers.WebsocketService) {
	router.Path("/connect/{id}").
		Methods(http.MethodGet, http.MethodOptions).
		Handler(controllers.ServeWebSocket(websocketService)).
		Name("connect")
}
