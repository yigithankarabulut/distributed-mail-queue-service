package userhandler_test

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	dtores "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/res"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/basehttphandler"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/transport/http/userhandler"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg"
	"net/http/httptest"
	"testing"
)

func Test_userHandler_AddRoutes(t *testing.T) {
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
	userHandler := userhandler.New(
		userhandler.WithBaseHttpHandler(basehttphandler),
		userhandler.WithUserService(mockUserService),
		userhandler.WithTaskService(mockTaskService),
	)
	{
		tc := "Case 1: Look for the number of routes in the fiber app"
		app := fiber.New()
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			return c.Next()
		}
		userHandler.AddRoutes(app)
		t.Run(tc, func(t *testing.T) {
			if len(app.Stack()) == 0 {
				t.Fatalf("expected %d, got %d", 3, len(app.Stack()))
			}
		})
	}
}

func Test_userHandler_Register(t *testing.T) {
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
	userHandler := userhandler.New(
		userhandler.WithBaseHttpHandler(basehttphandler),
		userhandler.WithUserService(mockUserService),
		userhandler.WithTaskService(mockTaskService),
	)
	{
		tc := "Case 1: Validation error in request and returns 400"
		mockValidator.errBindAndValidate = errors.New("validation error")
		app := fiber.New()
		app.Post("/api/v1/register", userHandler.Register)
		req := httptest.NewRequest("POST", "/api/v1/register", nil)
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
	}
	{
		tc := "Case 2: Error in register service and returns 500"
		mockUserService.errRegister = errors.New("error in register service")
		app := fiber.New()
		app.Post("/api/v1/register", userHandler.Register)
		req := httptest.NewRequest("POST", "/api/v1/register", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusInternalServerError {
				t.Fatalf("expected %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
			}
		})
		mockUserService.errRegister = nil
	}
	{
		tc := "Case 3: Successful registration"
		app := fiber.New()
		app.Post("/api/v1/register", userHandler.Register)
		req := httptest.NewRequest("POST", "/api/v1/register", nil)
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

func Test_userHandler_Login(t *testing.T) {
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
	userHandler := userhandler.New(
		userhandler.WithBaseHttpHandler(basehttphandler),
		userhandler.WithUserService(mockUserService),
		userhandler.WithTaskService(mockTaskService),
	)
	{
		{
			tc := "Case 1: Validation error in request and returns 400"
			mockValidator.errBindAndValidate = errors.New("validation error")

			app := fiber.New()
			app.Post("/api/v1/login", userHandler.Login)
			req := httptest.NewRequest("POST", "/api/v1/login", nil)
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
		}
		{
			tc := "Case 2: Error in login service and returns 500"
			mockUserService.errLogin = errors.New("error in login service")
			app := fiber.New()
			app.Post("/api/v1/login", userHandler.Login)
			req := httptest.NewRequest("POST", "/api/v1/login", nil)
			t.Run(tc, func(t *testing.T) {
				resp, err := app.Test(req)
				if err != nil {
					t.Fatalf("expected nil, got %v", err)
				}
				if resp.StatusCode != fiber.StatusInternalServerError {
					t.Fatalf("expected %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
				}
			})
			mockUserService.errLogin = nil
		}
		{
			tc := "Case 3: Successful login"
			mockUserService.resLogin = dtores.LoginResponse{
				Token: "token",
				ID:    1,
			}
			app := fiber.New()
			app.Post("/api/v1/login", userHandler.Login)
			req := httptest.NewRequest("POST", "/api/v1/login", nil)
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
}

func Test_userHandler_GetUser(t *testing.T) {
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
	userHandler := userhandler.New(
		userhandler.WithBaseHttpHandler(basehttphandler),
		userhandler.WithUserService(mockUserService),
		userhandler.WithTaskService(mockTaskService),
	)
	{
		tc := "Case 1: Bearer token not found in the request and returns 401"
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).SendString("bearer token not found")
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/user/:id", userHandler.GetUser)
		req := httptest.NewRequest("GET", "/api/v1/user/1", nil)
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
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		mockValidator.errBindAndValidate = errors.New("validation error")
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/user/:id", userHandler.GetUser)
		req := httptest.NewRequest("GET", "/api/v1/user/1", nil)
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
		tc := "Case 3: Error in get user service and returns 500"
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		mockUserService.errGetUser = errors.New("error in get user service")
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/user/:id", userHandler.GetUser)
		req := httptest.NewRequest("GET", "/api/v1/user/1", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusInternalServerError {
				t.Fatalf("expected %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
			}
		})
		mockUserService.errGetUser = nil
		mockMiddleware.errAuthMiddleware = nil
	}
	{
		tc := "Case 4: Successful get user"
		mockMiddleware.errAuthMiddleware = func(c *fiber.Ctx) error {
			c.Locals("userID", uint(1))
			return c.Next()
		}
		mockUserService.resGetUser = dtores.GetUserResponse{
			SmtpUsername: "test",
		}
		app := fiber.New()
		app.Use(mockMiddleware.AuthMiddleware())
		app.Get("/api/v1/user/:id", userHandler.GetUser)
		req := httptest.NewRequest("GET", "/api/v1/user/1", nil)
		t.Run(tc, func(t *testing.T) {
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("expected %d, got %d", fiber.StatusOK, resp.StatusCode)
			}
		})
		mockUserService.resGetUser = dtores.GetUserResponse{}
		mockMiddleware.errAuthMiddleware = nil
	}
}
