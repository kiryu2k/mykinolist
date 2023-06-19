package repository

import (
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/service"
)

type Repository struct {
	service.UserRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{&userRepository{db}}
}
