package repository

import (
	"context"
	"database/sql"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type listRepository struct {
	db *sql.DB
}

func (r *listRepository) Create(ctx context.Context, ownerID int64) (*model.ListInfo, error) {
	list := new(model.ListInfo)
	list.OwnerID = ownerID
	query := `INSERT INTO lists (owner_id) VALUES ($1) RETURNING id;`
	err := r.db.QueryRowContext(ctx, query, ownerID).Scan(&list.ListID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *listRepository) GetID(ctx context.Context, ownerID int64) (int64, error) {
	query := `SELECT id FROM lists WHERE owner_id = $1;`
	var id int64
	err := r.db.QueryRowContext(ctx, query, ownerID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
