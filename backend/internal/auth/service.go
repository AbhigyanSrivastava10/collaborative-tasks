package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/models"
)

type Service struct {
	db        *pgxpool.Pool
	jwtSecret string
}

func NewService(db *pgxpool.Pool, jwtSecret string) *Service {
	return &Service{db: db, jwtSecret: jwtSecret}
}

func (s *Service) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	var exists bool
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already in use")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Insert user
	var user models.User
	err = s.db.QueryRow(ctx,
		`INSERT INTO users (email, password, name, provider)
		 VALUES ($1, $2, $3, 'email')
		 RETURNING id, email, name, avatar_url, provider, created_at, updated_at`,
		req.Email, string(hashed), req.Name,
	).Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: user}, nil
}

func (s *Service) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	var user models.User
	var hashedPassword string

	err := s.db.QueryRow(ctx,
		`SELECT id, email, password, name, avatar_url, provider, created_at, updated_at
		 FROM users WHERE email=$1`,
		req.Email,
	).Scan(&user.ID, &user.Email, &hashedPassword, &user.Name, &user.AvatarURL, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: user}, nil
}

func (s *Service) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
