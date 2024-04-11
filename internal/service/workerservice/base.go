package workerservice

import (
	"context"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskqueue"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskstorage"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
)

type IWorker interface {
	TriggerWorker()
	HandleTask(ctx context.Context, task model.MailTaskQueue) error
}

type worker struct {
	id          uint32
	mailService mailservice.MailService
	taskStorage taskstorage.TaskStorer
	taskqueue   taskqueue.TaskQueue
	taskChannel chan model.MailTaskQueue
}

type Option func(*worker)

func WithID(id int) Option {
	return func(w *worker) {
		w.id = uint32(id)
	}
}

func WithTaskStorage(rds taskstorage.TaskStorer) Option {
	return func(w *worker) {
		w.taskStorage = rds
	}
}

func WithTaskQueue(rds taskqueue.TaskQueue) Option {
	return func(w *worker) {
		w.taskqueue = rds
	}
}

func WithChannel(ch chan model.MailTaskQueue) Option {
	return func(w *worker) {
		w.taskChannel = ch
	}
}

func New(opts ...Option) IWorker {
	w := &worker{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}
