package service

import (
	"context"
	"errors"
	"github.com/romakorinenko/task-manager/internal/constant"
	"github.com/romakorinenko/task-manager/internal/errs"
	"github.com/romakorinenko/task-manager/internal/repository"
	"testing"

	mockRepository "github.com/romakorinenko/task-manager/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTaskService_GetTaskRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	taskRepository := taskService.GetTaskRepository()

	require.Equal(t, taskRepo, taskRepository)
}

func TestTaskService_Create_TaskCreated(t *testing.T) {
	background := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	taskService := NewTaskService(taskRepo, userRepo)

	userRepo.EXPECT().GetByLogin(gomock.Any(), gomock.Any()).Return(&repository.User{ID: 1}, nil)
	taskRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(1, nil)

	task, err := taskService.Create(background, 1, "Title", "Desc", "user")
	require.NoError(t, err)
	require.Equal(t, 1, task)
}

func TestTaskService_Create_PriorityInvalid(t *testing.T) {
	ctx := context.Background()
	taskService := NewTaskService(nil, nil)

	taskID, err := taskService.Create(ctx, 0, "Title", "Desc", "user")
	require.Equal(t, errs.BadReqErr{}, err)
	require.Equal(t, 0, taskID)
}

func TestTaskService_Create_UserLoginIsNotExists(t *testing.T) {
	background := context.Background()
	ctrl := gomock.NewController(t)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	taskService := NewTaskService(nil, userRepo)

	userRepo.EXPECT().GetByLogin(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

	taskID, err := taskService.Create(background, 1, "Title", "Desc", "user")
	require.Equal(t, errs.BadReqErr{}, err)
	require.Equal(t, 0, taskID)
}

func TestTaskService_Create_InternalErr(t *testing.T) {
	background := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	userRepo := mockRepository.NewMockIUserRepo(ctrl)
	taskService := NewTaskService(taskRepo, userRepo)

	userRepo.EXPECT().GetByLogin(gomock.Any(), gomock.Any()).Return(&repository.User{ID: 1}, nil)
	taskRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(0, errors.New(""))

	task, err := taskService.Create(background, 1, "Title", "Desc", "user")
	require.Error(t, err)
	require.Equal(t, 0, task)
}

func TestTaskService_GetAllByUser_UserTasksReturned(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	user := &repository.User{ID: 1, Role: constant.UserRole}
	taskRepo.EXPECT().GetTasksWithLoginByUserID(gomock.Any(), gomock.Any()).Return([]repository.TaskWithLogin{}, nil)

	tasks, err := taskService.GetAllByUser(ctx, user)
	require.NoError(t, err)
	require.Equal(t, 0, len(tasks))
}

func TestTaskService_GetAllByUser_UsersTaskReceivingError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	user := &repository.User{ID: 1, Role: constant.UserRole}
	taskRepo.EXPECT().GetTasksWithLoginByUserID(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

	tasks, err := taskService.GetAllByUser(ctx, user)
	require.Error(t, err)
	require.Nil(t, tasks)
}

func TestTaskService_GetAllByUser_AllTasksReturned(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	user := &repository.User{ID: 1, Role: constant.AdminRole}
	taskRepo.EXPECT().GetTasksWithLogin(gomock.Any()).Return([]repository.TaskWithLogin{}, nil)

	tasks, err := taskService.GetAllByUser(ctx, user)
	require.NoError(t, err)
	require.Equal(t, 0, len(tasks))
}

func TestTaskService_GetAllByUser_AllTaskReceivingError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	user := &repository.User{ID: 1, Role: constant.AdminRole}
	taskRepo.EXPECT().GetTasksWithLogin(gomock.Any()).Return(nil, errors.New(""))

	tasks, err := taskService.GetAllByUser(ctx, user)
	require.Error(t, err)
	require.Nil(t, tasks)
}

func TestTaskService_Update_TaskUpdated(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	user := &repository.Task{}
	taskRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(user, nil)
	taskRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	err := taskService.Update(ctx, "title", "desc", "OPEN", 1, 1)
	require.NoError(t, err)
}

func TestTaskService_Update_InvalidDescription(t *testing.T) {
	ctx := context.Background()
	taskService := NewTaskService(nil, nil)

	err := taskService.Update(ctx, "title", "", "OPEN", 1, 1)
	require.Error(t, err)
}

func TestTaskService_Update_TaskNotFound(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	taskRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

	err := taskService.Update(ctx, "title", "desc", "OPEN", 1, 1)
	require.Error(t, err)
}

func TestTaskService_Update_InternalErr(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	user := &repository.Task{}
	taskRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(user, nil)
	taskRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New(""))

	err := taskService.Update(ctx, "title", "desc", "OPEN", 1, 1)
	require.Error(t, err)
}

func TestTaskService_GetByStatus_TasksReturned(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	taskRepo.EXPECT().GetByStatus(gomock.Any(), gomock.Any()).Return([]repository.Task{}, nil)

	tasks, err := taskService.GetByStatus(ctx, "OPEN")
	require.NoError(t, err)
	require.Equal(t, 0, len(tasks))
}

func TestTaskService_GetByStatus_StatusIsEmpty(t *testing.T) {
	ctx := context.Background()
	taskService := NewTaskService(nil, nil)

	tasks, err := taskService.GetByStatus(ctx, "")
	require.Error(t, err)
	require.Nil(t, tasks)
}

func TestTaskService_GetByStatus_WrongStatus(t *testing.T) {
	ctx := context.Background()
	taskService := NewTaskService(nil, nil)

	tasks, err := taskService.GetByStatus(ctx, "PPPPP")
	require.Error(t, err)
	require.Nil(t, tasks)
}

func TestTaskService_GetByPriority_TasksReturned(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	taskRepo := mockRepository.NewMockITaskRepo(ctrl)
	taskService := NewTaskService(taskRepo, nil)

	taskRepo.EXPECT().GetByPriority(gomock.Any(), gomock.Any()).Return([]repository.Task{}, nil)

	tasks, err := taskService.GetByPriority(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, 0, len(tasks))
}

func TestTaskService_GetByPriority_WrongStatus(t *testing.T) {
	ctx := context.Background()
	taskService := NewTaskService(nil, nil)

	tasks, err := taskService.GetByPriority(ctx, 0)
	require.Error(t, err)
	require.Nil(t, tasks)
}
