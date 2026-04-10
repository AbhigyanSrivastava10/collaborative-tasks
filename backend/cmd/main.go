package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"

	"github.com/abhigyansrivastava10/collaborative-tasks/backend/config"
	"github.com/abhigyansrivastava10/collaborative-tasks/backend/db"
	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/auth"
	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/board"
	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/task"
	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/ws"
)

func main() {
	cfg := config.Load()

	pool := db.Connect(cfg.DatabaseURL)
	defer pool.Close()

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v\n", err)
	}
	redisClient := redis.NewClient(redisOpts)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Unable to connect to Redis: %v\n", err)
	}
	fmt.Println("✅ Connected to Redis")
	defer redisClient.Close()

	// Services
	authService := auth.NewService(pool, cfg.JWTSecret)
	boardService := board.NewService(pool)
	taskService := task.NewService(pool)

	// WebSocket hub
	hub := ws.NewHub(redisClient)
	go hub.Run(context.Background())

	// Handlers
	authHandler := auth.NewHandler(authService)
	boardHandler := board.NewHandler(boardService)
	taskHandler := task.NewHandler(taskService)
	wsHandler := ws.NewHandler(hub)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Public routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authService.Middleware)

		// Boards
		r.Get("/api/boards", boardHandler.GetAll)
		r.Post("/api/boards", boardHandler.Create)
		r.Get("/api/boards/{id}", boardHandler.GetByID)
		r.Put("/api/boards/{id}", boardHandler.Update)
		r.Delete("/api/boards/{id}", boardHandler.Delete)

		// Tasks
		r.Get("/api/boards/{boardID}/tasks", taskHandler.GetByBoard)
		r.Post("/api/boards/{boardID}/tasks", taskHandler.Create)
		r.Put("/api/boards/{boardID}/tasks/{taskID}", taskHandler.Update)
		r.Delete("/api/boards/{boardID}/tasks/{taskID}", taskHandler.Delete)

		// WebSocket
		r.Get("/ws/boards/{boardID}", wsHandler.ServeWS)
	})

	addr := ":" + cfg.Port
	fmt.Printf("🚀 Server running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
