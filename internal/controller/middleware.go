package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (h *authHandler) identifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		tokenParts := strings.Split(tokenStr, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			resp := &errorResponse{"invalid authorization header"}
			writeJSONResponse(w, http.StatusUnauthorized, resp)
			return
		}
		userID, err := h.service.ParseAccessToken(tokenParts[1])
		ctx := context.WithValue(r.Context(), userIDKey{}, userID)
		r = r.WithContext(ctx)
		if err == nil {
			next.ServeHTTP(w, r)
			return
		}
		if err.Error() != "token expiration date has passed" {
			resp := &errorResponse{err.Error()}
			writeJSONResponse(w, http.StatusUnauthorized, resp)
			return
		}
		refreshToken, err := r.Cookie("refreshToken")
		if err != nil {
			resp := &errorResponse{err.Error()}
			writeJSONResponse(w, http.StatusBadRequest, resp)
			return
		}
		if id, err := h.service.ParseRefreshToken(refreshToken.Value); err != nil || id != userID {
			resp := &errorResponse{"token expiration date has passed"}
			writeJSONResponse(w, http.StatusBadRequest, resp)
			return
		}
		h.updateTokens(w, r)
		next.ServeHTTP(w, r)
	})
}

func (h *authHandler) updateTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey{}).(int64)
	tokens, err := h.service.UpdateTokens(userID)
	if err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusUnauthorized, resp)
		return
	}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	removeRefreshTokenCookie(w) // remove old cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Path:     "/auth",
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
	})
}
