package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kiryu-dev/mykinolist/internal/model"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

type listHandler struct {
	service service.ListService
}

func (h *listHandler) addMovie(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIDKey{}).(int64)
	req := new(model.Movie)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.service.AddMovie(ctx, req); err != nil {
		resp := &errorResponse{err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, resp)
		return
	}
	w.Write([]byte(fmt.Sprintf("Movie %s has successfully added to user's %d list", req.Title, id)))
}
