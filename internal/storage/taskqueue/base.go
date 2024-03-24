package taskqueue

import (
	"github.com/redis/go-redis/v9"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
)

type TaskQueue interface {
	PublishTask(channel string, task interface{}) error
	SubscribeTask(channel string) error
}

type taskQueue struct {
	rdb *redis.Client
	ch  chan model.MailTaskQueue
}

type Option func(*taskQueue)

func WithRedisClient(rdb *redis.Client) Option {
	return func(r *taskQueue) {
		r.rdb = rdb
	}
}

func WithChannel(ch chan model.MailTaskQueue) Option {
	return func(r *taskQueue) {
		r.ch = ch
	}
}

func New(opts ...Option) TaskQueue {
	queue := &taskQueue{}
	for _, opt := range opts {
		opt(queue)
	}
	return queue
}
