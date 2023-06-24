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

func (r *movieRepository) Add(ctx context.Context, movie *model.Movie) error {
	query := `
INSERT INTO list_titles (list_id, title, status_id, score, is_favorite)
VALUES ($1, $2, $3, $4, $5);
	`
	res, err := r.db.ExecContext(ctx, query, movie.ListID, movie.Title, movie.Status,
		movie.Score, movie.IsFavorite)
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
