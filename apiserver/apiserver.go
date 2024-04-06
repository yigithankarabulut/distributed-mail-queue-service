package apiserver

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
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
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/constant"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/postgres"
	redisclient "github.com/yigithankarabulut/distributed-mail-queue-service/pkg/redis"
	"github.com/yigithankarabulut/distributed-mail-queue-service/releaseinfo"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	QueueConsumerCount = 10
	WorkerCount        = 10
)

type HttpEndpoints interface {
	AddRoutes(router fiber.Router)
}

type apiServer struct {
	config    *config.Config
	app       *fiber.App
	handlers  []HttpEndpoints
	logLevel  slog.Level
	logger    *slog.Logger
	serverEnv string
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

// NewApiServer creates a new api server instance with the given options.
func NewApiServer(opts ...Option) error {
	apiserv := &apiServer{
		logLevel: slog.LevelInfo,
	}
	for _, opt := range opts {
		opt(apiserv)
	}
	if apiserv.config == nil {
		return fmt.Errorf("config is required")
	}

	if err := connectStorages(apiserv); err != nil {
		return fmt.Errorf("error connecting to storages: %w", err)
	}
	if err := initializeServer(apiserv); err != nil {
		return fmt.Errorf("error initializing server: %w", err)
	}

	initializeApp(apiserv)
	healthzCheck(apiserv)
	appendDepends(apiserv)
	return listenAndServe(apiserv)
}

// connectStorages connects to the storages and returns an error if any.
func connectStorages(apiserv *apiServer) error {
	if _, err := postgres.ConnectPQ(apiserv.config.Database); err != nil {
		return fmt.Errorf("error connecting to postgres: %w", err)
	}
	if _, err := redisclient.New(apiserv.config.Redis); err != nil {
		return fmt.Errorf("error connecting to redis: %w", err)
	}
	return nil
}

// initializeServer initializes the server with the given config and logger.
func initializeServer(apiserv *apiServer) error {
	if apiserv.logger == nil {
		logHandlerOpts := &slog.HandlerOptions{Level: apiserv.logLevel}
		logHandler := slog.NewJSONHandler(os.Stdout, logHandlerOpts)
		apiserv.logger = slog.New(logHandler)
	}
	slog.SetDefault(apiserv.logger)
	if apiserv.serverEnv == "" {
		apiserv.serverEnv = "development"
	}
	return nil
}

// initializeApp initializes the fiber app with the given logger and adds the http logging middleware.
func initializeApp(apiserv *apiServer) {
	var (
		corsConfig cors.Config
	)
	apiserv.app = fiber.New(fiber.Config{
		ReadTimeout:  constant.ServerReadTimeout,
		WriteTimeout: constant.ServerWriteTimeout,
		IdleTimeout:  constant.ServerIdleTimeout,
	})

	corsConfig.AllowOrigins = constant.AllowedOrigins
	corsConfig.AllowCredentials = false
	corsConfig.AllowHeaders = strings.Join(
		[]string{
			constant.ContentType,
			constant.Authorization,
		}, ",")
	corsConfig.AllowMethods = strings.Join(
		[]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}, ",")
	apiserv.app.Use(cors.New(corsConfig))
	apiserv.app.Use(recover.New())
	apiserv.app.Use(httpLoggingMiddleware(apiserv.logger, apiserv.app))
}

// healthzCheck adds health endpoints to the apiserver.
func healthzCheck(apiserv *apiServer) {
	apiserv.app.Get("/healthz/live", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"server":            apiserv.serverEnv,
			"version":           releaseinfo.Version,
			"build_information": releaseinfo.BuildInformation,
			"message":           "liveness is OK!, server is ready to accept connections",
		})
	})
	apiserv.app.Get("/healthz/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "pong",
		})
	})
}

// appendDepends create instances from all layers and append them together. Then add it to the fiber by calling the AddRoutes method of all its handlers.
func appendDepends(apiserv *apiServer) {
	packages := pkg.New()
	TaskQueueChannel := make(chan model.MailTaskQueue, 100)

	userStorage := userstorage.New(userstorage.WithUserDB(postgres.DB))
	taskStorage := taskstorage.New(taskstorage.WithTaskDB(postgres.DB))
	taskQueue := taskqueue.New(
		taskqueue.WithTaskChannel(TaskQueueChannel),
		taskqueue.WithConsumerCount(QueueConsumerCount),
		taskqueue.WithQueueName(constant.RedisMailQueueChannel),
		taskqueue.WithRedisClient(redisclient.GetRedisClient()),
	)

	userService := userservice.New(
		userservice.WithUserStorage(userStorage),
		userservice.WithTaskStorage(taskStorage),
		userservice.WithPackages(packages),
	)
	taskService := taskservice.New(
		taskservice.WithTaskStorage(taskStorage),
		taskservice.WithUserStorage(userStorage),
		taskservice.WithRedisClient(taskQueue),
	)
	if err := taskService.FindUnprocessedTasksAndEnqueue(context.Background()); err != nil {
		apiserv.logger.Error("error finding unprocessed tasks", "error", err)
	}

	if err := taskQueue.StartConsume(); err != nil {
		apiserv.logger.Error("error starting consume", "error", err)
	}
	s, err := gocron.NewScheduler()
	if err != nil {
		apiserv.logger.Error("error creating scheduler", "error", err)
	}

	cronJob := &CronJob{
		JobName:   "FindUnprocessedTasksAndEnqueue",
		Scheduler: s,
		Task:      gocron.NewTask(taskService.FindUnprocessedTasksAndEnqueue, context.TODO()),
		Duration:  gocron.DurationJob(10 * time.Second),
	}
	jobID := cronJob.AddJob()
	if jobID == uuid.Nil {
		apiserv.logger.Error("error adding job")
	}
	log.Infof("job id: %s", jobID)
	s.Start()

	workers := make([]workerservice.IWorker, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		workers[i] = workerservice.New(
			workerservice.WithID(i+1),
			workerservice.WithTaskStorage(taskStorage),
			workerservice.WithTaskQueue(taskQueue),
			workerservice.WithChannel(TaskQueueChannel),
		)
	}
	for _, worker := range workers {
		go worker.TriggerWorker()
	}

	baseHttpHandler := basehttphandler.New(
		basehttphandler.WithContextTimeout(constant.ContextCancelTimeout),
		basehttphandler.WithLogger(apiserv.logger),
		basehttphandler.WithPackages(packages),
	)
	userHandler := userhandler.New(
		userhandler.WithBaseHttpHandler(baseHttpHandler),
		userhandler.WithUserService(userService),
		userhandler.WithTaskService(taskService),
	)
	taskHandler := taskhandler.New(
		taskhandler.WithBaseHttpHandler(baseHttpHandler),
		taskhandler.WithTaskService(taskService),
		taskhandler.WithUserService(userService),
	)

	apiserv.handlers = append(apiserv.handlers, userHandler, taskHandler)
	for _, handler := range apiserv.handlers {
		handler.AddRoutes(apiserv.app)
	}
}

// listenAndServe starts the fiber app and listens for incoming requests. It also listens for shutdown signals and handles graceful shutdown.
func listenAndServe(apiserv *apiServer) error {
	shutdown := make(chan os.Signal, 1)
	apiErr := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		apiserv.logger.Info("starting api server", "listening on", apiserv.config.Port, "env", apiserv.serverEnv)
		apiErr <- apiserv.app.Listen(":" + apiserv.config.Port)
	}()

	select {
	case err := <-apiErr:
		return fmt.Errorf("error listening api server: %w", err)
	case <-shutdown:
		apiserv.logger.Info("starting shutdown", "pid", os.Getpid())
		ctx, cancel := context.WithTimeout(context.Background(), constant.ShutdownTimeout)
		defer cancel()
		if err := apiserv.app.ShutdownWithContext(ctx); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
		apiserv.logger.Info("shutdown complete", "pid", os.Getpid())
	}
	return nil
}
