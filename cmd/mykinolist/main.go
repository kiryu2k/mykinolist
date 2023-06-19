package main

import (
	"log"

	"github.com/kiryu-dev/mykinolist/internal/app"
	"github.com/kiryu-dev/mykinolist/internal/config"
)

const configPath = "./configs/config.yaml"

func main() {
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	app.Run(config)
}
