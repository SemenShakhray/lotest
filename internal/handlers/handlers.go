package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"lotest/internal/logger"
	"lotest/internal/models"

	"github.com/google/uuid"
)

type Handler struct {
	service Servicer
	Log     chan<- string
}

type Servicer interface {
	CreateTask(ctx context.Context, req models.RequestCreateTask) string
	GetAllTasks(ctx context.Context, status string) []models.Task
	GetTaskByID(ctx context.Context, id string) (models.Task, error)
}

func NewHandler(srv Servicer, logCh chan<- string) *Handler {
	return &Handler{
		service: srv,
		Log:     logCh,
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.RequestCreateTask

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Log <- logger.FormatLog(logger.Err, fmt.Sprintf("failed to decode request -%s", err.Error()))

		http.Error(w, `{"error":"failed to decode request"}`, http.StatusBadRequest)
		return
	}

	id := h.service.CreateTask(r.Context(), req)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"id": id}); err != nil {
		h.Log <- logger.FormatLog(logger.Err, fmt.Sprintf("failed to encode response -%s", err.Error()))

		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	tasks := h.service.GetAllTasks(r.Context(), status)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.Log <- logger.FormatLog(logger.Err, fmt.Sprintf("failed to encode response -%s", err.Error()))

		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	id := strings.TrimPrefix(path, "/task/")
	if id == "" {
		h.Log <- logger.FormatLog(logger.Err, "id is empty")

		http.Error(w, `{"error":"id is empty"}`, http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		h.Log <- logger.FormatLog(logger.Err, fmt.Sprintf("invalid id - %s", err.Error()))

		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrorTaskNotFound) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}

		http.Error(w, `{"error":"failed to get task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.Log <- logger.FormatLog(logger.Err, fmt.Sprintf("failed to encode response -%s", err.Error()))

		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
