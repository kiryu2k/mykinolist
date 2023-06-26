package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kiryu-dev/mykinolist/internal/model"
)

type listService struct {
	searcher MovieSearcher
	movie    MovieRepositroy
	list     ListRepository
}

type MovieSearcher interface {
	Search(context.Context, string) (*model.SearchResult, error)
	SearchByID(context.Context, int64) (*model.Movie, error)
}

type MovieRepositroy interface {
	Add(context.Context, *model.ListUnit) error
	GetAll(context.Context, int64) ([]*model.ListUnit, error)
	GetByID(context.Context, *model.ListUnit) error
	Update(context.Context, *model.ListUnitPatch) error
	Delete(context.Context, *model.ListUnit) error
}

/* Add the first found movie by the specified title to the [kino]list */
func (s *listService) AddMovie(movie *model.ListUnit) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if len(movie.Name) == 0 {
		return fmt.Errorf("empty movie name")
	}
	searchResult, err := s.searcher.Search(ctx, movie.Name)
	if err != nil {
		return err
	}
	movie.Movie = searchResult.Docs[0]
	return s.movie.Add(ctx, movie)
}

func (s *listService) GetMovies(userID int64) ([]*model.ListUnit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	movies, err := s.movie.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i, movieInfo := range movies {
		movie, err := s.searcher.SearchByID(ctx, movieInfo.ID)
		if err != nil {
			return nil, err
		}
		movies[i].Name = movie.Name
	}
	return movies, nil
}

func (s *listService) UpdateMovie(movie *model.ListUnitPatch) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := movie.Validate(); err != nil {
		return err
	}
	listID, err := s.list.GetID(ctx, *movie.OwnerID)
	if err != nil {
		return err
	}
	movie.ListID = &listID
	return s.movie.Update(ctx, movie)
}

func (s *listService) DeleteMovie(movie *model.ListUnit) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var (
		errChan = make(chan error, 2)
		wg      = &sync.WaitGroup{}
	)
	wg.Add(2)
	go func() {
		errChan <- s.movie.GetByID(ctx, movie)
		wg.Done()
	}()
	go func() {
		movieInfo, err := s.searcher.SearchByID(ctx, movie.ID)
		if err == nil {
			movie.Name = movieInfo.Name
		}
		errChan <- err
		wg.Done()
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return s.movie.Delete(ctx, movie)
}
