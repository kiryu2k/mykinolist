package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/model"
	"github.com/kiryu-dev/mykinolist/internal/service"
	"github.com/sirupsen/logrus"
)

type authHandler struct {
	service service.AuthService
}

func InitAuthRoutes(router *mux.Router, s service.AuthService) {
	handler := &authHandler{service: s}
	router.HandleFunc("/auth/signup", handler.signUp).Methods(http.MethodPost)
	router.HandleFunc("/auth/signin", handler.signIn).Methods(http.MethodGet)
}

func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	user := new(model.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logrus.Fatal(err.Error())
	}
	defer r.Body.Close()
	if err := h.service.SignUp(user); err != nil {
		logrus.Fatal(err.Error())
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {

}
