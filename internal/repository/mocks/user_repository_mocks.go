// Code generated by MockGen. DO NOT EDIT.
// Source: user_repository.go
//
// Generated by this command:
//
//	mockgen -source=user_repository.go -destination=mocks/user_repository_mocks.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	repository "github.com/romakorinenko/task-manager/internal/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockIUserRepo is a mock of IUserRepo interface.
type MockIUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIUserRepoMockRecorder
}

// MockIUserRepoMockRecorder is the mock recorder for MockIUserRepo.
type MockIUserRepoMockRecorder struct {
	mock *MockIUserRepo
}

// NewMockIUserRepo creates a new mock instance.
func NewMockIUserRepo(ctrl *gomock.Controller) *MockIUserRepo {
	mock := &MockIUserRepo{ctrl: ctrl}
	mock.recorder = &MockIUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserRepo) EXPECT() *MockIUserRepoMockRecorder {
	return m.recorder
}

// BlockByID mocks base method.
func (m *MockIUserRepo) BlockByID(ctx context.Context, userID string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockByID", ctx, userID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// BlockByID indicates an expected call of BlockByID.
func (mr *MockIUserRepoMockRecorder) BlockByID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockByID", reflect.TypeOf((*MockIUserRepo)(nil).BlockByID), ctx, userID)
}

// Create mocks base method.
func (m *MockIUserRepo) Create(ctx context.Context, user *repository.User) *repository.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(*repository.User)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIUserRepoMockRecorder) Create(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIUserRepo)(nil).Create), ctx, user)
}

// GetAll mocks base method.
func (m *MockIUserRepo) GetAll(ctx context.Context) []repository.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]repository.User)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockIUserRepoMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockIUserRepo)(nil).GetAll), ctx)
}

// GetByLogin mocks base method.
func (m *MockIUserRepo) GetByLogin(ctx context.Context, userLogin string) (*repository.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByLogin", ctx, userLogin)
	ret0, _ := ret[0].(*repository.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByLogin indicates an expected call of GetByLogin.
func (mr *MockIUserRepoMockRecorder) GetByLogin(ctx, userLogin any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByLogin", reflect.TypeOf((*MockIUserRepo)(nil).GetByLogin), ctx, userLogin)
}
