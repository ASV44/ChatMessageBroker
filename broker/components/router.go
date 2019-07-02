package broker

import (
	"github.com/gorilla/mux"
	"net/http"
)

func New() *mux.Router {
	router := mux.NewRouter()
	addWebSocketRoutes(router)

	return router
}

func addWebSocketRoutes(router *mux.Router) {
	router.Path("/connect/{id}").
		Methods(http.MethodGet, http.MethodOptions).
		HandlerFunc(acceptConnections).
		Name("connect")
}