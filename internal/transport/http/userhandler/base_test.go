package userhandler_test

import (
	"context"
	"github.com/gofiber/fiber/v2"
	dtoreq "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	dtores "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/res"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/response"
	"log/slog"
	"time"
)

type mockUserService struct {
	errRegister error
	errLogin    error
	errGetUser  error
	resLogin    dtores.LoginResponse
	resGetUser  dtores.GetUserResponse
}

func (m *mockUserService) Register(ctx context.Context, req dtoreq.RegisterRequest) error {
	return m.errRegister
}

func (m *mockUserService) Login(ctx context.Context, req dtoreq.LoginRequest) (dtores.LoginResponse, error) {
	return m.resLogin, m.errLogin
}

func (m *mockUserService) GetUser(ctx context.Context, req dtoreq.GetUserRequest) (dtores.GetUserResponse, error) {
	return m.resGetUser, m.errGetUser
}

type mockTaskService struct {
	errEnqueueMailTask         error
	errGetAllQueuedTasks       error
	errGetAllFailedQueuedTasks error
	resEnqueueMailTask         dtores.TaskEnqueueResponse
	resGetAllQueuedTasks       dtores.GetAllQueuedTasksResponse
	resGetAllFailedQueuedTasks dtores.GetAllFailedTasksResponse
}

func (m *mockTaskService) EnqueueMailTask(ctx context.Context, request dtoreq.TaskEnqueueRequest) (dtores.TaskEnqueueResponse, error) {
	return m.resEnqueueMailTask, m.errEnqueueMailTask
}

func (m *mockTaskService) GetAllQueuedTasks(ctx context.Context, request dtoreq.GetAllQueuedTasksRequest) (dtores.GetAllQueuedTasksResponse, error) {
	return m.resGetAllQueuedTasks, m.errGetAllQueuedTasks
}

func (m *mockTaskService) GetAllFailedQueuedTasks(ctx context.Context, request dtoreq.GetAllFailedTasksRequest) (dtores.GetAllFailedTasksResponse, error) {
	return m.resGetAllFailedQueuedTasks, m.errGetAllFailedQueuedTasks
}

func (m *mockTaskService) FindUnprocessedTasksAndEnqueue() {
	return
}

type mockJwtUtils struct {
	errGenerateToken error
	resGenerateToken string
}

func (m *mockJwtUtils) GenerateJwtToken(userID uint, expiration time.Duration) (string, error) {
	return m.resGenerateToken, m.errGenerateToken
}

type mockPassUtils struct {
	errHashPassword error
	errCompareHash  error
	resHashPassword string
}

func (m *mockPassUtils) HashPassword(password string) (string, error) {
	return m.resHashPassword, m.errHashPassword
}

func (m *mockPassUtils) ComparePassword(hash, password string) error {
	return m.errCompareHash
}

type mockValidator struct {
	errBindAndValidate error
}

func (m *mockValidator) BindAndValidate(c *fiber.Ctx, data interface{}) error {
	return m.errBindAndValidate
}

type mockResponse struct {
	errBasicError response.ErrorResponse
	errData       response.DataResponse
}

func (m *mockResponse) BasicError(d interface{}, status int) response.ErrorResponse {
	return m.errBasicError
}

func (m *mockResponse) Data(status int, data interface{}) response.DataResponse {
	return m.errData
}

type mockMiddleware struct {
	errAuthMiddleware        fiber.Handler
	errHttpLoggingMiddleware fiber.Handler
}

func (m *mockMiddleware) AuthMiddleware() fiber.Handler {
	return m.errAuthMiddleware
}

func (m *mockMiddleware) HttpLoggingMiddleware(logger *slog.Logger, app *fiber.App) fiber.Handler {
	return m.errHttpLoggingMiddleware
}
