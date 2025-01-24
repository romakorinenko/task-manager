package service

import (
	"context"
	"time"

	"github.com/romakorinenko/task-manager/internal/constant"
	"github.com/romakorinenko/task-manager/internal/errs"
	"github.com/romakorinenko/task-manager/internal/repository"
)

type ITaskService interface {
	GetTaskRepository() repository.ITaskRepo
	Create(ctx context.Context, priority int, title, description, userLogin string) (int, error)
	Update(ctx context.Context,
		title, description, status string,
		priority, ID int,
	) error
	GetAllByUser(ctx context.Context, user *repository.User) ([]repository.TaskWithLogin, error)
	GetByStatus(ctx context.Context, status string) ([]repository.Task, error)
	GetByPriority(ctx context.Context, priority int) ([]repository.Task, error)
}

type TaskService struct {
	TaskRepository repository.ITaskRepo
	userRepository repository.IUserRepo
}

func NewTaskService(taskRepository repository.ITaskRepo, userRepository repository.IUserRepo) *TaskService {
	return &TaskService{
		TaskRepository: taskRepository,
		userRepository: userRepository,
	}
}

func (t *TaskService) GetTaskRepository() repository.ITaskRepo {
	return t.TaskRepository
}

func (t *TaskService) GetAllByUser(ctx context.Context, user *repository.User) ([]repository.TaskWithLogin, error) {
	userRole := user.Role
	if userRole == constant.AdminRole {
		return t.TaskRepository.GetTasksWithLogin(ctx)
	}

	return t.TaskRepository.GetTasksWithLoginByUserID(ctx, user.ID)
}

func (t *TaskService) Create(ctx context.Context, priority int, title, description, userLogin string) (int, error) {
	if title == "" || description == "" || userLogin == "" || priority < 1 || priority > 4 {
		return 0, errs.BadReqErr{}
	}

	user, err := t.userRepository.GetByLogin(ctx, userLogin)
	if err != nil {
		return 0, errs.BadReqErr{}
	}

	now := time.Now()
	taskForCreate := &repository.Task{
		Title:       title,
		Description: description,
		Priority:    priority,
		UserID:      user.ID,
		Status:      constant.OpenTaskStatus,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return t.TaskRepository.Create(ctx, taskForCreate)
}

func (t *TaskService) Update(ctx context.Context,
	title, description, status string,
	priority, ID int,
) error {
	if title == "" || description == "" || status == "" || priority < 1 || priority > 4 {
		return errs.BadReqErr{}
	}

	taskForUpdate, err := t.TaskRepository.GetByID(ctx, ID)
	if err != nil {
		return errs.BadReqErr{}
	}

	taskForUpdate.ID = ID
	taskForUpdate.Title = title
	taskForUpdate.Description = description
	taskForUpdate.Priority = priority
	taskForUpdate.Status = status

	return t.TaskRepository.Update(ctx, taskForUpdate)
}

func (t *TaskService) GetByStatus(ctx context.Context, status string) ([]repository.Task, error) {
	if status == "" {
		return nil, errs.BadReqErr{}
	}

	var isStatus bool
	for _, taskStatus := range constant.TaskStatuses {
		if taskStatus == status {
			isStatus = true
		}
	}
	if !isStatus {
		return nil, errs.BadReqErr{}
	}

	return t.TaskRepository.GetByStatus(ctx, status)
}

func (t *TaskService) GetByPriority(ctx context.Context, priority int) ([]repository.Task, error) {
	if priority < 1 || priority > 4 {
		return nil, errs.BadReqErr{}
	}

	return t.TaskRepository.GetByPriority(ctx, priority)
}
