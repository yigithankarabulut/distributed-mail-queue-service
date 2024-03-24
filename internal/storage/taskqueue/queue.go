package taskqueue

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"time"
)

func (r *taskQueue) PublishTask(channel string, task interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	if err := r.rdb.LPush(ctx, channel, taskJson).Err(); err != nil {
		return err
	}
	log.Infof("Publishing task to channel: %s", channel)
	return nil
}

func (r *taskQueue) SubscribeTask(channel string) error {
	var (
		ctx  context.Context
		task model.MailTaskQueue
	)

	ctx = context.Background()
	log.Infof("Subscribing to channel: %s", channel)
	for {
		msg, err := r.rdb.LPop(ctx, channel).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			return err
		}
		log.Infof("Received message to channel: %s", channel)
		if err := json.Unmarshal([]byte(msg), &task); err != nil {
			log.Errorf("Error unmarshalling task: %v", err)
			continue
		}
		r.ch <- task
		log.Infof("Task sent to internal channel.")
	}
}
