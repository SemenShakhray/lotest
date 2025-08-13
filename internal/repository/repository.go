package repository

import (
	"context"
	"lotest/internal/models"
	"lotest/internal/service"
	"strings"

	"github.com/google/uuid"
)

type Repo struct {
	Tasks map[string]models.Task
}

func NewRepo() service.Repositorer {
	return &Repo{
		Tasks: make(map[string]models.Task),
	}
}

func (r *Repo) CreateTask(ctx context.Context, req models.RequestCreateTask) string {
	id := uuid.New().String()

	r.Tasks[id] = models.Task{
		ID:          id,
		Status:      models.StatusPending,
		Title:       req.Title,
		Description: req.Description,
	}

	return id
}

func (r *Repo) GetAllTasks(ctx context.Context, status string) []models.Task {
	var tasks []models.Task

	if status == "" {
		for _, task := range r.Tasks {
			tasks = append(tasks, task)
		}

		return tasks
	}

	for _, task := range r.Tasks {
		if strings.Contains(string(task.Status), status) {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (r *Repo) GetTaskByID(ctx context.Context, id string) (models.Task, error) {
	task, ok := r.Tasks[id]
	if !ok {
		return models.Task{}, models.ErrorTaskNotFound
	}

	return task, nil
}
