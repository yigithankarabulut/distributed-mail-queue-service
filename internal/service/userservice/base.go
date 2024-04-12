package userservice

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/res"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/userstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg"
)

type UserService interface {
	Register(ctx context.Context, req dtoreq.RegisterRequest) error
	Login(ctx context.Context, req dtoreq.LoginRequest) (dtores.LoginResponse, error)
	GetUser(ctx context.Context, req dtoreq.GetUserRequest) (dtores.GetUserResponse, error)
}

type userService struct {
	*pkg.Packages
	userStorage userstorage.UserStorer
	taskStorage taskstorage.TaskStorer
	mailService mailservice.MailService
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

func WithMailService(mailService mailservice.MailService) Option {
	return func(u *userService) {
		u.mailService = mailService
	}
}

func New(opts ...Option) UserService {
	service := &userService{}
	for _, opt := range opts {
		opt(service)
	}
	return service
}
