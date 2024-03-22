package userservice

import (
	"context"
	"fmt"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/dto"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/model"
)

// TODO: add custom error types

func (s *userService) Register(ctx context.Context, req dto.RegisterUserRequest) error {
	var (
		user model.User
		err  error
	)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if _, err = s.userStorage.GetByEmail(ctx, user.Email); err == nil {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
		hashPwd, err := s.PassUtils.HashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}
		req.Password = hashPwd
		user = req.ConvertToUser()
		if err = s.userStorage.Insert(ctx, user); err != nil {
			return fmt.Errorf("error inserting user: %w", err)
		}
		return nil
	}
}

func (s *userService) GetUser(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	var (
		res dto.GetUserResponse
	)
	select {
	case <-ctx.Done():
		return res, ctx.Err()
	default:
		user, err := s.userStorage.GetByID(ctx, req.ID)
		if err != nil {
			return res, fmt.Errorf("error getting user: %w", err)
		}
		res.FromUser(user)
		return res, nil
	}
}

func (s *userService) UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (dto.UpdateUserResponse, error) {
	return dto.UpdateUserResponse{}, nil
}
