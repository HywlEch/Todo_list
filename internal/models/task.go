package models

import (
	"time"
)

type Task struct {
	ID         int          `json:"id"`
	Title      string       `json:"title"`
	Content    string 		`json:"content"`
	Done       bool   		`json:"done"`
	CreatedAt  time.Time	`json:"created_at"`
	UpdatedAt  time.Time    `json:"update_at"`
}