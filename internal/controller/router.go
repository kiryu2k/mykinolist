package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

func New(auth service.AuthService) *mux.Router {
	router := mux.NewRouter()
	initAuthRoutes(router.PathPrefix("/auth").Subrouter(), auth)
	// InitListRoutes(router.PathPrefix("/list").Subrouter(), list)
	return router
}

func initAuthRoutes(router *mux.Router, s service.AuthService) {
	handler := &authHandler{service: s}
	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/signup", handler.signUp)
	postRouter.HandleFunc("/signin", handler.signIn)
	postRouter.HandleFunc("/signout", handler.signOut)
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/{id:[0-9]+}", handler.getUser).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id:[0-9]+}", handler.deleteUser).Methods(http.MethodDelete)
	userRouter.Use(handler.IdentifyUser)
}
