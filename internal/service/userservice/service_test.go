package userservice_test

import (
	"context"
	"errors"
	dtoreq "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/userservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg"
	"testing"
)

func Test_userService_Register(t *testing.T) {
	mockUserStorer := &mockUserStorer{}
	mockTaskStorer := &mockTaskStorer{}
	mockPassUtils := &mockPassUtils{}
	mockJwtUtils := &mockJwtUtils{}
	mockMailService := &mockMailService{}
	userService := userservice.New(
		userservice.WithUserStorage(mockUserStorer),
		userservice.WithTaskStorage(mockTaskStorer),
		userservice.WithMailService(mockMailService),
		userservice.WithPackages(pkg.New(
			pkg.WithPassUtils(mockPassUtils),
			pkg.WithJwtUtils(mockJwtUtils),
		)),
	)
	{
		tc := "Case 1: Context Cancelled And Should Return Error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := userService.Register(ctx, dtoreq.RegisterRequest{})
		want := "context canceled"
		t.Run(tc, func(t *testing.T) {
			if err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
	}
	{
		tc := "Case 2: Email Already Exists And Should Return Error"
		mockUserStorer.errGetByEmail = nil
		err := userService.Register(context.Background(), dtoreq.RegisterRequest{})

		want := "email already exists"
		t.Run(tc, func(t *testing.T) {
			if err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
	}
	{
		tc := "Case 3: Error Hashing Password And Should Return Error"
		mockUserStorer.errGetByEmail = errors.New("not found")
		mockPassUtils.errHashPassword = errors.New("hash error")

		err := userService.Register(context.Background(), dtoreq.RegisterRequest{})
		want := "error hashing password: hash error"
		t.Run(tc, func(t *testing.T) {
			if err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
		mockPassUtils.errHashPassword = nil
		mockUserStorer.errGetByEmail = nil
	}
	{
		tc := "Case 4: Error Adding Test Task And Should Return Error"
		mockUserStorer.errGetByEmail = errors.New("not found")
		mockMailService.errAddTask = errors.New("invalid task")

		err := userService.Register(context.Background(), dtoreq.RegisterRequest{})
		want := "error adding test task: invalid task"
		t.Run(tc, func(t *testing.T) {
			if err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
		mockMailService.errAddTask = nil
		mockUserStorer.errGetByEmail = nil
	}
	{
		tc := "Case 5: Error Sending Test Mail And Should Return Error"
		mockUserStorer.errGetByEmail = errors.New("not found")
		mockMailService.errSendMail = errors.New("invalid mail")

		err := userService.Register(context.Background(), dtoreq.RegisterRequest{})
		want := "error sending test mail: invalid mail"
		t.Run(tc, func(t *testing.T) {
			if err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
		mockMailService.errSendMail = nil
		mockUserStorer.errGetByEmail = nil
	}
	{
		tc := "Case 6: Error Inserting User And Should Return Error"
		mockUserStorer.errGetByEmail = errors.New("not found")
		mockUserStorer.errInsert = errors.New("insert error")

		err := userService.Register(context.Background(), dtoreq.RegisterRequest{})
		want := "error inserting user: insert error"
		t.Run(tc, func(t *testing.T) {
			if err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
		mockUserStorer.errInsert = nil
		mockUserStorer.errGetByEmail = nil
	}
	{
		tc := "Case 7: Success And Should Return Nil"
		mockUserStorer.errGetByEmail = errors.New("not found")
		mockUserStorer.errInsert = nil

		err := userService.Register(context.Background(), dtoreq.RegisterRequest{})
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("Expected error to be nil but got %v", err)
			}
		})
	}
}

func Test_userService_Login(t *testing.T) {
	mockUserStorer := &mockUserStorer{}
	mockTaskStorer := &mockTaskStorer{}
	mockPassUtils := &mockPassUtils{}
	mockJwtUtils := &mockJwtUtils{}
	userService := userservice.New(
		userservice.WithUserStorage(mockUserStorer),
		userservice.WithTaskStorage(mockTaskStorer),
		userservice.WithPackages(pkg.New(
			pkg.WithPassUtils(mockPassUtils),
			pkg.WithJwtUtils(mockJwtUtils),
		)),
	)
	{
		tc := "Case 1: Context Cancelled And Should Return Error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := userService.Login(ctx, dtoreq.LoginRequest{})
		want := "context canceled"
		t.Run(tc, func(t *testing.T) {
			if err != nil && err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
	}
	{
		tc := "Case 2: Error Getting User By Email And Should Return Error"
		mockUserStorer.errGetByEmail = errors.New("not found")

		_, err := userService.Login(context.Background(), dtoreq.LoginRequest{})
		want := "error getting user: not found"
		t.Run(tc, func(t *testing.T) {
			if err != nil && err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
		mockUserStorer.errGetByEmail = nil
	}
	{
		tc := "Case 3: Error Comparing Password And Should Return Error"
		mockPassUtils.errComparePassword = errors.New("invalid password")

		_, err := userService.Login(context.Background(), dtoreq.LoginRequest{})
		want := "invalid password"
		t.Run(tc, func(t *testing.T) {
			if err != nil && err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
		mockPassUtils.errComparePassword = nil
	}
	{
		tc := "Case 4: Generate Token Error And Should Return Error"
		mockJwtUtils.errGenerateJwtToken = errors.New("token error")

		_, err := userService.Login(context.Background(), dtoreq.LoginRequest{})
		want := "error generating token: token error"
		t.Run(tc, func(t *testing.T) {
			if err != nil && err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
		mockJwtUtils.errGenerateJwtToken = nil
	}
	{
		tc := "Case 5: Success And Should Return Nil"
		_, err := userService.Login(context.Background(), dtoreq.LoginRequest{})
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("Expected error to be nil but got %v", err)
			}
		})
	}
}

func Test_userService_GetUser(t *testing.T) {
	mockUserStorer := &mockUserStorer{}
	userService := userservice.New(
		userservice.WithUserStorage(mockUserStorer),
	)
	{
		tc := "Case 1: Context Cancelled And Should Return Error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := userService.GetUser(ctx, dtoreq.GetUserRequest{})
		want := "context canceled"
		t.Run(tc, func(t *testing.T) {
			if err != nil && err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
	}
	{
		tc := "Case 2: Error Getting User By ID And Should Return Error"
		mockUserStorer.errGetByID = errors.New("not found")

		_, err := userService.GetUser(context.Background(), dtoreq.GetUserRequest{})
		want := "error getting user: not found"
		t.Run(tc, func(t *testing.T) {
			if err != nil && err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
		mockUserStorer.errGetByID = nil
	}
	{
		tc := "Case 3: Success And Should Return Nil"
		_, err := userService.GetUser(context.Background(), dtoreq.GetUserRequest{})
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("Expected error to be nil but got %v", err)
			}
		})
	}
}
