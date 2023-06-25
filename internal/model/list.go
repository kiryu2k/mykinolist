package model

import (
	"fmt"
	"strings"
)

var titleStatus = [...]string{"watching", "completed", "on-hold", "dropped", "plan to watch"}

type ListInfo struct {
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
	Status     string `json:"status"`
	Score      uint8  `json:"score"`
	IsFavorite bool   `json:"is_favorite"`
	ListInfo   `json:"-"`
}

func (u *ListUnit) Validate() error {
	if u.Score > 10 {
		return fmt.Errorf("score cannot be greater than 10")
	}
	for _, status := range titleStatus {
		if strings.EqualFold(status, u.Status) {
			u.Status = status
			return nil
		}
	}
	return fmt.Errorf("invalid title status")
}
