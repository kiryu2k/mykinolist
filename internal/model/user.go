package model

import (
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	CreatedOn      time.Time `json:"created_on"`
	LastLogin      time.Time `json:"last_login"`
}
