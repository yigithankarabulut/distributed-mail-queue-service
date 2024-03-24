package taskhandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	"github.com/yigithankarabulut/distributed-mail-queue-service/releaseinfo"
)

func (h *taskHandler) AddRoutes(r fiber.Router) {
	r.Use(h.Packages.JwtUtils.AuthMiddleware())
	r.Post(releaseinfo.EnqueueMailApiPath, h.EnqueueTask)
	r.Get(releaseinfo.GetAllQueuedMailTasksApiPath, h.GetAllQueuedTasks)
	r.Get(releaseinfo.GetAllFailedQueuedMailApiPath, h.GetAllFailedQueuedTasks)
}

func (h *taskHandler) EnqueueTask(c *fiber.Ctx) error {
	var (
		req dtoreq.TaskEnqueueRequest
	)
	req.UserID = c.Locals("userID").(uint)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	res, err := h.taskService.EnqueueMailTask(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusOK, res))
}

func (h *taskHandler) GetAllQueuedTasks(c *fiber.Ctx) error {
	var (
		req dtoreq.GetAllQueuedTasksRequest
	)
	req.UserID = c.Locals("userID").(uint)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	res, err := h.taskService.GetAllQueuedTasks(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusOK, res))
}

func (h *taskHandler) GetAllFailedQueuedTasks(c *fiber.Ctx) error {
	var (
		req dtoreq.GetAllFailedTasksRequest
	)
	req.UserID = c.Locals("userID").(uint)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	res, err := h.taskService.GetAllFailedQueuedTasks(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusOK, res))
}
