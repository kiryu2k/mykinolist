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
	_, err := r.findTokenByID(ctx, userToken.UserID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		/* there's no user in token table, so we need to add him */
		err = r.addUserToken(ctx, userToken)
	} else {
		err = r.updateToken(ctx, userToken)
	}
	return err
}

func (r *tokenRepository) findTokenByID(ctx context.Context, id int64) (*model.UserToken, error) {
	userToken := new(model.UserToken)
	query := `
SELECT * FROM tokens WHERE user_id = $1;
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&userToken.UserID, &userToken.RefreshToken)
	if err != nil {
		return nil, err
	}
	return userToken, nil
}

func (r *tokenRepository) updateToken(ctx context.Context, userToken *model.UserToken) error {
	query := `
UPDATE tokens
SET refresh_token = $1
WHERE user_id = $2;
	`
	_, err := r.db.ExecContext(ctx, query, userToken.RefreshToken, userToken.UserID)
	return err
}

func (r *tokenRepository) addUserToken(ctx context.Context, userToken *model.UserToken) error {
	query := `
INSERT INTO tokens (user_id, refresh_token)
VALUES ($1, $2);
	`
	_, err := r.db.ExecContext(ctx, query, userToken.UserID, userToken.RefreshToken)
	return err
}
