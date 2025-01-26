package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/romakorinenko/task-manager/internal/errs"
	"github.com/romakorinenko/task-manager/internal/repository"
)

type IUserService interface {
	GetUserRepository() repository.IUserRepo
	Create(ctx context.Context, user *repository.User) error
	GetByLogin(ctx context.Context, userLogin string) *repository.User
	GetAll(ctx context.Context) []repository.User
}

type UserService struct {
	userRepository repository.IUserRepo
}

func NewUserService(userRepository repository.IUserRepo) *UserService {
	return &UserService{userRepository: userRepository}
}

func (u *UserService) GetUserRepository() repository.IUserRepo {
	return u.userRepository
}

func (u *UserService) Create(ctx context.Context, user *repository.User) error {
	userFromDB, err := u.userRepository.GetByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if userFromDB != nil {
		return errs.UserExistsErr{}
	}

	createdUser := u.userRepository.Create(ctx, user)
	if createdUser == nil {
		return errors.New("internal server error. user is not created")
	}

	return nil
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
