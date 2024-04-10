package workerservice

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/constant"
)

func (c *worker) TriggerWorker() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for task := range c.taskChannel {
		go func() {
			if err := c.SendMail(ctx, task); err != nil {
				log.Errorf("worker %d error sending mail: %v", c.id, err)
			}
		}()
	}
}

func (c *worker) SendMail(ctx context.Context, task model.MailTaskQueue) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.mailService = mailservice.New(
			mailservice.WithTask(task),
		)
		log.Infof("worker %d sending mail to %s", c.id, task.RecipientEmail)
		err := c.mailService.SendMail(c.mailService.NewDialer(), c.mailService.NewMessage())
		if err != nil {
			return c.handleError(ctx, task, err)
		}
		task.Status = constant.StatusSuccess
		if err := c.taskStorage.Update(ctx, task); err != nil {
			log.Errorf("worker %d error updating task: %v", c.id, err)
		}
		log.Infof("worker %d sent mail to %s", c.id, task.RecipientEmail)
	}
	return nil
}

func (c *worker) handleError(ctx context.Context, task model.MailTaskQueue, err error) error {
	log.Errorf("worker %d error sending mail to %s: %v", c.id, task.RecipientEmail, err)
	task.TryCount++
	if task.TryCount >= constant.MaxTryCount {
		task.Status = constant.StatusCancelled
		if err := c.taskStorage.Update(ctx, task); err != nil {
			log.Errorf("worker %d error updating task: %v", c.id, err)
		}
		return fmt.Errorf("task %d cancelled after %d tries", task.ID, task.TryCount)
	}
	task.Status = constant.StatusFailed
	if err := c.taskStorage.Update(ctx, task); err != nil {
		log.Errorf("worker %d error updating task: %v", c.id, err)
	}
	if err := c.taskqueue.PublishTask(task); err != nil {
		log.Errorf("worker %d error publishing task: %v", c.id, err)
	}
	return nil
}
