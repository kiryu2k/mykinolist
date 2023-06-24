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

type userIDKey struct{}

type authHandler struct {
	service service.AuthService
}

func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	req := new(model.SignUpUserDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	startTime := time.Now()
	list, err := h.service.SignUp(req)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, list)
}

func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	req := new(model.SignInUserDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	startTime := time.Now()
	tokens, err := h.service.SignIn(req)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Path:     "/auth",
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
	})
	writeJSONResponse(w, http.StatusOK, tokens)
}

func (h *authHandler) signOut(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.SignOut(refreshToken.Value); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	removeRefreshTokenCookie(w)
	w.Write([]byte("u've successfully logged out"))
}

func (h *authHandler) getUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	idFromCtx := r.Context().Value(userIDKey{}).(int64)
	if id != idFromCtx {
		writeErrorJSON(w, http.StatusForbidden, "cannot get other's account info")
		return
	}
	startTime := time.Now()
	user, err := h.service.GetUser(id)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, user)
}

func (h *authHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	idFromCtx := r.Context().Value(userIDKey{}).(int64)
	if id != idFromCtx {
		writeErrorJSON(w, http.StatusForbidden, "cannot delete someone else's account")
		return
	}
	startTime := time.Now()
	user, err := h.service.Delete(id)
	log.Printf("elapsed time: %v", time.Since(startTime))
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, user)
}

func removeRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/auth",
		MaxAge:   -1,
		HttpOnly: true,
	})
}
