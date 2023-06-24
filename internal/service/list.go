package service

import (
	"context"
	"fmt"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type listService struct {
	searcher MovieSearcher
	repo     MovieRepositroy
}

type MovieSearcher interface {
	Search(context.Context, string) (*model.SearchResult, error)
}

type MovieRepositroy interface {
	Add(context.Context, *model.ListUnit) error
}

/* Add the first found movie by the specified title to the [kino]list */
func (s *listService) AddMovie(ctx context.Context, movie *model.ListUnit) error {
	if len(movie.Name) == 0 {
		return fmt.Errorf("empty movie name")
	}
	searchResult, err := s.searcher.Search(ctx, movie.Name)
	if err != nil {
		return err
	}
	movie.Movie = searchResult.Docs[0]
	return s.repo.Add(ctx, movie)
}
