package workerservice

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/constant"
)

func (c *worker) TriggerWorker() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case <-c.done:
			return fmt.Errorf("worker %d done", c.id)
		case task, ok := <-c.taskChannel:
			if !ok {
				return fmt.Errorf("worker %d task channel closed", c.id)
			}
			if err := c.HandleTask(ctx, task); err != nil {
				log.Errorf("worker %d error handling task: %v", c.id, err)
			}
		}
	}
}

func (c *worker) HandleTask(ctx context.Context, task model.MailTaskQueue) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := c.mailService.AddTask(task); err != nil {
			return fmt.Errorf("worker %d error adding task: %v", c.id, err)
		}
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
	if err := c.taskqueue.PublishTask(ctx, task); err != nil {
		log.Errorf("worker %d error publishing task: %v", c.id, err)
	}
	return nil
}
