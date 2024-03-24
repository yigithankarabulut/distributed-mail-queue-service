package workerservice

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/constant"
)

func (c *worker) TriggerWorker() {
	go func() {
		c.PlayWorker()
	}()
}

func (c *worker) PlayWorker() {
	for task := range c.ch {
		c.SendMail(task)
	}
}

func (c *worker) SendMail(task model.MailTaskQueue) {
	mailService := mailservice.New(
		mailservice.WithTask(task),
	)
	log.Infof("Worker %d sending mail to %s", c.id, task.RecipientEmail)
	err := mailService.SendMail(mailService.NewDialer(), mailService.NewMessage())
	if err != nil {
		log.Errorf("Error sending mail: %v", err)
		task.TryCount++
		if task.TryCount >= 3 {
			task.Status = constant.StatusFailed
			c.db.Save(&task)
			return
		}
		if err := c.taskqueue.PublishTask(constant.RedisMailQueueChannel, task); err != nil {
			log.Errorf("Error publishing task: %v", err)
		}
		task.Status = constant.StatusQueued
		c.db.Save(&task)
		return
	}
	task.Status = constant.StatusSuccess
	c.db.Save(&task)
	log.Infof("Worker %d sent mail to %s", c.id, task.RecipientEmail)
}
