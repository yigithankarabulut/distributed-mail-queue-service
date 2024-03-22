package taskhandler

import "github.com/gofiber/fiber/v2"

func (h *taskHandler) AddRoutes(r fiber.Router) {
	task := r.Group("/task")
	task.Post("/enqueue", h.EnqueueTask)
	task.Get("/queue", h.GetAllQueuedTasks)
	task.Get("/queue/fail", h.GetAllFailedQueuedTasks)
}

func (h *taskHandler) EnqueueTask(c *fiber.Ctx) error {
	return nil
}

func (h *taskHandler) GetAllQueuedTasks(c *fiber.Ctx) error {
	return nil
}

func (h *taskHandler) GetAllFailedQueuedTasks(c *fiber.Ctx) error {
	return nil
}
