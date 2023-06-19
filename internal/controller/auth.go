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

const cookieMaxAge = 30 * 24 * 60 * 60 // 30 days

type authHandler struct {
	service service.AuthService
}

func InitAuthRoutes(router *mux.Router, s service.AuthService) {
	handler := &authHandler{service: s}
	router.HandleFunc("/auth/signup", handler.signUp).Methods(http.MethodPost)
	router.HandleFunc("/auth/signin", handler.signIn).Methods(http.MethodGet)
}

func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	req := new(model.SignUpUserDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	defer r.Body.Close()
	startTime := time.Now()
	user, tokens, err := h.service.SignUp(req)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	w.Header().Add("X-Authorization", tokens.AccessToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
	})
	writeJSONResponse(w, http.StatusOK, map[string]any{"user": user, "tokens": tokens})
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	req := new(model.SignInUserDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	defer r.Body.Close()
	startTime := time.Now()
	user, tokens, err := h.service.SignIn(req)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	w.Header().Add("X-Authorization", tokens.AccessToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
	})
	writeJSONResponse(w, http.StatusOK, map[string]any{"user": user, "tokens": tokens})
}
