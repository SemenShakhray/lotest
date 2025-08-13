package service

import (
	"context"
	"fmt"
	"lotest/internal/handlers"
	"lotest/internal/logger"
	"lotest/internal/models"
)

type Service struct {
	repo Repositorer
	log  chan<- string
}

func NewService(repo Repositorer, log chan<- string) handlers.Servicer {
	return &Service{
		repo: repo,
		log:  log,
	}
}

type Repositorer interface {
	CreateTask(ctx context.Context, req models.RequestCreateTask) string
	GetAllTasks(ctx context.Context, status string) []models.Task
	GetTaskByID(ctx context.Context, id string) (models.Task, error)
}

func (s *Service) CreateTask(ctx context.Context, req models.RequestCreateTask) string {
	s.log <- logger.FormatLog(logger.Info, "creating task")

	if req.Title == "" {
		s.log <- logger.FormatLog(logger.Warn, "title is empty")
	}

	if req.Description == "" {
		s.log <- logger.FormatLog(logger.Warn, "description is empty")
	}

	taskID := s.repo.CreateTask(ctx, req)

	s.log <- logger.FormatLog(logger.Info, fmt.Sprintf("task created with id %s", taskID))

	return taskID
}

func (s *Service) GetAllTasks(ctx context.Context, status string) []models.Task {
	s.log <- logger.FormatLog(logger.Info, fmt.Sprintf("getting all tasks with status [%s]", status))

	tasks := s.repo.GetAllTasks(ctx, status)

	s.log <- logger.FormatLog(logger.Info, "got all tasks")

	return tasks
}

func (s *Service) GetTaskByID(ctx context.Context, id string) (models.Task, error) {
	s.log <- logger.FormatLog(logger.Info, fmt.Sprintf("getting task with id %s", id))

	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		s.log <- logger.FormatLog(logger.Err, fmt.Sprintf("failed to get task by id - %s", err.Error()))

		return models.Task{}, err
	}

	s.log <- logger.FormatLog(logger.Info, fmt.Sprintf("got task with id %s", id))

	return task, nil
}
