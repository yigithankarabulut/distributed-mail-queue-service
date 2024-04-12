package taskservice_test

import (
	"bytes"
	"context"
	"errors"
	dtoreq "github.com/yigithankarabulut/distributed-mail-queue-service/internal/dto/req"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/taskservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"log"
	"regexp"
	"strings"
	"testing"
)

func Test_taskService_EnqueueMailTask(t *testing.T) {
	mockTaskStorer := &mockTaskStorer{}
	mockUserStorer := &mockUserStorer{}
	mockTaskQueue := &mockTaskQueue{}
	mockTaskService := taskservice.New(
		taskservice.WithTaskStorage(mockTaskStorer),
		taskservice.WithUserStorage(mockUserStorer),
		taskservice.WithRedisClient(mockTaskQueue),
	)
	{
		tc := "Case 1: Context is done and returns context error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := mockTaskService.EnqueueMailTask(ctx, dtoreq.TaskEnqueueRequest{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, context.Canceled) {
				t.Errorf("%s: expected %v but got %v", tc, context.Canceled, err)
			}
		})
	}
	{
		tc := "Case 2: TaskStorage Insert returns error"
		mockTaskStorer.errInsert = errors.New("insert error")
		_, err := mockTaskService.EnqueueMailTask(context.Background(), dtoreq.TaskEnqueueRequest{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, mockTaskStorer.errInsert) {
				t.Errorf("%s: expected %v but got %v", tc, mockTaskStorer.errInsert, err)
			}
		})
		mockTaskStorer.errInsert = nil
	}
	{
		tc := "Case 3: UserStorage GetByID returns error"
		mockUserStorer.errGetByID = errors.New("get by id error")
		_, err := mockTaskService.EnqueueMailTask(context.Background(), dtoreq.TaskEnqueueRequest{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, mockUserStorer.errGetByID) {
				t.Errorf("%s: expected %v but got %v", tc, mockUserStorer.errGetByID, err)
			}
		})
		mockUserStorer.errGetByID = nil
	}
	{
		tc := "Case 4: RedisClient PublishTask returns error"
		mockTaskQueue.errPublishTask = errors.New("publish task error")
		_, err := mockTaskService.EnqueueMailTask(context.Background(), dtoreq.TaskEnqueueRequest{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, mockTaskQueue.errPublishTask) {
				t.Errorf("%s: expected %v but got %v", tc, mockTaskQueue.errPublishTask, err)
			}
		})
		mockTaskQueue.errPublishTask = nil
	}
	{
		tc := "Case 5: Success"
		_, err := mockTaskService.EnqueueMailTask(context.Background(), dtoreq.TaskEnqueueRequest{})
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("%s: expected nil but got %v", tc, err)
			}
		})
	}
}

func Test_taskService_GetAllQueuedTasks(t *testing.T) {
	mockTaskStorer := &mockTaskStorer{}
	mockUserStorer := &mockUserStorer{}
	mockTaskQueue := &mockTaskQueue{}
	mockTaskService := taskservice.New(
		taskservice.WithTaskStorage(mockTaskStorer),
		taskservice.WithUserStorage(mockUserStorer),
		taskservice.WithRedisClient(mockTaskQueue),
	)
	{
		tc := "Case 1: Context is done and returns context error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := mockTaskService.GetAllQueuedTasks(ctx, dtoreq.GetAllQueuedTasksRequest{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, context.Canceled) {
				t.Errorf("%s: expected %v but got %v", tc, context.Canceled, err)
			}
		})
	}
	{
		tc := "Case 2: TaskStorage GetAll returns error"
		mockTaskStorer.errGetAll = errors.New("get all error")
		_, err := mockTaskService.GetAllQueuedTasks(context.Background(), dtoreq.GetAllQueuedTasksRequest{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, mockTaskStorer.errGetAll) {
				t.Errorf("%s: expected %v but got %v", tc, mockTaskStorer.errGetAll, err)
			}
		})
		mockTaskStorer.errGetAll = nil
	}
	{
		tc := "Case 3: Success"
		_, err := mockTaskService.GetAllQueuedTasks(context.Background(), dtoreq.GetAllQueuedTasksRequest{})
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("%s: expected nil but got %v", tc, err)
			}
		})
	}
}

func Test_taskService_GetAllFailedQueuedTasks(t *testing.T) {
	mockTaskStorer := &mockTaskStorer{}
	mockUserStorer := &mockUserStorer{}
	mockTaskQueue := &mockTaskQueue{}
	mockTaskService := taskservice.New(
		taskservice.WithTaskStorage(mockTaskStorer),
		taskservice.WithUserStorage(mockUserStorer),
		taskservice.WithRedisClient(mockTaskQueue),
	)
	{
		tc := "Case 1: Context is done and returns context error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := mockTaskService.GetAllFailedQueuedTasks(ctx, dtoreq.GetAllFailedTasksRequest{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, context.Canceled) {
				t.Errorf("%s: expected %v but got %v", tc, context.Canceled, err)
			}
		})
	}
	{
		tc := "Case 2: TaskStorage GetAllByStatusWithUserID returns error"
		mockTaskStorer.errGetAllByStatusWithUserID = errors.New("get all by status with user id error")
		_, err := mockTaskService.GetAllFailedQueuedTasks(context.Background(), dtoreq.GetAllFailedTasksRequest{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, mockTaskStorer.errGetAllByStatusWithUserID) {
				t.Errorf("%s: expected %v but got %v", tc, mockTaskStorer.errGetAllByStatusWithUserID, err)
			}
		})
		mockTaskStorer.errGetAllByStatusWithUserID = nil
	}
	{
		tc := "Case 3: Success"
		_, err := mockTaskService.GetAllFailedQueuedTasks(context.Background(), dtoreq.GetAllFailedTasksRequest{})
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("%s: expected nil but got %v", tc, err)
			}
		})
	}
}

func Test_taskService_FindUnprocessedTasksAndEnqueue(t *testing.T) {
	mockTaskStorer := &mockTaskStorer{}
	mockUserStorer := &mockUserStorer{}
	mockTaskQueue := &mockTaskQueue{}
	mockTaskService := taskservice.New(
		taskservice.WithTaskStorage(mockTaskStorer),
		taskservice.WithUserStorage(mockUserStorer),
		taskservice.WithRedisClient(mockTaskQueue),
	)
	{
		tc := "Case 1: TaskStorage GetAllByUnprocessedTasks returns error"
		mockTaskStorer.errGetAllByUnprocessedTasks = errors.New("get all by unprocessed tasks error")
		var buf bytes.Buffer
		log.SetOutput(&buf)
		mockTaskService.FindUnprocessedTasksAndEnqueue()

		logContents := buf.String()
		firstLog := " Finding unprocessed tasks and enqueueing...\n"
		expectedLog := " error finding unprocessed tasks: " + mockTaskStorer.errGetAllByUnprocessedTasks.Error() + "\n"
		fullLog := firstLog + expectedLog

		logContents = removeTimeInfo(logContents)
		t.Run(tc, func(t *testing.T) {
			if !strings.Contains(logContents, fullLog) {
				t.Errorf("%s: expected log:\n%s but got:\n%s", tc, fullLog, logContents)
			}
		})
		mockTaskStorer.errGetAllByUnprocessedTasks = nil
	}
	{
		tc := "Case 2: RedisClient PublishTask returns error"
		mockTaskQueue.errPublishTask = errors.New("throw error")
		mockTaskStorer.taskModelArr = []model.MailTaskQueue{{
			UserID: 1,
		}}
		var buf bytes.Buffer
		log.SetOutput(&buf)
		mockTaskService.FindUnprocessedTasksAndEnqueue()

		logContents := buf.String()
		firstLog := " Finding unprocessed tasks and enqueueing...\n"
		expectedLog := " error publishing task: " + mockTaskQueue.errPublishTask.Error() + "\n"
		lastLog := " 1 unprocessed tasks enqueued"
		fullLog := firstLog + expectedLog + lastLog

		logContents = removeTimeInfo(logContents)
		t.Run(tc, func(t *testing.T) {
			if !strings.Contains(logContents, fullLog) {
				t.Errorf("%s: expected log:\n%s but got:\n%s", tc, fullLog, logContents)
			}
		})
		mockTaskQueue.errPublishTask = nil
	}
	{
		tc := "Case 3: 5 unprocessed task found and enqueued"
		mockTaskStorer.taskModelArr = []model.MailTaskQueue{
			{UserID: 1},
			{UserID: 2},
			{UserID: 3},
			{UserID: 4},
			{UserID: 5},
		}
		var buf bytes.Buffer
		log.SetOutput(&buf)
		mockTaskService.FindUnprocessedTasksAndEnqueue()

		logContents := buf.String()
		firstLog := " Finding unprocessed tasks and enqueueing...\n"
		lastLog := " 5 unprocessed tasks enqueued"
		fullLog := firstLog + lastLog

		logContents = removeTimeInfo(logContents)
		t.Run(tc, func(t *testing.T) {
			if !strings.Contains(logContents, fullLog) {
				t.Errorf("%s: expected log:\n%s but got:\n%s", tc, fullLog, logContents)
			}
		})
	}
}

func removeTimeInfo(logContents string) string {
	return regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}`).ReplaceAllString(logContents, "")
}
