package entities

import "time"

type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

type Tasks []Task
