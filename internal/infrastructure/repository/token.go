package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type tokenRepository struct {
	db *sql.DB
}

func (r *tokenRepository) Save(ctx context.Context, userToken *model.UserToken) error {
	_, err := r.findByID(ctx, userToken.UserID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		/* there's no user in token table, so we need to add him */
		err = r.add(ctx, userToken)
	} else {
		err = r.update(ctx, userToken)
	}
	return err
}

func (r *tokenRepository) Remove(ctx context.Context, refreshToken string) error {
	query := `DELETE FROM tokens WHERE refresh_token = $1;`
	res, err := r.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("invalid count of deleted tokens %d", count)
	}
	return nil
}

func (r *tokenRepository) findByID(ctx context.Context, id int64) (*model.UserToken, error) {
	userToken := new(model.UserToken)
	query := `SELECT * FROM tokens WHERE user_id = $1;`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&userToken.UserID, &userToken.RefreshToken)
	if err != nil {
		return nil, err
	}
	return userToken, nil
}

func (r *tokenRepository) update(ctx context.Context, userToken *model.UserToken) error {
	query := `
UPDATE tokens
SET refresh_token = $1
WHERE user_id = $2;
	`
	res, err := r.db.ExecContext(ctx, query, userToken.RefreshToken, userToken.UserID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("invalid count of updated tokens %d", count)
	}
	return nil
}

func (r *tokenRepository) add(ctx context.Context, userToken *model.UserToken) error {
	query := `
INSERT INTO tokens (user_id, refresh_token)
VALUES ($1, $2);
	`
	_, err := r.db.ExecContext(ctx, query, userToken.UserID, userToken.RefreshToken)
	return err
}
