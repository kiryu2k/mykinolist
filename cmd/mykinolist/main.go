package main

import (
	"os"

	"github.com/kiryu-dev/mykinolist/internal/app"
	"github.com/kiryu-dev/mykinolist/internal/infrastructure/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

const configPath = "./configs/config.yaml"

func main() {
	if err := loadConfigFile(); err != nil {
		logrus.Fatal(err.Error())
	}
	if err := gotenv.Load(); err != nil {
		logrus.Fatal(err.Error())
	}
	app.Run(viper.GetString("port"), &repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
}

func loadConfigFile() error {
	viper.SetConfigFile(configPath)
	return viper.ReadInConfig()
}
