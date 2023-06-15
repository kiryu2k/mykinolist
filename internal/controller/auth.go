package controller

import (
	"encoding/json"
	"net/http"

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

type signUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	req := new(signUpRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
	}
	defer r.Body.Close()
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.service.SignUp(user); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	writeJSONResponse(w, http.StatusOK, user)
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {

}
