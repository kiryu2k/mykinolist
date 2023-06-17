package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/model"
	"github.com/kiryu-dev/mykinolist/internal/service"
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
	startTime := time.Now()
	req := new(model.SignUpUserDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	defer r.Body.Close()
	user, err := h.service.SignUp(req)
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	writeJSONResponse(w, http.StatusOK, user)
	log.Printf("elapsed time: %v", time.Since(startTime))
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	req := new(model.SignInUserDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	defer r.Body.Close()
	user, err := h.service.SignIn(req)
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	writeJSONResponse(w, http.StatusOK, user)
}
