// Code generated by MockGen. DO NOT EDIT.
// Source: task_repository.go
//
// Generated by this command:
//
//	mockgen -source=task_repository.go -destination=mocks/task_repository_mocks.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	repository "github.com/romakorinenko/task-manager/internal/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockITaskRepo is a mock of ITaskRepo interface.
type MockITaskRepo struct {
	ctrl     *gomock.Controller
	recorder *MockITaskRepoMockRecorder
}

// MockITaskRepoMockRecorder is the mock recorder for MockITaskRepo.
type MockITaskRepoMockRecorder struct {
	mock *MockITaskRepo
}

// NewMockITaskRepo creates a new mock instance.
func NewMockITaskRepo(ctrl *gomock.Controller) *MockITaskRepo {
	mock := &MockITaskRepo{ctrl: ctrl}
	mock.recorder = &MockITaskRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITaskRepo) EXPECT() *MockITaskRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockITaskRepo) Create(ctx context.Context, task *repository.Task) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, task)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockITaskRepoMockRecorder) Create(ctx, task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockITaskRepo)(nil).Create), ctx, task)
}

// DeleteByID mocks base method.
func (m *MockITaskRepo) DeleteByID(ctx context.Context, taskID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockITaskRepoMockRecorder) DeleteByID(ctx, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockITaskRepo)(nil).DeleteByID), ctx, taskID)
}

// GetByID mocks base method.
func (m *MockITaskRepo) GetByID(ctx context.Context, taskID int) (*repository.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, taskID)
	ret0, _ := ret[0].(*repository.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockITaskRepoMockRecorder) GetByID(ctx, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockITaskRepo)(nil).GetByID), ctx, taskID)
}

// GetByPriority mocks base method.
func (m *MockITaskRepo) GetByPriority(ctx context.Context, priority int) ([]repository.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPriority", ctx, priority)
	ret0, _ := ret[0].([]repository.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPriority indicates an expected call of GetByPriority.
func (mr *MockITaskRepoMockRecorder) GetByPriority(ctx, priority any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPriority", reflect.TypeOf((*MockITaskRepo)(nil).GetByPriority), ctx, priority)
}

// GetByStatus mocks base method.
func (m *MockITaskRepo) GetByStatus(ctx context.Context, status string) ([]repository.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByStatus", ctx, status)
	ret0, _ := ret[0].([]repository.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByStatus indicates an expected call of GetByStatus.
func (mr *MockITaskRepoMockRecorder) GetByStatus(ctx, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByStatus", reflect.TypeOf((*MockITaskRepo)(nil).GetByStatus), ctx, status)
}

// GetByUserLogin mocks base method.
func (m *MockITaskRepo) GetByUserLogin(ctx context.Context, userLogin string) ([]repository.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserLogin", ctx, userLogin)
	ret0, _ := ret[0].([]repository.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserLogin indicates an expected call of GetByUserLogin.
func (mr *MockITaskRepoMockRecorder) GetByUserLogin(ctx, userLogin any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserLogin", reflect.TypeOf((*MockITaskRepo)(nil).GetByUserLogin), ctx, userLogin)
}

// GetTaskWithLoginByID mocks base method.
func (m *MockITaskRepo) GetTaskWithLoginByID(ctx context.Context, taskID int) (*repository.TaskWithLogin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskWithLoginByID", ctx, taskID)
	ret0, _ := ret[0].(*repository.TaskWithLogin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskWithLoginByID indicates an expected call of GetTaskWithLoginByID.
func (mr *MockITaskRepoMockRecorder) GetTaskWithLoginByID(ctx, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskWithLoginByID", reflect.TypeOf((*MockITaskRepo)(nil).GetTaskWithLoginByID), ctx, taskID)
}

// GetTasksWithLogin mocks base method.
func (m *MockITaskRepo) GetTasksWithLogin(ctx context.Context) ([]repository.TaskWithLogin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksWithLogin", ctx)
	ret0, _ := ret[0].([]repository.TaskWithLogin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksWithLogin indicates an expected call of GetTasksWithLogin.
func (mr *MockITaskRepoMockRecorder) GetTasksWithLogin(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksWithLogin", reflect.TypeOf((*MockITaskRepo)(nil).GetTasksWithLogin), ctx)
}

// GetTasksWithLoginByUserID mocks base method.
func (m *MockITaskRepo) GetTasksWithLoginByUserID(ctx context.Context, userID int) ([]repository.TaskWithLogin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksWithLoginByUserID", ctx, userID)
	ret0, _ := ret[0].([]repository.TaskWithLogin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksWithLoginByUserID indicates an expected call of GetTasksWithLoginByUserID.
func (mr *MockITaskRepoMockRecorder) GetTasksWithLoginByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksWithLoginByUserID", reflect.TypeOf((*MockITaskRepo)(nil).GetTasksWithLoginByUserID), ctx, userID)
}

// Update mocks base method.
func (m *MockITaskRepo) Update(ctx context.Context, task *repository.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockITaskRepoMockRecorder) Update(ctx, task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockITaskRepo)(nil).Update), ctx, task)
}
