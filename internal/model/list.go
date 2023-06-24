package model

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

type Movie struct {
	List
	Title      string `json:"title"`
	Status     int    `json:"status_id"`
	Score      uint8  `json:"score"`
	IsFavorite bool   `json:"is_favorite"`
}
