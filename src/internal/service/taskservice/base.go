package taskservice

import (
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/storage/taskstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/storage/userstorage"
)

type TaskService interface {
}

type taskService struct {
	taskStorage taskstorage.TaskStorer
	userStorage userstorage.UserStorer
}

type Option func(*taskService)

func WithTaskStorage(taskStorage taskstorage.TaskStorer) Option {
	return func(t *taskService) {
		t.taskStorage = taskStorage
	}
}

func WithUserStorage(userStorage userstorage.UserStorer) Option {
	return func(t *taskService) {
		t.userStorage = userStorage
	}
}

func New(opts ...Option) TaskService {
	service := &taskService{}
	for _, opt := range opts {
		opt(service)
	}
	return service
}
