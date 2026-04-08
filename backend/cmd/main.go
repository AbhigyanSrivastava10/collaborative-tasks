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
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect to PostgreSQL
	pool := db.Connect(cfg.DatabaseURL)
	defer pool.Close()

	// Connect to Redis
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

	// Set up services and handlers
	authService := auth.NewService(pool, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	// Set up router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Auth routes
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	// Protected routes (we'll add more here in future commits)
	r.Group(func(r chi.Router) {
		r.Use(authService.Middleware)
		r.Get("/api/me", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(auth.UserIDKey)
			respondJSON(w, http.StatusOK, map[string]any{"user_id": userID})
		})
	})

	// Start server
	addr := ":" + cfg.Port
	fmt.Printf("🚀 Server running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, "%v", data)
}