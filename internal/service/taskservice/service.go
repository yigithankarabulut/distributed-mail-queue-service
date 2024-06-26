package taskservice

import (
	"context"
	dtoreq "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	dtores "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/res"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/constant"
	"log"
)

func (s *taskService) EnqueueMailTask(ctx context.Context, request dtoreq.TaskEnqueueRequest) (dtores.TaskEnqueueResponse, error) {
	var (
		task model.MailTaskQueue
		res  dtores.TaskEnqueueResponse
	)
	select {
	case <-ctx.Done():
		return dtores.TaskEnqueueResponse{}, ctx.Err()
	default:
		task = request.ConvertToMailTaskQueue()
		task, err := s.taskStorage.Insert(ctx, task)
		if err != nil {
			return dtores.TaskEnqueueResponse{}, err
		}
		task.User, err = s.userStorage.GetByID(ctx, task.UserID)
		if err != nil {
			return dtores.TaskEnqueueResponse{}, err
		}
		if err := s.redisClient.PublishTask(ctx, task); err != nil {
			return dtores.TaskEnqueueResponse{}, err
		}
		res.TaskID = task.ID
		return res, nil
	}
}

func (s *taskService) GetAllQueuedTasks(ctx context.Context, request dtoreq.GetAllQueuedTasksRequest) (dtores.GetAllQueuedTasksResponse, error) {
	var (
		res dtores.GetAllQueuedTasksResponse
	)
	select {
	case <-ctx.Done():
		return dtores.GetAllQueuedTasksResponse{}, ctx.Err()
	default:
		tasks, err := s.taskStorage.GetAll(ctx, request.UserID)
		if err != nil {
			return dtores.GetAllQueuedTasksResponse{}, err
		}
		res.ToMailTaskQueue(tasks)
		return res, nil
	}
}

func (s *taskService) GetAllFailedQueuedTasks(ctx context.Context, request dtoreq.GetAllFailedTasksRequest) (dtores.GetAllFailedTasksResponse, error) {
	var (
		res dtores.GetAllFailedTasksResponse
	)
	select {
	case <-ctx.Done():
		return dtores.GetAllFailedTasksResponse{}, ctx.Err()
	default:
		tasks, err := s.taskStorage.GetAllByStatusWithUserID(ctx, constant.StatusCancelled, request.UserID)
		if err != nil {
			return dtores.GetAllFailedTasksResponse{}, err
		}
		res.ToMailTaskQueue(tasks)
		return res, nil
	}
}

func (s *taskService) FindUnprocessedTasksAndEnqueue() {
	var (
		tasks []model.MailTaskQueue
		err   error
	)
	log.Println("Finding unprocessed tasks and enqueueing...")
	ctx, cancel := context.WithTimeout(context.Background(), constant.TaskCancelTimeout)
	defer cancel()
	tasks, err = s.taskStorage.GetAllByUnprocessedTasks(ctx)
	if err != nil {
		log.Printf("error finding unprocessed tasks: %v", err)
		return
	}
	for _, task := range tasks {
		if err := s.redisClient.PublishTask(ctx, task); err != nil {
			log.Printf("error publishing task: %v", err)
		}
	}
	log.Printf("%d unprocessed tasks enqueued", len(tasks))
}
