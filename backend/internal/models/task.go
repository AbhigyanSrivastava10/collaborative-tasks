package models

import "time"

type Task struct {
	ID          string     `json:"id"`
	BoardID     string     `json:"board_id"`
	AssignedTo  *string    `json:"assigned_to"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Position    int        `json:"position"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssignedTo  *string    `json:"assigned_to"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Position    int        `json:"position"`
	AssignedTo  *string    `json:"assigned_to"`
	DueDate     *time.Time `json:"due_date"`
}

