package dto

import (
	"time"

	"github.com/romakorinenko/task-manager/internal/repository"
)

type TaskTemplateData struct {
	ID          int
	Title       string
	Description string
	Priority    int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserLogin   string
}

type TasksTemplateData struct {
	Tasks []TaskTemplateData
}

func TaskToTaskTemplateData(task *repository.Task, userLogin string) *TaskTemplateData {
	return &TaskTemplateData{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		UserLogin:   userLogin,
	}
}

type UsersTemplateData struct {
	Users []repository.User
}
