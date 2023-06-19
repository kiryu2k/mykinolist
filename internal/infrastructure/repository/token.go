package repository

import (
	"context"
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type tokenRepository struct {
	db *sql.DB
}

func (r *tokenRepository) Save(ctx context.Context, userToken *model.UserToken) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = findTokenByID(ctx, tx, userToken.UserID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		/* there's no user in token table, so we need to add him */
		err = addUserToken(ctx, tx, userToken)
	} else {
		err = updateToken(ctx, tx, userToken)
	}
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func findTokenByID(ctx context.Context, tx *sql.Tx, id int64) (*model.UserToken, error) {
	userToken := new(model.UserToken)
	query := `
SELECT * FROM tokens WHERE user_id = $1;
	`
	err := tx.QueryRowContext(ctx, query, id).Scan(&userToken.UserID, &userToken.RefreshToken)
	if err != nil {
		return nil, err
	}
	return userToken, nil
}

func updateToken(ctx context.Context, tx *sql.Tx, userToken *model.UserToken) error {
	query := `
UPDATE tokens
SET refresh_token = $1
WHERE user_id = $2;
	`
	_, err := tx.ExecContext(ctx, query, userToken.RefreshToken, userToken.UserID)
	return err
}

func addUserToken(ctx context.Context, tx *sql.Tx, userToken *model.UserToken) error {
	query := `
INSERT INTO tokens (user_id, refresh_token)
VALUES ($1, $2);
	`
	_, err := tx.ExecContext(ctx, query, userToken.UserID, userToken.RefreshToken)
	return err
}
