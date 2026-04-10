package task

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/models"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	task, err := h.service.Create(r.Context(), boardID, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetByBoard(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")
	tasks, err := h.service.GetByBoard(r.Context(), boardID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}
	respondJSON(w, http.StatusOK, tasks)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	task, err := h.service.Update(r.Context(), taskID, req)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, task)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if err := h.service.Delete(r.Context(), taskID); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Task deleted"})
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
