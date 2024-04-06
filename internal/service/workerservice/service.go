package workerservice

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/constant"
	"time"
)

const (
	TaskTimeout = 5 * time.Second
)

func (c *worker) TriggerWorker() {
	ctx, cancel := context.WithTimeout(context.Background(), TaskTimeout)
	defer cancel()
	for task := range c.taskChannel {
		if err := c.SendMail(ctx, task); err != nil {
			log.Errorf("worker %d error sending mail: %v", c.id, err)
		}
	}
}

func (c *worker) SendMail(ctx context.Context, task model.MailTaskQueue) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		mailService := mailservice.New(
			mailservice.WithTask(task),
		)
		log.Infof("worker %d sending mail to %s", c.id, task.RecipientEmail)
		err := mailService.SendMail(mailService.NewDialer(), mailService.NewMessage())
		if err != nil {
			log.Errorf("worker %d error sending mail to %s: %v", c.id, task.RecipientEmail, err)
			task.TryCount++
			task.Status = constant.StatusFailed
			if task.TryCount < 3 {
				if err := c.taskqueue.PublishTask(task); err != nil {
					log.Errorf("worker %d error publishing task: %v", c.id, err)
				}
			}
			if err := c.taskStorage.Update(ctx, task); err != nil {
				log.Errorf("worker %d error updating task: %v", c.id, err)
			}
			return err
		}
		task.Status = constant.StatusSuccess
		if err := c.taskStorage.Update(ctx, task); err != nil {
			log.Errorf("worker %d error updating task: %v", c.id, err)
		}
		log.Infof("worker %d sent mail to %s", c.id, task.RecipientEmail)
	}
	return nil
}
