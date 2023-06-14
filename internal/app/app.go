package app

import (
	"net/http"

	"github.com/kiryu-dev/mykinolist/internal/controller"
	"github.com/kiryu-dev/mykinolist/internal/infrastructure/repository"
	"github.com/kiryu-dev/mykinolist/internal/service"
	"github.com/sirupsen/logrus"
)

func Run(port string, config *repository.Config) {
	db, err := repository.NewPostgresDB(config)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	repo := repository.New(db)
	service := service.New(repo)
	controller := controller.New(service)
	if err := http.ListenAndServe(port, controller); err != nil {
		logrus.Fatal(err.Error())
	}
}
