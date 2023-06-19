package repository

import (
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/mykinolist/internal/config"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
