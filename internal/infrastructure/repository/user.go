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

func (r *userRepository) FindUserByName(username string) (*model.User, error) {
	query := `
SELECT * FROM users
WHERE username = $1;
	`
	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return scanUser(rows)
	}
	return nil, err
}

func (r *userRepository) UpdateLastLogin(user *model.User) error {
	query := `
UPDATE users
SET last_login = $1
WHERE id = $2
	`
	_, err := r.db.Query(query, user.LastLogin, user.ID)
	return err
}

func scanUser(rows *sql.Rows) (*model.User, error) {
	u := new(model.User)
	err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedOn, &u.LastLogin)
	return u, err
}
