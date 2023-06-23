package app

import (
	"log"
	"net/http"

	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/controller"
	"github.com/kiryu-dev/mykinolist/internal/infrastructure/repository"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

func Run(config *config.Config) {
	db, err := repository.NewPostgresDB(config.DB)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	repo := repository.New(db)
	services := service.New(repo.UserRepository, repo.TokenRepository, config)
	controller := controller.New(services.AuthService, services.ListService)
	if err := http.ListenAndServe(config.ListeningPort, controller); err != nil {
		log.Fatal(err.Error())
	}
}
