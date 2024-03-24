package userhandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	"github.com/yigithankarabulut/distributed-mail-queue-service/releaseinfo"
)

func (h *userHandler) AddRoutes(r fiber.Router) {
	r.Post(releaseinfo.RegisterUserApiPath, h.Register)
	r.Post(releaseinfo.LoginUserApiPath, h.Login)
	r.Use(h.Packages.JwtUtils.AuthMiddleware())
	r.Get(releaseinfo.GetUserApiPath, h.GetUser)
	r.Put(releaseinfo.UpdateUserApiPath, h.UpdateUser)
}

func (h *userHandler) Register(c *fiber.Ctx) error {
	var (
		req dtoreq.RegisterRequest
	)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	if err := h.userService.Register(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusCreated, "user registered successfully"))
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	var (
		req dtoreq.LoginRequest
	)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	res, err := h.userService.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusOK, res))
}

func (h *userHandler) GetUser(c *fiber.Ctx) error {
	var (
		req dtoreq.GetUserRequest
	)
	req.UserID = c.Locals("userID").(uint)
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
		req dtoreq.UpdateUserRequest
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
