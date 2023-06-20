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
	id, err := h.service.SignUp(req)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusInternalServerError, resp)
		return
	}
	writeJSONResponse(w, http.StatusOK, map[string]int64{"id": id})
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
	tokens, err := h.service.SignIn(req)
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
	writeJSONResponse(w, http.StatusOK, tokens)
}

func (h *authHandler) signOut(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	if err := h.service.SignOut(refreshToken.Value); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	/* set to remove */
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	w.Write([]byte("u've successfully logged out"))
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
