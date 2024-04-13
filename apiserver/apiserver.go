package apiserver

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
	"time"
)

// initializeStorages initializes the storages with the given database connection.
func (s *apiServer) initializeStorages() {
	s.instances.userStorage = userstorage.New(userstorage.WithUserDB(postgres.DB))
	s.instances.taskStorage = taskstorage.New(taskstorage.WithTaskDB(postgres.DB))
	s.instances.taskQueue = taskqueue.New(
		taskqueue.WithTaskChannel(s.taskChannel),
		taskqueue.WithConsumerCount(constant.QueueConsumerCount),
		taskqueue.WithQueueName(constant.RedisMailQueueChannel),
		taskqueue.WithRedisClient(redisclient.GetRedisClient()),
	)
}

// initializeServices initializes the services with the given storages and packages.
func (s *apiServer) initializeServices() {
	s.instances.cronService = cron.NewCronService()
	s.instances.cronService.Start()
	s.instances.userService = userservice.New(
		userservice.WithUserStorage(s.instances.userStorage),
		userservice.WithTaskStorage(s.instances.taskStorage),
		userservice.WithPackages(s.instances.packages),
		userservice.WithMailService(mailservice.New()),
	)
	s.instances.taskService = taskservice.New(
		taskservice.WithTaskStorage(s.instances.taskStorage),
		taskservice.WithUserStorage(s.instances.userStorage),
		taskservice.WithRedisClient(s.instances.taskQueue),
	)
	handleUnprocessedJob := cron.CronJob{
		Name:     "FindUnprocessedTasksAndEnqueue",
		Schedule: "@every 5m",
		Func:     s.instances.taskService.FindUnprocessedTasksAndEnqueue,
	}
	if err := s.instances.cronService.RegisterJob(handleUnprocessedJob); err != nil {
		s.logger.Error("error registering cron job", "error", err)
	}
}

// initializeWorkers initializes the queue consumers and workers. It also triggers the workers.
func (s *apiServer) initializeWorkers() {
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errCh := s.instances.taskQueue.StartConsume(ctx)
		for {
			select {
			case <-s.done:
				log.Info("task queue consumers done")
				return
			case err := <-errCh:
				log.Errorf("task queue consumer error: %v", err)
			}
		}
	}()
	for i := 0; i < constant.WorkerCount; i++ {
		s.instances.workers[i] = workerservice.New(
			workerservice.WithID(i+1),
			workerservice.WithTaskStorage(s.instances.taskStorage),
			workerservice.WithTaskQueue(s.instances.taskQueue),
			workerservice.WithChannel(s.taskChannel),
			workerservice.WithDoneChannel(s.done),
			workerservice.WithMailService(mailservice.New()),
		)
	}
	for _, worker := range s.instances.workers {
		go func(w workerservice.IWorker) {
			if err := w.TriggerWorker(); err != nil {
				log.Errorf("error triggering worker: %v", err)
			}
		}(worker)
	}
}

// initializeHandlers initializes the handlers with the given base http handler and adds the routes to the fiber app.
func (s *apiServer) initializeHandlers() {
	baseHttpHandler := basehttphandler.New(
		basehttphandler.WithContextTimeout(constant.ContextCancelTimeout),
		basehttphandler.WithLogger(s.logger),
		basehttphandler.WithPackages(s.instances.packages),
	)
	userHandler := userhandler.New(
		userhandler.WithBaseHttpHandler(baseHttpHandler),
		userhandler.WithUserService(s.instances.userService),
		userhandler.WithTaskService(s.instances.taskService),
	)
	taskHandler := taskhandler.New(
		taskhandler.WithBaseHttpHandler(baseHttpHandler),
		taskhandler.WithTaskService(s.instances.taskService),
		taskhandler.WithUserService(s.instances.userService),
	)
	s.handlers = append(s.handlers, userHandler, taskHandler)
	for _, handler := range s.handlers {
		handler.AddRoutes(s.app)
	}
}

// createInstance creates a new instance of the server dependencies.
func (s *apiServer) createInstance() {
	s.instances = new(Instances)
	s.done = make(chan struct{})
	s.taskChannel = make(chan model.MailTaskQueue, constant.WorkerCount)
	s.instances.workers = make([]workerservice.IWorker, constant.WorkerCount)
	s.instances.packages = pkg.New(
		pkg.WithValidator(validator.New()),
		pkg.WithJwtUtils(jwtutils.New()),
		pkg.WithPassUtils(passutils.New()),
		pkg.WithResponse(response.New()),
		pkg.WithMiddleware(middleware.New()),
	)
	s.initializeStorages()
	s.initializeServices()
	s.initializeWorkers()
	s.initializeHandlers()
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
	corsConfig.AllowHeaders = strings.Join([]string{
		constant.ContentType, constant.Authorization}, ",")
	corsConfig.AllowMethods = strings.Join([]string{
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodPut,
		fiber.MethodDelete,
		fiber.MethodPatch,
		fiber.MethodOptions}, ",")
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
	s.createInstance()
	return s.listenAndServe()
}

// listenAndServe starts the fiber app and listens for incoming requests. It also listens for shutdown signals and handles graceful shutdown.
func (s *apiServer) listenAndServe() error {
	shutdown := make(chan os.Signal, 2)
	apiErr := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.logger.Info("starting api server", "listening on", s.config.Port, "env", s.serverEnv)
		if err := s.app.Listen(":" + s.config.Port); err != nil {
			apiErr <- err
		}
	}()
	closeChan := func() {
		close(s.done)
		close(shutdown)
		close(apiErr)
		close(s.taskChannel)
	}

	select {
	case err := <-apiErr:
		closeChan()
		return fmt.Errorf("error listening api server: %w", err)
	case <-shutdown:
		s.logger.Info("starting shutdown", "pid", os.Getpid())
		s.done <- struct{}{}
		ctx, cancel := context.WithTimeout(context.Background(), constant.ShutdownTimeout)
		defer cancel()
		if err := s.app.ShutdownWithContext(ctx); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
		closeChan()
		time.Sleep(constant.ShutdownTimeout)
		s.logger.Info("shutdown complete", "pid", os.Getpid())
	}
	return nil
}
