package userservice

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/dto"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/storage/taskstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/storage/userstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/pkg"
)

type UserService interface {
	Register(ctx context.Context, req dto.RegisterUserRequest) error
	GetUser(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error)
	UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (dto.UpdateUserResponse, error)
}

type userService struct {
	*pkg.Packages
	userStorage userstorage.UserStorer
	taskStorage taskstorage.TaskStorer
}

type Option func(*userService)

func WithUserStorage(userStorage userstorage.UserStorer) Option {
	return func(u *userService) {
		u.userStorage = userStorage
	}
}

func WithTaskStorage(taskStorage taskstorage.TaskStorer) Option {
	return func(u *userService) {
		u.taskStorage = taskStorage
	}
}

func WithPackages(packages *pkg.Packages) Option {
	return func(u *userService) {
		u.Packages = packages
	}
}

func New(opts ...Option) UserService {
	service := &userService{}
	for _, opt := range opts {
		opt(service)
	}
	return service
}
