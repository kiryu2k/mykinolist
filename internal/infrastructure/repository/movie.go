package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type movieRepository struct {
	db *sql.DB
}

func (r *movieRepository) Add(ctx context.Context, movie *model.ListUnit) error {
	query := `
INSERT INTO list_titles (list_id, title_id, status_id, score, is_favorite)
SELECT id, $1, $2, $3, $4
FROM lists WHERE id = $5;
	`
	res, err := r.db.ExecContext(ctx, query, movie.ID, movie.Status,
		movie.Score, movie.IsFavorite, movie.OwnerID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("invalid count of added titles %d", count)
	}
	return nil
}
