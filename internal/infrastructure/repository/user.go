package repository

import (
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type userRepository struct {
	db *sql.DB
}

// TODO: querycontext to cancel query
func (r *userRepository) CreateAccount(user *model.User) error {
	query := `
INSERT INTO users (username, email, encrypted_password, created_on, last_login)
VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	rows, err := r.db.Query(query, user.Username, user.Email, user.Password,
		user.CreatedOn, user.LastLogin)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&user.ID); err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) LoginAccount(user *model.User) error {
	return nil
}
