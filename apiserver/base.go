package apiserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yigithankarabulut/distributed-mail-queue-service/config"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"log/slog"
)

type HttpEndpoints interface {
	AddRoutes(router fiber.Router)
}

type ApiServer interface {
	Run() error
}

type apiServer struct {
	config      *config.Config
	app         *fiber.App
	handlers    []HttpEndpoints
	logLevel    slog.Level
	logger      *slog.Logger
	serverEnv   string
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
