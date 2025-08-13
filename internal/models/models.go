package models

import "errors"

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

var ErrorTaskNotFound = errors.New("task not found")

type Task struct {
	ID          string     `json:"id"`
	Status      TaskStatus `json:"status"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
}

type RequestCreateTask struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}
