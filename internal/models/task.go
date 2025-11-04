package models

import (
	"time"
)

type Task struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Done      bool      `json:"done" db:"done"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserID    int       `json:"user_id" db:"user_id"` 
}