package taskqueue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"sync"
	"time"
)

func (r *taskQueue) PublishTask(ctx context.Context, task interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
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
}

func (r *taskQueue) SubscribeTask(ctx context.Context, consumerID int) error {
	var (
		task model.MailTaskQueue
	)
	log.Infof("consumer %d subscribed to channel: %s", consumerID, r.queueName)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("consumer %d done: %v", consumerID, ctx.Err())
		default:
			timeout := 1 * time.Second
			result, err := r.rdb.BRPop(ctx, timeout, r.queueName).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				return err
			}
			if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
				log.Errorf("consumer %d error unmarshalling task: %v", consumerID, err)
				continue
			}
			log.Infof("consumer %d received task id: %d", consumerID, task.ID)
			r.taskChannel <- task
			log.Infof("consumer %d sent task to internal channel", consumerID)
		}
	}
}

func (r *taskQueue) StartConsume(ctx context.Context) <-chan error {
	errCh := make(chan error, r.consumerCount)
	wg := sync.WaitGroup{}
	for i := 0; i < r.consumerCount; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			if err := r.SubscribeTask(ctx, consumerID+1); err != nil {
				log.Errorf("error consuming task: %v", err)
				errCh <- err
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()
	return errCh
}
