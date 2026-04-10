package board

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

func (s *Service) Create(ctx context.Context, ownerID string, req models.CreateBoardRequest) (*models.Board, error) {
	if req.Name == "" {
		return nil, errors.New("board name is required")
	}

	var board models.Board
	err := s.db.QueryRow(ctx,
		`INSERT INTO boards (name, description, owner_id)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, description, owner_id, created_at, updated_at`,
		req.Name, req.Description, ownerID,
	).Scan(&board.ID, &board.Name, &board.Description, &board.OwnerID, &board.CreatedAt, &board.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Auto-add owner as admin member
	_, err = s.db.Exec(ctx,
		`INSERT INTO board_members (board_id, user_id, role) VALUES ($1, $2, 'admin')`,
		board.ID, ownerID,
	)
	if err != nil {
		return nil, err
	}

	return &board, nil
}

func (s *Service) GetAll(ctx context.Context, userID string) ([]models.Board, error) {
	rows, err := s.db.Query(ctx,
		`SELECT b.id, b.name, b.description, b.owner_id, b.created_at, b.updated_at
		 FROM boards b
		 JOIN board_members bm ON bm.board_id = b.id
		 WHERE bm.user_id = $1
		 ORDER BY b.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []models.Board
	for rows.Next() {
		var b models.Board
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.OwnerID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, nil
}

func (s *Service) GetByID(ctx context.Context, boardID, userID string) (*models.Board, error) {
	var board models.Board
	err := s.db.QueryRow(ctx,
		`SELECT b.id, b.name, b.description, b.owner_id, b.created_at, b.updated_at
		 FROM boards b
		 JOIN board_members bm ON bm.board_id = b.id
		 WHERE b.id = $1 AND bm.user_id = $2`,
		boardID, userID,
	).Scan(&board.ID, &board.Name, &board.Description, &board.OwnerID, &board.CreatedAt, &board.UpdatedAt)
	if err != nil {
		return nil, errors.New("board not found")
	}
	return &board, nil
}

func (s *Service) Update(ctx context.Context, boardID, userID string, req models.UpdateBoardRequest) (*models.Board, error) {
	var board models.Board
	err := s.db.QueryRow(ctx,
		`UPDATE boards SET name=$1, description=$2, updated_at=NOW()
		 WHERE id=$3 AND owner_id=$4
		 RETURNING id, name, description, owner_id, created_at, updated_at`,
		req.Name, req.Description, boardID, userID,
	).Scan(&board.ID, &board.Name, &board.Description, &board.OwnerID, &board.CreatedAt, &board.UpdatedAt)
	if err != nil {
		return nil, errors.New("board not found or you are not the owner")
	}
	return &board, nil
}

func (s *Service) Delete(ctx context.Context, boardID, userID string) error {
	result, err := s.db.Exec(ctx,
		`DELETE FROM boards WHERE id=$1 AND owner_id=$2`,
		boardID, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("board not found or you are not the owner")
	}
	return nil
}
