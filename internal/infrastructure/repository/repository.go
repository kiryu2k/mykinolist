package repository

import (
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/service"
)

type Repository struct {
	service.UserRepository
	service.TokenRepository
	service.ListRepository
	service.MovieRepositroy
}

func New(db *sql.DB) *Repository {
	return &Repository{
		&userRepository{db},
		&tokenRepository{db},
		&listRepository{db},
		&movieRepository{db},
	}
}
