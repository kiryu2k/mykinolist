package controller

import (
	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

func New(auth service.AuthService) *mux.Router {
	router := mux.NewRouter()
	InitAuthRoutes(router, auth)
	// InitListRoutes(router, list)
	return router
}
