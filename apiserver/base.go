package apiserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yigithankarabulut/distributed-mail-queue-service/config"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/taskservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/userservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/workerservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskqueue"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/userstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/basehttphandler"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/taskhandler"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/userhandler"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/cron"
	"log/slog"
)

type HttpEndpoints interface {
	AddRoutes(router fiber.Router)
}

type ApiServer interface {
	Run() error
}

type Instances struct {
	packages        *pkg.Packages
	taskQueue       taskqueue.TaskQueue
	userStorage     userstorage.UserStorer
	taskStorage     taskstorage.TaskStorer
	cronService     *cron.CronService
	userService     userservice.UserService
	taskService     taskservice.TaskService
	workers         []workerservice.IWorker
	basehttphandler *basehttphandler.BaseHttpHandler
	userHandler     userhandler.UserHandler
	taskHandler     taskhandler.TaskHandler
}

type apiServer struct {
	serverEnv   string
	app         *fiber.App
	logLevel    slog.Level
	logger      *slog.Logger
	config      *config.Config
	handlers    []HttpEndpoints
	instances   *Instances
	done        chan struct{}
	taskChannel chan model.MailTaskQueue
}

type Option func(*apiServer)

// WithLogLevel sets the log level option.
func WithLogLevel(level string) Option {
	return func(s *apiServer) {
		switch level {
		case "DEBUG":
			s.logLevel = slog.LevelDebug
		case "INFO":
			s.logLevel = slog.LevelInfo
		case "WARN":
			s.logLevel = slog.LevelWarn
		case "ERROR":
			s.logLevel = slog.LevelError
		default:
			s.logLevel = slog.LevelInfo
		}
	}
}

// WithLogger sets the logger option.
func WithLogger(logger *slog.Logger) Option {
	return func(s *apiServer) {
		s.logger = logger
	}
}

// WithServerEnv sets the server environment option.
func WithServerEnv(env string) Option {
	return func(s *apiServer) {
		s.serverEnv = env
	}
}

// WithConfig sets the config option.
func WithConfig(config *config.Config) Option {
	return func(s *apiServer) {
		s.config = config
	}
}

// New creates a new ApiServer instance with the given options.
func New(opts ...Option) ApiServer {
	apiserv := &apiServer{
		logLevel: slog.LevelInfo,
	}
	for _, opt := range opts {
		opt(apiserv)
	}
	return apiserv
}
