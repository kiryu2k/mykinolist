package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/kiryu-dev/mykinolist/internal/service"
)

type authMiddleware struct {
	service service.AuthService
}

func (m *authMiddleware) identifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		tokenParts := strings.Split(tokenStr, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			writeErrorJSON(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}
		userID, err := m.service.ParseAccessToken(tokenParts[1])
		ctx := context.WithValue(r.Context(), userIDKey{}, userID)
		r = r.WithContext(ctx)
		if err == nil {
			next.ServeHTTP(w, r)
			return
		}
		if err.Error() != "token expiration date has passed" {
			writeErrorJSON(w, http.StatusUnauthorized, err.Error())
			return
		}
		refreshToken, err := r.Cookie("refreshToken")
		if err != nil {
			writeErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		if id, err := m.service.ParseRefreshToken(refreshToken.Value); err != nil || id != userID {
			writeErrorJSON(w, http.StatusBadRequest, "token expiration date has passed")
			return
		}
		m.updateTokens(w, r)
		next.ServeHTTP(w, r)
	})
}

func (m *authMiddleware) updateTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey{}).(int64)
	tokens, err := m.service.UpdateTokens(userID)
	if err != nil {
		writeErrorJSON(w, http.StatusUnauthorized, err.Error())
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
