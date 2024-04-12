package apiserver

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
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
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/cron"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/jwtutils"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/middleware"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/passutils"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/postgres"
	redisclient "github.com/yigithankarabulut/distributed-mail-queue-service/pkg/redis"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/response"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/validator"
	"github.com/yigithankarabulut/distributed-mail-queue-service/releaseinfo"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	QueueConsumerCount = 10
	WorkerCount        = 10
)

// appendDepends create instances from all layers and append them together. Then add it to the fiber by calling the AddRoutes method of all its handlers.
func (s *apiServer) appendDepends() {
	var (
		taskQueueChannel chan model.MailTaskQueue
		cronService      *cron.CronService
		packages         *pkg.Packages
		taskQueue        taskqueue.TaskQueue
		userStorage      userstorage.UserStorer
		taskStorage      taskstorage.TaskStorer
		userService      userservice.UserService
		taskService      taskservice.TaskService
		workers          []workerservice.IWorker
		done             chan struct{}
	)

	taskQueueChannel = make(chan model.MailTaskQueue, WorkerCount)
	done = make(chan struct{}, WorkerCount)
	workers = make([]workerservice.IWorker, WorkerCount)
	cronService = cron.NewCronService()
	cronService.Start()

	packages = pkg.New(
		pkg.WithValidator(validator.New()),
		pkg.WithJwtUtils(jwtutils.New()),
		pkg.WithPassUtils(passutils.New()),
		pkg.WithResponse(response.New()),
		pkg.WithMiddleware(middleware.New()),
	)
	userStorage = userstorage.New(userstorage.WithUserDB(postgres.DB))
	taskStorage = taskstorage.New(taskstorage.WithTaskDB(postgres.DB))
	taskQueue = taskqueue.New(
		taskqueue.WithTaskChannel(taskQueueChannel),
		taskqueue.WithConsumerCount(QueueConsumerCount),
		taskqueue.WithQueueName(constant.RedisMailQueueChannel),
		taskqueue.WithRedisClient(redisclient.GetRedisClient()),
	)
	userService = userservice.New(
		userservice.WithUserStorage(userStorage),
		userservice.WithTaskStorage(taskStorage),
		userservice.WithPackages(packages),
		userservice.WithMailService(mailservice.New()),
	)
	taskService = taskservice.New(
		taskservice.WithTaskStorage(taskStorage),
		taskservice.WithUserStorage(userStorage),
		taskservice.WithRedisClient(taskQueue),
	)

	go func() {
		errCount := 0
		errCh := taskQueue.StartConsume(context.Background())
		select {
		case err := <-errCh:
			s.logger.Error("error consuming task", "error", err)
			errCount++
			if errCount == QueueConsumerCount {
				s.logger.Error("all queue consumers failed, shutting down")
			}
		}
	}()
	for i := 0; i < WorkerCount; i++ {
		workers[i] = workerservice.New(
			workerservice.WithID(i+1),
			workerservice.WithTaskStorage(taskStorage),
			workerservice.WithTaskQueue(taskQueue),
			workerservice.WithChannel(taskQueueChannel),
			workerservice.WithDoneChannel(done),
			workerservice.WithMailService(mailservice.New()),
		)
	}
	for _, worker := range workers {
		go func(w workerservice.IWorker) {
			if err := w.TriggerWorker(); err != nil {
				s.logger.Error("error triggering worker", "error", err)
			}
		}(worker)
	}

	handleUnprocessedJob := cron.CronJob{
		Name:     "FindUnprocessedTasksAndEnqueue",
		Schedule: "@every 5m",
		Func:     taskService.FindUnprocessedTasksAndEnqueue,
	}
	if err := cronService.RegisterJob(handleUnprocessedJob); err != nil {
		s.logger.Error("error registering cron job", "error", err)
	}

	// Create http handlers instances and add them to the fiber app.
	baseHttpHandler := basehttphandler.New(
		basehttphandler.WithContextTimeout(constant.ContextCancelTimeout),
		basehttphandler.WithLogger(s.logger),
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
	s.handlers = append(s.handlers, userHandler, taskHandler)
	for _, handler := range s.handlers {
		handler.AddRoutes(s.app)
	}
}

// healthzCheck adds health endpoints to the apiserver.
func (s *apiServer) healthzCheck() {
	s.app.Get("/healthz/live", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"server":            s.serverEnv,
			"version":           releaseinfo.Version,
			"build_information": releaseinfo.BuildInformation,
			"message":           "liveness is OK!, server is ready to accept connections",
		})
	})
	s.app.Get("/healthz/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "pong",
		})
	})
}

// initializeApp initializes the fiber app with the given logger and adds the http logging middleware.
func (s *apiServer) initializeApp() {
	var (
		corsConfig cors.Config
	)
	s.app = fiber.New(fiber.Config{
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
	s.app.Use(cors.New(corsConfig))
	s.app.Use(recover.New())
	s.app.Use(middleware.HttpLoggingMiddleware(s.logger, s.app))
}

// initializeServer initializes the server with the given config and logger.
func (s *apiServer) initializeServer() error {
	if s.logger == nil {
		logHandlerOpts := &slog.HandlerOptions{Level: s.logLevel}
		logHandler := slog.NewJSONHandler(os.Stdout, logHandlerOpts)
		s.logger = slog.New(logHandler)
	}
	slog.SetDefault(s.logger)
	if s.serverEnv == "" {
		s.serverEnv = "development"
	}
	return nil
}

// connectStorages connects to the storages and returns an error if any.
func (s *apiServer) connectStorages() error {
	if _, err := postgres.ConnectPQ(s.config.Database); err != nil {
		return fmt.Errorf("error connecting to postgres: %w", err)
	}
	if _, err := redisclient.New(s.config.Redis); err != nil {
		return fmt.Errorf("error connecting to redis: %w", err)
	}
	return nil
}

// Run starts the server and listens for incoming requests. It also initializes the server and connects to the storages.
func (s *apiServer) Run() error {
	if s.config == nil {
		return fmt.Errorf("config is required")
	}
	if err := s.connectStorages(); err != nil {
		return fmt.Errorf("error connecting to storages: %w", err)
	}
	if err := s.initializeServer(); err != nil {
		return fmt.Errorf("error initializing server: %w", err)
	}

	s.initializeApp()
	s.healthzCheck()
	s.appendDepends()
	return s.listenAndServe()
}

// listenAndServe starts the fiber app and listens for incoming requests. It also listens for shutdown signals and handles graceful shutdown.
func (s *apiServer) listenAndServe() error {
	shutdown := make(chan os.Signal, 1)
	apiErr := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.logger.Info("starting api server", "listening on", s.config.Port, "env", s.serverEnv)
		apiErr <- s.app.Listen(":" + s.config.Port)
	}()

	select {
	case err := <-apiErr:
		return fmt.Errorf("error listening api server: %w", err)
	case <-shutdown:
		s.logger.Info("starting shutdown", "pid", os.Getpid())
		ctx, cancel := context.WithTimeout(context.Background(), constant.ShutdownTimeout)
		defer cancel()
		if err := s.app.ShutdownWithContext(ctx); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
		s.logger.Info("shutdown complete", "pid", os.Getpid())
	}
	return nil
}
