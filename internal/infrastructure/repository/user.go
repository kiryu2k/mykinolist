package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type userRepository struct {
	db *sql.DB
}

// TODO: querycontext to cancel query
func (r *userRepository) CreateAccount(ctx context.Context, user *model.User) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
INSERT INTO users (username, email, encrypted_password, created_on, last_login)
VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	err := r.db.QueryRowContext(queryCtx, query,
		user.Username, user.Email, user.EncryptedPassword,
		user.CreatedOn, user.LastLogin).Scan(&user.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("failed to create a user")
	}
	return err
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
SELECT * FROM users
WHERE email = $1;
	`
	rows, err := r.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return scanUser(rows)
	}
	return nil, err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, user *model.User) error {
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
	err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.EncryptedPassword, &u.CreatedOn, &u.LastLogin)
	return u, err
}
