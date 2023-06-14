package repository

import (
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/service"
)

type Repository struct {
	service.AuthRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{&userRepository{db}}
}
