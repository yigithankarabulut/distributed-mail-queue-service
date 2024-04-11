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

func (r *taskQueue) PublishTask(task interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	if err := r.rdb.LPush(ctx, r.queueName, taskJson).Err(); err != nil {
		return err
	}
	log.Infof("publishing task to channel: %s", r.queueName)
	return nil
}

func (r *taskQueue) SubscribeTask(ctx context.Context, consumerID int) error {
	var (
		task model.MailTaskQueue
	)
	log.Infof("consumer %d subscribed to channel: %s", consumerID, r.queueName)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			result, err := r.rdb.BRPop(ctx, 0, r.queueName).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				return err
			}
			if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
				log.Errorf("Consumer %d error unmarshalling task: %v", consumerID, err)
				continue
			}
			log.Infof("consumer %d received task id: %d", consumerID, task.ID)
			r.taskChannel <- task
			log.Infof("consumer %d sent task to internal channel", consumerID)
		}
	}

}

func (r *taskQueue) StartConsume(ctx context.Context) <-chan error {
	errCh := make(chan error)
	for i := 0; i < r.consumerCount; i++ {
		go func(consumerID int) {
			if err := r.SubscribeTask(ctx, consumerID); err != nil {
				errCh <- err
			}
		}(i)
	}
	return errCh
}
