package main

import "time"

// TaskStatus enumerates allowed statuses.
type TaskStatus string

const (
	StatusPending   TaskStatus = "Pending"
	StatusInProgress TaskStatus = "InProgress"
	StatusCompleted TaskStatus = "Completed"
	StatusCancelled TaskStatus = "Cancelled"
)

// Task is the domain model.
type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
