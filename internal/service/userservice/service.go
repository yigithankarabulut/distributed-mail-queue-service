package userservice

import (
	"context"
	"fmt"
	dtoreq "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	dtores "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/res"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"time"
)

func (s *userService) Register(ctx context.Context, req dtoreq.RegisterRequest) error {
	var (
		user model.User
		tx   = s.userStorage.CreateTx()
	)
	defer s.userStorage.RollbackTx(tx)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if _, err := s.userStorage.GetByEmail(ctx, req.Email); err == nil {
			return fmt.Errorf("email already exists")
		}
		hashPwd, err := s.Packages.PassUtils.HashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}
		req.Password = hashPwd
		if err = s.userStorage.Insert(ctx, user); err != nil {
			return fmt.Errorf("error inserting user: %w", err)
		}
		s.userStorage.CommitTx(tx)
		return nil
	}
}

func (s *userService) Login(ctx context.Context, req dtoreq.LoginRequest) (dtores.LoginResponse, error) {
	var (
		res dtores.LoginResponse
	)
	select {
	case <-ctx.Done():
		return res, ctx.Err()
	default:
		user, err := s.userStorage.GetByEmail(ctx, req.Email)
		if err != nil {
			return res, fmt.Errorf("error getting user: %w", err)
		}
		if err := s.PassUtils.ComparePassword(user.Password, req.Password); err != nil {
			return res, err
		}

		token, err := s.JwtUtils.GenerateJwtToken(user.ID, 12*time.Hour)
		if err != nil {
			return res, fmt.Errorf("error generating token: %w", err)
		}
		res.ID = user.ID
		res.Token = token
		return res, nil
	}
}

func (s *userService) GetUser(ctx context.Context, req dtoreq.GetUserRequest) (dtores.GetUserResponse, error) {
	var (
		res dtores.GetUserResponse
	)
	select {
	case <-ctx.Done():
		return res, ctx.Err()
	default:
		user, err := s.userStorage.GetByID(ctx, req.UserID)
		if err != nil {
			return res, fmt.Errorf("error getting user: %w", err)
		}
		res.FromUser(user)
		return res, nil
	}
}
