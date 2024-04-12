package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type IMiddleware interface {
	AuthMiddleware() fiber.Handler
}

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}
