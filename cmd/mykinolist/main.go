package main

import (
	"log"

	"github.com/kiryu-dev/mykinolist/internal/app"
	"github.com/spf13/viper"
)

const configPath = "./configs/config.yaml"

func main() {
	if err := loadConfigFile(); err != nil {
		log.Fatal(err.Error())
	}
	app.Run(viper.GetString("port"))
}

func loadConfigFile() error {
	viper.SetConfigFile(configPath)
	return viper.ReadInConfig()
}
