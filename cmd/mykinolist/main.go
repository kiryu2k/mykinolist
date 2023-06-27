package main

import (
	"log"

	"github.com/kiryu-dev/mykinolist/internal/app"
	"github.com/kiryu-dev/mykinolist/internal/config"
)

// @title MyKinoList API
// @version 1.0
// @description MyKinoList is an API that provides the ability to create accounts to keep a list of movies you are watching, have watched, or may want to watch sometime later, as well as to rate them and add them to your «favorites».\n\nThis API uses JSONWebTokens to give access to authorized users. A third-party, unofficial Kinopoisk API is used to retrieve information about movies.\n\nThe documentation uses an access token for authorization, which only lasts 30 seconds.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey AccessToken
// @in header
// @name Authorization

const configPath = "./configs/config.yaml"

func main() {
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	app.Run(config)
}
