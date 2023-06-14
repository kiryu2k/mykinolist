package repository

import (
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) CreateAccount(user *model.User) error {
	return nil
}

func (r *userRepository) LoginAccount(user *model.User) error {
	return nil
}
