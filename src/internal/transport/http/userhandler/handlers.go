package userhandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/internal/dto"
)

func (h *userHandler) AddRoutes(r fiber.Router) {
	r.Post("/register", h.Register)
	user := r.Group("/user")
	user.Use(h.Packages.JwtUtils.AuthMiddleware())
	user.Get("/:id/details", h.GetUser)
	user.Put("/:id/update", h.UpdateUser)
}

func (h *userHandler) Register(c *fiber.Ctx) error {
	var (
		req dto.RegisterUserRequest
	)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	user, err := h.userService.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusCreated, user))
}

func (h *userHandler) GetUser(c *fiber.Ctx) error {
	var (
		req dto.GetUserRequest
	)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	user, err := h.userService.GetUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusOK, user))
}

func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	var (
		req dto.UpdateUserRequest
	)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	user, err := h.userService.UpdateUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusOK, user))
}
