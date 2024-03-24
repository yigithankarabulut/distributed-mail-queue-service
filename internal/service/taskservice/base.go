package taskservice

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/res"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskqueue"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/userstorage"
)

type TaskService interface {
	EnqueueMailTask(ctx context.Context, request dtoreq.TaskEnqueueRequest) (dtores.TaskEnqueueResponse, error)
	GetAllQueuedTasks(ctx context.Context, request dtoreq.GetAllQueuedTasksRequest) (dtores.GetAllQueuedTasksResponse, error)
	GetAllFailedQueuedTasks(ctx context.Context, request dtoreq.GetAllFailedTasksRequest) (dtores.GetAllFailedTasksResponse, error)
}

type taskService struct {
	taskStorage taskstorage.TaskStorer
	userStorage userstorage.UserStorer
	redisClient taskqueue.TaskQueue
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

func WithRedisClient(redisClient taskqueue.TaskQueue) Option {
	return func(t *taskService) {
		t.redisClient = redisClient
	}
}

func New(opts ...Option) TaskService {
	service := &taskService{}
	for _, opt := range opts {
		opt(service)
	}
	return service
}
