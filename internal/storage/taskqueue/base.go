package taskqueue

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
)

type TaskQueue interface {
	PublishTask(task interface{}) error
	SubscribeTask(ctx context.Context, consumerID int) error
	StartConsume(ctx context.Context) <-chan error
}

type taskQueue struct {
	consumerCount int
	queueName     string
	rdb           *redis.Client
	taskChannel   chan model.MailTaskQueue
}

type Option func(*taskQueue)

func WithConsumerCount(count int) Option {
	return func(r *taskQueue) {
		r.consumerCount = count
	}
}

func WithQueueName(name string) Option {
	return func(r *taskQueue) {
		r.queueName = name
	}
}

func WithRedisClient(rdb *redis.Client) Option {
	return func(r *taskQueue) {
		r.rdb = rdb
	}
}

func WithTaskChannel(ch chan model.MailTaskQueue) Option {
	return func(r *taskQueue) {
		r.taskChannel = ch
	}
}

func New(opts ...Option) TaskQueue {
	queue := &taskQueue{}
	for _, opt := range opts {
		opt(queue)
	}
	return queue
}
