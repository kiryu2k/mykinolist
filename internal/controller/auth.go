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

const cookieMaxAge = 30 * 24 * 60 * 60 // 30 days

type userIDKey struct{}

type authHandler struct {
	service service.AuthService
}

// SignUp godoc
// @Summary      Sign up an account
// @Description  Create an account and an empty movie list, which is linked to the account ID
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body model.SignUpUserDTO true "account info"
// @Success      200      {object}  model.ListInfo
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /auth/signup [post]
func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	req := new(model.SignUpUserDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	list, err := h.service.SignUp(req)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, list)
}

// SignIn godoc
// @Summary      Sign in to account
// @Description  Log in to an existing account, update access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body model.SignInUserDTO true "account info"
// @Success      200      {object}  model.Tokens
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /auth/signin [post]
func (h *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	req := new(model.SignInUserDTO)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	tokens, err := h.service.SignIn(req)
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

// SignOut godoc
// @Summary      Sign out of account
// @Description  Sign out of account by deleting refresh token from cookies
// @Tags         auth
// @Produce      json
// @Success      200      {string}	string
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /auth/signout [post]
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("u've successfully logged out"))
}

// GetUser godoc
// @Summary      Get user info
// @Security	 AccessToken
// @Description  Get user info: id, username, email, last login time, etc.
// @Tags         user
// @Produce      json
// @Param 		 id path int true "User ID"
// @Success      200      {object}  model.User
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /user/{id} [get]
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
	user, err := h.service.GetUser(id)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSONResponse(w, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete account
// @Security	 AccessToken
// @Description  Delete user account
// @Tags         user
// @Produce      json
// @Param 		 id path int true "User ID"
// @Success      200      {object}  model.User
// @Failure      400,404  {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Failure      default  {object}  errorResponse
// @Router       /user/{id} [delete]
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
	user, err := h.service.Delete(id)
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
