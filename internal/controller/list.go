package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/model"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

type listHandler struct {
	service service.ListService
}

// AddMovie godoc
// @Summary      Add movie to list
// @Security	 AccessToken
// @Description  Use a third-party API to search for movie information by title and, if successful, add the movie to the list. You can add the movie to your favorites, rate it, and specify movie status (watching, plan to watch, etc.)
// @Tags         list
// @Accept       json
// @Produce      json
// @Param 		 input body model.ListUnit true "movie info"
// @Success      200      {string}  string
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /list [post]
func (h *listHandler) addMovie(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIDKey{}).(int64)
	req := new(model.ListUnit)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	r.Body.Close()
	req.OwnerID = id
	if err := h.service.AddMovie(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Movie %s has successfully added to user's %d list", req.Name, id)))
}

// GetMovies godoc
// @Summary      Get movies
// @Security	 AccessToken
// @Description  Get all movies from list
// @Tags         list
// @Produce      json
// @Success      200      {array}   model.ListUnit
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /list [get]
func (h *listHandler) getMovies(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(userIDKey{}).(int64)
	movies, err := h.service.GetMovies(id)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, movies)
}

// UpdateMovie godoc
// @Summary      Update movie info
// @Security	 AccessToken
// @Description  Update part of the information about the added movie. For example, you can add the movie to your favorites or change the rating of the movie.
// @Tags         list
// @Produce      json
// @Param 		 id path int true "Movie ID"
// @Success      200      {string}  string
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /list/{id} [patch]
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
	defer r.Body.Close()
	req.OwnerID = &userID
	req.MovieID = &movieID
	if err := h.service.UpdateMovie(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("movie data has been updated"))
}

// DeleteMovie godoc
// @Summary      Delete movie
// @Security	 AccessToken
// @Description  Delete movie from list
// @Tags         list
// @Produce      json
// @Param 		 id path int true "Movie ID"
// @Success      200      {object}  model.ListUnit
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /list/{id} [delete]
func (h *listHandler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey{}).(int64)
	idStr := mux.Vars(r)["id"]
	movieID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	req := &model.ListUnit{
		Movie:    model.Movie{ID: movieID},
		ListInfo: model.ListInfo{OwnerID: userID},
	}
	if err := h.service.DeleteMovie(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, req)
}
