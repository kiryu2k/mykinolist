package model

import "fmt"

const (
	Watching = iota + 1
	Completed
	OnHold
	Dropped
	PlanToWatch
)

type List struct {
	ListID  int64 `json:"list_id"`
	OwnerID int64 `json:"user_id"`
}

type SearchResult struct {
	Docs []Movie `json:"docs"`
}

type Movie struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ListUnit struct {
	Movie
	Status     int   `json:"status_id"`
	Score      uint8 `json:"score"`
	IsFavorite bool  `json:"is_favorite"`
	List       `json:"-"`
}

func (u *ListUnit) Validate() error {
	if u.Score > 10 {
		return fmt.Errorf("score cannot be greater than 10")
	}
	if u.Status < Watching || u.Status > PlanToWatch {
		return fmt.Errorf("invalid title status")
	}
	return nil
}
