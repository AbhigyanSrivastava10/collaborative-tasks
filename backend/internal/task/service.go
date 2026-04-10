package task

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/models"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, boardID string, req models.CreateTaskRequest) (*models.Task, error) {
	if req.Title == "" {
		return nil, errors.New("task title is required")
	}
	if req.Status == "" {
		req.Status = "todo"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	var task models.Task
	err := s.db.QueryRow(ctx,
		`INSERT INTO tasks (board_id, title, description, status, priority, assigned_to, due_date)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, board_id, assigned_to, title, description, status, priority, position, due_date, created_at, updated_at`,
		boardID, req.Title, req.Description, req.Status, req.Priority, req.AssignedTo, req.DueDate,
	).Scan(&task.ID, &task.BoardID, &task.AssignedTo, &task.Title, &task.Description,
		&task.Status, &task.Priority, &task.Position, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *Service) GetByBoard(ctx context.Context, boardID string) ([]models.Task, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, board_id, assigned_to, title, description, status, priority, position, due_date, created_at, updated_at
		 FROM tasks WHERE board_id=$1 ORDER BY position ASC, created_at ASC`,
		boardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.BoardID, &t.AssignedTo, &t.Title, &t.Description,
			&t.Status, &t.Priority, &t.Position, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Service) Update(ctx context.Context, taskID string, req models.UpdateTaskRequest) (*models.Task, error) {
	var task models.Task
	err := s.db.QueryRow(ctx,
		`UPDATE tasks
		 SET title=$1, description=$2, status=$3, priority=$4, position=$5, assigned_to=$6, due_date=$7, updated_at=NOW()
		 WHERE id=$8
		 RETURNING id, board_id, assigned_to, title, description, status, priority, position, due_date, created_at, updated_at`,
		req.Title, req.Description, req.Status, req.Priority, req.Position, req.AssignedTo, req.DueDate, taskID,
	).Scan(&task.ID, &task.BoardID, &task.AssignedTo, &task.Title, &task.Description,
		&task.Status, &task.Priority, &task.Position, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, errors.New("task not found")
	}
	return &task, nil
}

func (s *Service) Delete(ctx context.Context, taskID string) error {
	result, err := s.db.Exec(ctx, `DELETE FROM tasks WHERE id=$1`, taskID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("task not found")
	}
	return nil
}
