package service

import (
	"context"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type listService struct{}

func (s *listService) AddMovie(ctx context.Context, movie *model.Movie) error {
	return nil
}
