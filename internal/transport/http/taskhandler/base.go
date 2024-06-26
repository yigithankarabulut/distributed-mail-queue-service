package taskhandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/taskservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/userservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/basehttphandler"
)

// TaskHandler is the interface for task handler.
type TaskHandler interface {
	AddRoutes(router fiber.Router)
	EnqueueTask(c *fiber.Ctx) error
	GetAllQueuedTasks(c *fiber.Ctx) error
	GetAllFailedQueuedTasks(c *fiber.Ctx) error
}

// taskHandler is the handler for http requests.
type taskHandler struct {
	*basehttphandler.BaseHttpHandler
	userService userservice.UserService
	taskService taskservice.TaskService
}

// Option is the option type for task handler.
type Option func(*taskHandler)

// WithBaseHttpHandler sets the base http handler option.
func WithBaseHttpHandler(handler *basehttphandler.BaseHttpHandler) Option {
	return func(h *taskHandler) {
		h.BaseHttpHandler = handler
	}
}

// WithUserService sets the user service option.
func WithUserService(service userservice.UserService) Option {
	return func(h *taskHandler) {
		h.userService = service
	}
}

// WithTaskService sets the task service option.
func WithTaskService(service taskservice.TaskService) Option {
	return func(h *taskHandler) {
		h.taskService = service
	}
}

// New creates a new http handler with the given options.
func New(opts ...Option) TaskHandler {
	h := &taskHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
