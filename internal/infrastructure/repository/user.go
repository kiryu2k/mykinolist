package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) CreateAccount(ctx context.Context, user *model.User) error {
	query := `
INSERT INTO users (username, email, hashed_password, created_on, last_login)
VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.HashedPassword,
		user.CreatedOn, user.LastLogin).Scan(&user.ID)
	return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = $1;`
	user := new(model.User)
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.HashedPassword, &user.CreatedOn, &user.LastLogin,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, user *model.User) error {
	query := `
UPDATE users
SET last_login = $1
WHERE id = $2;
`
	_, err := r.db.ExecContext(ctx, query, user.LastLogin, user.ID)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT * FROM users WHERE id = $1;`
	user := new(model.User)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.HashedPassword, &user.CreatedOn, &user.LastLogin,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteAccount(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1;`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("invalid count of deleted users %d", count)
	}
	return nil
}
