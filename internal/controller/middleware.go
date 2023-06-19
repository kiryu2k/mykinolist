package controller

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

func (h *authHandler) IdentifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("X-Authorization")
		tokenParts := strings.Split(tokenStr, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			resp := &errorResponse{"invalid x-authorization header"}
			writeJSONResponse(w, http.StatusUnauthorized, resp)
			return
		}
		userID, err := model.ParseAccessToken(tokenParts[1], os.Getenv("JWT_ACCESS_SECRET_KEY")) // TODO getter of secret key
		if err != nil {
			resp := &errorResponse{err.Error()}
			writeJSONResponse(w, http.StatusUnauthorized, resp)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
