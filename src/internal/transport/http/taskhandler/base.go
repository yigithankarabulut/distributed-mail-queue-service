package taskhandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/service/taskservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/service/userservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/transport/http/basehttphandler"
)

// Endpoints is the interface for http endpoints.
type Endpoints interface {
	AddRoutes(router fiber.Router)
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
func New(opts ...Option) Endpoints {
	h := &taskHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
