package workerservice_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/workerservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/constant"
	"strings"
	"testing"
)

func Test_worker_TriggerWorker(t *testing.T) {
	mockTaskStorer := &mockTaskStorer{}
	mockMailService := &mockMailService{}
	mockTaskQueue := &mockTaskQueue{}
	{
		done := make(chan struct{})
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
			workerservice.WithDoneChannel(done),
		)
		tc := "Case 1: Done channel is closed and returns error"
		go func() {
			close(done)
		}()
		err := mockWorkerService.TriggerWorker()
		want := "worker 1 done"
		t.Run(tc, func(t *testing.T) {
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s: expected %v but got %v", tc, want, err)
			}
		})
	}
	{
		done := make(chan struct{})
		taskChannel := make(chan model.MailTaskQueue)
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
			workerservice.WithChannel(taskChannel),
			workerservice.WithDoneChannel(done),
		)
		tc := "Case 2: Task channel is closed and returns error"
		go func() {
			close(taskChannel)
		}()
		err := mockWorkerService.TriggerWorker()
		want := "worker 1 task channel closed"
		t.Run(tc, func(t *testing.T) {
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s: expected %v but got %v", tc, want, err)
			}
		})
	}
	{
		done := make(chan struct{})
		taskChannel := make(chan model.MailTaskQueue)
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
			workerservice.WithChannel(taskChannel),
			workerservice.WithDoneChannel(done),
		)
		tc := "Case 3: Task MaxTryCount reached and print logs"

		mockMailService.errSendMail = errors.New("send mail error")
		var buf bytes.Buffer
		log.SetOutput(&buf)

		go func() {
			taskChannel <- model.MailTaskQueue{
				TryCount:       4,
				RecipientEmail: "test@test.com",
			}
			close(taskChannel)
			close(done)
		}()
		_ = mockWorkerService.TriggerWorker()

		logContents := buf.String()
		expectedLogs := []string{
			"worker 1 sending mail to test@test.com",
			"worker 1 error sending mail to test@test.com: send mail error",
			"worker 1 error handling task: task 0 cancelled after 5 tries",
		}
		t.Run(tc, func(t *testing.T) {
			for _, expectedLog := range expectedLogs {
				if !strings.Contains(logContents, expectedLog) {
					t.Errorf("Expected log \"%s\" not found in log contents:\n%s", expectedLog, logContents)
				}
			}
		})
		mockMailService.errSendMail = nil
	}
}

func Test_worker_HandleTask(t *testing.T) {
	mockTaskStorer := &mockTaskStorer{}
	mockMailService := &mockMailService{}
	mockTaskQueue := &mockTaskQueue{}
	{
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
		)
		tc := "Case 1: Context is cancelled and returns error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := mockWorkerService.HandleTask(ctx, model.MailTaskQueue{})
		want := "context canceled"
		t.Run(tc, func(t *testing.T) {
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s: expected %v but got %v", tc, want, err)
			}
		})
	}
	{
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
		)
		tc := "Case 2: Invalid task and MailService.AddTask returns error"
		mockMailService.errAddTask = errors.New("add task error")
		err := mockWorkerService.HandleTask(context.Background(), model.MailTaskQueue{})
		want := "worker 1 error adding task: add task error"
		t.Run(tc, func(t *testing.T) {
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s: expected %v but got %v", tc, want, err)
			}
		})
		mockMailService.errAddTask = nil
	}
	{
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
		)
		tc := "Case 3: SendMail returns error and handleError return max try count error"
		mockMailService.errSendMail = errors.New("send mail error")
		err := mockWorkerService.HandleTask(context.Background(), model.MailTaskQueue{
			TryCount: constant.MaxTryCount,
		})
		want := "task 0 cancelled after 4 tries"
		t.Run(tc, func(t *testing.T) {
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("%s: expected %v but got %v", tc, want, err)
			}
		})
		mockMailService.errSendMail = nil
	}
	{
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
		)
		tc := "Case 4: TaskStorage.Update returns error and print logs"
		mockTaskStorer.errUpdate = errors.New("update error")
		var buf bytes.Buffer
		log.SetOutput(&buf)

		_ = mockWorkerService.HandleTask(context.Background(), model.MailTaskQueue{})
		want := "worker 1 error updating task: update error"
		logContents := buf.String()

		t.Run(tc, func(t *testing.T) {
			if !strings.Contains(logContents, want) {
				t.Errorf("Expected log \"%s\" not found in log contents:\n%s", want, logContents)
			}
		})
	}
	{
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
		)
		tc := "Case 5: Task TryCount is less than MaxTryCount but TaskStorage.Update returns error and print logs"
		mockTaskStorer.errUpdate = errors.New("update error")
		mockMailService.errSendMail = errors.New("send mail error")
		var buf bytes.Buffer
		log.SetOutput(&buf)

		_ = mockWorkerService.HandleTask(context.Background(), model.MailTaskQueue{
			TryCount: 1,
		})
		want := "worker 1 error updating task: update error"
		logContents := buf.String()

		t.Run(tc, func(t *testing.T) {
			if !strings.Contains(logContents, want) {
				t.Errorf("Expected log \"%s\" not found in log contents:\n%s", want, logContents)
			}
		})
		mockTaskStorer.errUpdate = nil
		mockMailService.errSendMail = nil
	}
	{
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
		)
		tc := "Case 6: Task TryCount is less than MaxTryCount but TaskQueue.PublishTask returns error and print logs"
		mockTaskQueue.errPublishTask = errors.New("publish task error")
		mockMailService.errSendMail = errors.New("send mail error")
		var buf bytes.Buffer
		log.SetOutput(&buf)

		_ = mockWorkerService.HandleTask(context.Background(), model.MailTaskQueue{
			TryCount: 1,
		})
		want := "worker 1 error publishing task: publish task error"
		logContents := buf.String()

		t.Run(tc, func(t *testing.T) {
			if !strings.Contains(logContents, want) {
				t.Errorf("Expected log \"%s\" not found in log contents:\n%s", want, logContents)
			}
		})
		mockTaskQueue.errPublishTask = nil
		mockMailService.errSendMail = nil
	}
	{
		mockWorkerService := workerservice.New(
			workerservice.WithID(1),
			workerservice.WithTaskStorage(mockTaskStorer),
			workerservice.WithTaskQueue(mockTaskQueue),
			workerservice.WithMailService(mockMailService),
		)
		tc := "Case 5: All operations are successful"
		var buf bytes.Buffer
		log.SetOutput(&buf)

		_ = mockWorkerService.HandleTask(context.Background(), model.MailTaskQueue{
			RecipientEmail: "test@test.com",
		})
		logContents := buf.String()
		expectedLogs := []string{
			"worker 1 sending mail to test@test.com",
			"worker 1 sent mail to test@test.com",
		}

		t.Run(tc, func(t *testing.T) {
			for _, expectedLog := range expectedLogs {
				if !strings.Contains(logContents, expectedLog) {
					t.Errorf("Expected log \"%s\" not found in log contents:\n%s", expectedLog, logContents)
				}
			}
		})
	}
}
