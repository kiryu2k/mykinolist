package repository

import (
	"context"
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) CreateAccount(ctx context.Context, user *model.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	query := `
INSERT INTO users (username, email, Hashed_password, created_on, last_login)
VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	err = tx.QueryRowContext(ctx, query,
		user.Username, user.Email, user.HashedPassword,
		user.CreatedOn, user.LastLogin).Scan(&user.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
SELECT * FROM users
WHERE email = $1;
	`
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	query := `
UPDATE users
SET last_login = $1
WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, query, user.LastLogin, user.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}
