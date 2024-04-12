package taskhandler_test

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	dtores "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/res"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/basehttphandler"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/taskhandler"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg"
	"net/http/httptest"
	"testing"
)

func Test_taskHandler_AddRoutes(t *testing.T) {
	mockTaskService := &mockTaskService{}
	mockUserService := &mockUserService{}
	mockJwtUtils := &mockJwtUtils{}
	mockValidator := &mockValidator{}
	mockPassUtils := &mockPassUtils{}
	mockResponse := &mockResponse{}
	mockMiddleware := &mockMiddleware{}
	pkgs := pkg.New(
		pkg.WithJwtUtils(mockJwtUtils),
		pkg.WithValidator(mockValidator),
		pkg.WithPassUtils(mockPassUtils),
		pkg.WithResponse(mockResponse),
		pkg.WithMiddleware(mockMiddleware),
	)
	basehttphandler := basehttphandler.New(
		basehttphandler.WithPackages(pkgs),
		basehttphandler.WithContextTimeout(10),
		basehttphandler.WithLogger(nil),
	)
	taskHandler := taskhandler.New(
		taskhandler.WithTaskService(mockTaskService),
		taskhandler.WithUserService(mockUserService),
		taskhandler.WithBaseHttpHandler(basehttphandler),
	)
	{
		tc := "Case 1: Look for the number of routes in the fiber app"
		app := fiber.New()
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			return c.Next()
		}
		taskHandler.AddRoutes(app)
		t.Run(tc, func(t *testing.T) {
			if len(app.Stack()) == 0 {
				t.Fatalf("expected %d, got %d", 3, len(app.Stack()))
			}
		})
	}
}

func Test_taskHandler_EnqueueTask(t *testing.T) {
	mockTaskService := &mockTaskService{}
	mockUserService := &mockUserService{}
	mockJwtUtils := &mockJwtUtils{}
	mockValidator := &mockValidator{}
	mockPassUtils := &mockPassUtils{}
	mockResponse := &mockResponse{}
	mockMiddleware := &mockMiddleware{}
	pkgs := pkg.New(
		pkg.WithJwtUtils(mockJwtUtils),
		pkg.WithValidator(mockValidator),
		pkg.WithPassUtils(mockPassUtils),
		pkg.WithResponse(mockResponse),
		pkg.WithMiddleware(mockMiddleware),
	)
	basehttphandler := basehttphandler.New(
		basehttphandler.WithPackages(pkgs),
		basehttphandler.WithContextTimeout(10),
		basehttphandler.WithLogger(nil),
	)
	taskHandler := taskhandler.New(
		taskhandler.WithTaskService(mockTaskService),
		taskhandler.WithUserService(mockUserService),
		taskhandler.WithBaseHttpHandler(basehttphandler),
	)
	{
		tc := "Case 1: Bearer token not found in request and returns 401"
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(mockResponse.BasicError("missing token", fiber.StatusUnauthorized))
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Post("/api/v1/task/enqueue", taskHandler.EnqueueTask)
		req := httptest.NewRequest("POST", "/api/v1/task/enqueue", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Fatalf("expected %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
			}
		})
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 2: Validation error in request and returns 400"
		mockValidator.errBindAndValidate = errors.New("validation error")
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Post("/api/v1/task/enqueue", taskHandler.EnqueueTask)
		req := httptest.NewRequest("POST", "/api/v1/task/enqueue", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusBadRequest {
				t.Fatalf("expected %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
			}
		})
		mockValidator.errBindAndValidate = nil
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 3: Task service returns error and returns 500"
		mockTaskService.errEnqueueMailTask = errors.New("task service error")
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Post("/api/v1/task/enqueue", taskHandler.EnqueueTask)
		req := httptest.NewRequest("POST", "/api/v1/task/enqueue", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusInternalServerError {
				t.Fatalf("expected %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
			}
		})
		mockTaskService.errEnqueueMailTask = nil
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 4: Success"
		mockTaskService.resEnqueueMailTask = dtores.TaskEnqueueResponse{
			TaskID: 1,
		}
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Post("/api/v1/task/enqueue", taskHandler.EnqueueTask)
		req := httptest.NewRequest("POST", "/api/v1/task/enqueue", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("expected %d, got %d", fiber.StatusOK, resp.StatusCode)
			}
		})
	}
}

func Test_taskHandler_GetAllQueuedTasks(t *testing.T) {
	mockTaskService := &mockTaskService{}
	mockUserService := &mockUserService{}
	mockJwtUtils := &mockJwtUtils{}
	mockValidator := &mockValidator{}
	mockPassUtils := &mockPassUtils{}
	mockResponse := &mockResponse{}
	mockMiddleware := &mockMiddleware{}
	pkgs := pkg.New(
		pkg.WithJwtUtils(mockJwtUtils),
		pkg.WithValidator(mockValidator),
		pkg.WithPassUtils(mockPassUtils),
		pkg.WithResponse(mockResponse),
		pkg.WithMiddleware(mockMiddleware),
	)
	basehttphandler := basehttphandler.New(
		basehttphandler.WithPackages(pkgs),
		basehttphandler.WithContextTimeout(10),
		basehttphandler.WithLogger(nil),
	)
	taskHandler := taskhandler.New(
		taskhandler.WithTaskService(mockTaskService),
		taskhandler.WithUserService(mockUserService),
		taskhandler.WithBaseHttpHandler(basehttphandler),
	)
	{
		tc := "Case 1: Bearer token not found in request and returns 401"
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(mockResponse.BasicError("missing token", fiber.StatusUnauthorized))
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/task/failed", taskHandler.GetAllFailedQueuedTasks)
		req := httptest.NewRequest("GET", "/api/v1/task/failed", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Fatalf("expected %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
			}
		})
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 2: Validation error in request and returns 400"
		mockValidator.errBindAndValidate = errors.New("validation error")
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/task/failed", taskHandler.GetAllFailedQueuedTasks)
		req := httptest.NewRequest("GET", "/api/v1/task/failed", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusBadRequest {
				t.Fatalf("expected %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
			}
		})
		mockValidator.errBindAndValidate = nil
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 3: Task service returns error and returns 500"
		mockTaskService.errGetAllFailedQueuedTasks = errors.New("task service error")
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/task/failed", taskHandler.GetAllFailedQueuedTasks)
		req := httptest.NewRequest("GET", "/api/v1/task/failed", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusInternalServerError {
				t.Fatalf("expected %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
			}
		})
		mockTaskService.errGetAllFailedQueuedTasks = nil
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 4: Success"
		mockTaskService.resGetAllFailedQueuedTasks = dtores.GetAllFailedTasksResponse{
			Tasks: []dtores.BaseTaskResponse{
				{
					TaskID: 1,
				},
			},
		}
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/task/failed", taskHandler.GetAllFailedQueuedTasks)
		req := httptest.NewRequest("GET", "/api/v1/task/failed", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("expected %d, got %d", fiber.StatusOK, resp.StatusCode)
			}
		})
	}
}

func Test_taskHandler_GetAllFailedQueuedTasks(t *testing.T) {
	mockTaskService := &mockTaskService{}
	mockUserService := &mockUserService{}
	mockJwtUtils := &mockJwtUtils{}
	mockValidator := &mockValidator{}
	mockPassUtils := &mockPassUtils{}
	mockResponse := &mockResponse{}
	mockMiddleware := &mockMiddleware{}
	pkgs := pkg.New(
		pkg.WithJwtUtils(mockJwtUtils),
		pkg.WithValidator(mockValidator),
		pkg.WithPassUtils(mockPassUtils),
		pkg.WithResponse(mockResponse),
		pkg.WithMiddleware(mockMiddleware),
	)
	basehttphandler := basehttphandler.New(
		basehttphandler.WithPackages(pkgs),
		basehttphandler.WithContextTimeout(10),
		basehttphandler.WithLogger(nil),
	)
	taskHandler := taskhandler.New(
		taskhandler.WithTaskService(mockTaskService),
		taskhandler.WithUserService(mockUserService),
		taskhandler.WithBaseHttpHandler(basehttphandler),
	)
	{
		tc := "Case 1: Bearer token not found in request and returns 401"
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(mockResponse.BasicError("missing token", fiber.StatusUnauthorized))
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/task/failed", taskHandler.GetAllFailedQueuedTasks)
		req := httptest.NewRequest("GET", "/api/v1/task/failed", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Fatalf("expected %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
			}
		})
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 2: Validation error in request and returns 400"
		mockValidator.errBindAndValidate = errors.New("validation error")
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Post("/api/v1/task/enqueue", taskHandler.EnqueueTask)
		req := httptest.NewRequest("POST", "/api/v1/task/enqueue", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusBadRequest {
				t.Fatalf("expected %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
			}
		})
		mockValidator.errBindAndValidate = nil
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 3: Task service returns error and returns 500"
		mockTaskService.errGetAllFailedQueuedTasks = errors.New("task service error")
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/task/failed", taskHandler.GetAllFailedQueuedTasks)
		req := httptest.NewRequest("GET", "/api/v1/task/failed", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusInternalServerError {
				t.Fatalf("expected %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
			}
		})
		mockTaskService.errGetAllFailedQueuedTasks = nil
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 4: Success"
		mockTaskService.resGetAllFailedQueuedTasks = dtores.GetAllFailedTasksResponse{
			Tasks: []dtores.BaseTaskResponse{
				{
					TaskID:         1,
					RecipientEmail: "test@test.com",
					Status:         4,
					TryCount:       4,
				},
			},
		}
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/task/failed", taskHandler.GetAllFailedQueuedTasks)
		req := httptest.NewRequest("GET", "/api/v1/task/failed", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("expected %d, got %d", fiber.StatusOK, resp.StatusCode)
			}
		})
	}
}
