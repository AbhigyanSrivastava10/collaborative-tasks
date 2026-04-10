package board

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/auth"
	"github.com/abhigyansrivastava10/collaborative-tasks/backend/internal/models"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	var req models.CreateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	board, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, board)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	boards, err := h.service.GetAll(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch boards")
		return
	}
	respondJSON(w, http.StatusOK, boards)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	boardID := chi.URLParam(r, "id")
	board, err := h.service.GetByID(r.Context(), boardID, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, board)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	boardID := chi.URLParam(r, "id")
	var req models.UpdateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	board, err := h.service.Update(r.Context(), boardID, userID, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, board)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	boardID := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), boardID, userID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Board deleted"})
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
