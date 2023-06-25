package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/model"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

type listHandler struct {
	service service.ListService
}

func (h *listHandler) addMovie(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIDKey{}).(int64)
	req := new(model.ListUnit)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	req.OwnerID = id
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.service.AddMovie(ctx, req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Write([]byte(fmt.Sprintf("Movie %s has successfully added to user's %d list", req.Name, id)))
}

func (h *listHandler) getMovies(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIDKey{}).(int64)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	movies, err := h.service.GetMovies(ctx, id)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, movies)
}

func (h *listHandler) updateMovie(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey{}).(int64)
	idStr := mux.Vars(r)["id"]
	movieID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	req := new(model.ListUnitPatch)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req.OwnerID = &userID
	req.MovieID = &movieID
	if err := h.service.UpdateMovie(ctx, req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Write([]byte("movie data has been updated"))
}
