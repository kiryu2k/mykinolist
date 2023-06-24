package repository

import (
	"context"
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type listRepository struct {
	db *sql.DB
}

func (r *listRepository) Create(ctx context.Context, ownerID int64) (*model.List, error) {
	list := new(model.List)
	list.OwnerID = ownerID
	query := `INSERT INTO lists (owner_id) VALUES ($1) RETURNING id;`
	err := r.db.QueryRowContext(ctx, query, ownerID).Scan(&list.ListID)
	if err != nil {
		return nil, err
	}
	return list, nil
}
