package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

func New(auth service.AuthService, list service.ListService) *mux.Router {
	var (
		router      = mux.NewRouter()
		authHandler = &authHandler{service: auth}
		listHandler = &listHandler{service: list}
		middleware  = &authMiddleware{service: auth}
		authRouter  = router.PathPrefix("/auth").Subrouter()
		userRouter  = router.PathPrefix("/user").Subrouter()
		listRouter  = router.PathPrefix("/list").Subrouter()
	)
	{
		authRouter.HandleFunc("/signup", authHandler.signUp).Methods(http.MethodPost)
		authRouter.HandleFunc("/signin", authHandler.signIn).Methods(http.MethodPost)
		authRouter.HandleFunc("/signout", authHandler.signOut).Methods(http.MethodPost)
	}
	{
		userRouter.Use(middleware.identifyUser)
		userRouter.HandleFunc("/{id:[0-9]+}", authHandler.getUser).Methods(http.MethodGet)
		userRouter.HandleFunc("/{id:[0-9]+}", authHandler.deleteUser).Methods(http.MethodDelete)
	}
	{
		listRouter.Use(middleware.identifyUser)
		listRouter.HandleFunc("", listHandler.addMovie).Methods(http.MethodPost)
		listRouter.HandleFunc("", listHandler.getMovies).Methods(http.MethodGet)
		listRouter.HandleFunc("/{id:[0-9]+}", listHandler.updateMovie).Methods(http.MethodPatch)
	}
	return router
}
