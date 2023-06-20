package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/model"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

const cookieMaxAge = 30 * 24 * 60 * 60 // 30 days

type authHandler struct {
	service service.AuthService
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
	user, err := h.service.SignUp(req)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	writeJSONResponse(w, http.StatusOK, user)
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
	w.Header().Add("X-Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
	})
	writeJSONResponse(w, http.StatusOK, map[string]any{"user": user, "tokens": tokens})
}

func (h *authHandler) getUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	idFromCtx := r.Context().Value("userID")
	if id != idFromCtx {
		resp := &errorResponse{"cannot get other's account info"}
		writeJSONResponse(w, http.StatusForbidden, resp)
		return
	}
	startTime := time.Now()
	user, err := h.service.GetUser(id)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	writeJSONResponse(w, http.StatusOK, user)
}
