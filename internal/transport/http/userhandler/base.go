package userhandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/taskservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/userservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/basehttphandler"
)

// UserHandler is the interface for user handler.
type UserHandler interface {
	AddRoutes(router fiber.Router)
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
}

// userHandler is the handler for http requests.
type userHandler struct {
	*basehttphandler.BaseHttpHandler
	userService userservice.UserService
	taskService taskservice.TaskService
}

// Option is the option type for user handler.
type Option func(*userHandler)

// WithBaseHttpHandler sets the base http handler option.
func WithBaseHttpHandler(handler *basehttphandler.BaseHttpHandler) Option {
	return func(h *userHandler) {
		h.BaseHttpHandler = handler
	}
}

// WithUserService sets the user service option.
func WithUserService(service userservice.UserService) Option {
	return func(h *userHandler) {
		h.userService = service
	}
}

// WithTaskService sets the task service option.
func WithTaskService(service taskservice.TaskService) Option {
	return func(h *userHandler) {
		h.taskService = service
	}
}

// New creates a new http handler with the given options.
func New(opts ...Option) UserHandler {
	h := &userHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
