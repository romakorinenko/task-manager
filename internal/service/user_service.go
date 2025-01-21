package service

import (
	"context"
	"errors"
	"github.com/romakorinenko/task-manager/internal/repository"
)

type IUserService interface {
	Create(ctx context.Context, user *repository.User) error
	BlockByID(ctx context.Context, userID string) bool
	GetByID(ctx context.Context, userID int) *repository.User
	GetByLogin(ctx context.Context, userLogin string) *repository.User
	GetAll(ctx context.Context) []repository.User
}

type UserService struct {
	userRepository repository.IUserRepo
}

func NewUserService(userRepository repository.IUserRepo) *UserService {
	return &UserService{userRepository: userRepository}
}

func (u *UserService) Create(ctx context.Context, user *repository.User) error {
	if u.userRepository.Create(ctx, user) == nil {
		return errors.New("user has not been created")
	}
	return nil
}

func (u *UserService) BlockByID(ctx context.Context, userID string) bool {

	return u.userRepository.BlockByID(ctx, userID)
}

func (u *UserService) GetByID(ctx context.Context, userID int) *repository.User {
	return u.userRepository.GetByID(ctx, userID)
}

func (u *UserService) GetByLogin(ctx context.Context, userLogin string) *repository.User {
	user, err := u.userRepository.GetByLogin(ctx, userLogin)
	if err != nil {
		return nil
	}

	return user
}

func (u *UserService) GetAll(ctx context.Context) []repository.User {
	return u.userRepository.GetAll(ctx)
}
