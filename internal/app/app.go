package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kiryu-dev/mykinolist/internal/config"
	"github.com/kiryu-dev/mykinolist/internal/controller"
	"github.com/kiryu-dev/mykinolist/internal/infrastructure/repository"
	"github.com/kiryu-dev/mykinolist/internal/infrastructure/webapi"
	"github.com/kiryu-dev/mykinolist/internal/service"
)

func Run(config *config.Config) {
	db, err := repository.NewPostgresDB(config.DB)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()
	var (
		repo     = repository.New(db)
		services = service.New(
			repo.UserRepository,
			repo.TokenRepository,
			repo.ListRepository,
			repo.MovieRepositroy,
			webapi.New(config.KinopoiskAPIKey),
			config,
		)
		controller = controller.New(services.AuthService, services.ListService)
	)
	server := http.Server{
		Addr:    config.ListeningPort,
		Handler: controller,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()
	log.Println("mykinolist 😼")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("see you again, dude 😼")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
}
