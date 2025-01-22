package service

import (
	"context"
	"github.com/romakorinenko/task-manager/internal/repository"
)

type ITaskService interface {
	Create(ctx context.Context, task *repository.Task) (*repository.Task, error)
	Update(ctx context.Context, task *repository.Task) (*repository.Task, error)
	DeleteByID(ctx context.Context, taskID int) error
	GetByID(ctx context.Context, taskID int) (*repository.Task, error)
	GetByUserID(ctx context.Context, userID int) ([]repository.Task, error)
	GetByUserLogin(ctx context.Context, userLogin string) ([]repository.Task, error)
	GetAll(ctx context.Context) ([]repository.Task, error)
	GetByStatus(ctx context.Context, status string) ([]repository.Task, error)
	GetByPriority(ctx context.Context, priority int) ([]repository.Task, error)
}

type TaskService struct {
	taskRepository repository.ITaskRepo
}

func NewTaskService(taskRepository repository.ITaskRepo) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (t *TaskService) Create(ctx context.Context, task *repository.Task) (*repository.Task, error) {
	return t.taskRepository.Create(ctx, task)
}

func (t *TaskService) Update(ctx context.Context, task *repository.Task) (*repository.Task, error) {
	return t.taskRepository.Update(ctx, task)
}

func (t *TaskService) DeleteByID(ctx context.Context, taskID int) error {
	return t.taskRepository.DeleteByID(ctx, taskID)
}

func (t *TaskService) GetByID(ctx context.Context, taskID int) (*repository.Task, error) {
	return t.taskRepository.GetByID(ctx, taskID)
}

func (t *TaskService) GetByUserID(ctx context.Context, userID int) ([]repository.Task, error) {
	return t.taskRepository.GetByUserID(ctx, userID)
}

func (t *TaskService) GetByUserLogin(ctx context.Context, userLogin string) ([]repository.Task, error) {
	return t.taskRepository.GetByUserLogin(ctx, userLogin)
}

func (t *TaskService) GetAll(ctx context.Context) ([]repository.Task, error) {
	return t.taskRepository.GetAll(ctx)
}

func (t *TaskService) GetByStatus(ctx context.Context, status string) ([]repository.Task, error) {
	return t.taskRepository.GetByStatus(ctx, status)
}

func (t *TaskService) GetByPriority(ctx context.Context, priority int) ([]repository.Task, error) {
	return t.taskRepository.GetByPriority(ctx, priority)
}
