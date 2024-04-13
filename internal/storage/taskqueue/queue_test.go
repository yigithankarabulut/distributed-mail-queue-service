package taskqueue_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/gofiber/fiber/v2/log"
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/storage/taskqueue"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_taskQueue_PublishTask(t *testing.T) {
	rdb, mockClient := redismock.NewClientMock()
	taskQueue := taskqueue.New(
		taskqueue.WithConsumerCount(1),
		taskqueue.WithQueueName("testQueue"),
		taskqueue.WithRedisClient(rdb),
		taskqueue.WithTaskChannel(make(chan model.MailTaskQueue)),
	)
	{
		tc := "Case 1: Context Cancelled And Return Error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := taskQueue.PublishTask(ctx, model.MailTaskQueue{})
		t.Run(tc, func(t *testing.T) {
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Expected error to be context.Canceled, got %v", err)
			}
		})
	}
	{
		tc := "Case 2: Invalid Task Model And JSON Marshal Error"
		ctx := context.Background()
		err := taskQueue.PublishTask(ctx, "invalid task model")
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	}
	{
		tc := "Case 3: Redis LPUSH Error And Return Error"
		ctx := context.Background()
		mockClient.ExpectLPush("testQueue", "test").SetErr(errors.New("error"))
		err := taskQueue.PublishTask(ctx, "test")
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
		mockClient.ClearExpect()
	}
	{
		tc := "Case 4: Valid Task Model, Print Log And Return Nil"
		ctx := context.Background()
		var buf bytes.Buffer
		log.SetOutput(&buf)

		expectedTask := model.MailTaskQueue{UserID: 1}
		expectedJson, _ := json.Marshal(expectedTask)
		mockClient.ExpectLPush("testQueue", expectedJson).SetVal(1)

		err := taskQueue.PublishTask(ctx, expectedTask)
		want := "publishing task to channel: testQueue"
		logContents := buf.String()
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("Expected nil, got %v", err)
			}
			if !strings.Contains(logContents, want) {
				t.Errorf("Expected log \"%s\" not found in log contents:\n%s", want, logContents)
			}
		})
		mockClient.ClearExpect()
	}
}

func Test_taskQueue_SubscribeTask(t *testing.T) {
	rdb, mockClient := redismock.NewClientMock()
	taskCh := make(chan model.MailTaskQueue)
	taskQueue := taskqueue.New(
		taskqueue.WithConsumerCount(1),
		taskqueue.WithQueueName("testQueue"),
		taskqueue.WithRedisClient(rdb),
		taskqueue.WithTaskChannel(taskCh),
	)
	{
		tc := "Case 1: Context Cancelled And Return Error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := taskQueue.SubscribeTask(ctx, 1)
		want := "consumer 1 done: context canceled"
		t.Run(tc, func(t *testing.T) {
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("Expected error, got nil")
			}
		})
	}
	{
		tc := "Case 2: Redis BRPOP Error And Return Error"
		ctx := context.Background()
		mockClient.ExpectBRPop(time.Second, "testQueue").SetErr(errors.New("error"))
		err := taskQueue.SubscribeTask(ctx, 1)
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
		mockClient.ClearExpect()
	}
	{
		tc := "Case 3: Redis BRPOP Return Redis.Nil Error And Continue But JSON Unmarshal Error And Print Log"
		ctx := context.Background()
		mockClient.ExpectBRPop(time.Second, "testQueue").SetVal([]string{"test", "invalid json"})
		var buf bytes.Buffer
		log.SetOutput(&buf)

		err := taskQueue.SubscribeTask(ctx, 1)
		want := "consumer 1 error unmarshalling task"
		logContents := buf.String()
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("Expected nil, got %v", err)
			}
			if !strings.Contains(logContents, want) {
				t.Errorf("Expected log \"%s\" not found in log contents:\n%s", want, logContents)
			}
		})
		mockClient.ClearExpect()
	}
	{
		tc := "Case 4: Valid Task Model Send Task To Channel"
		ctx := context.Background()

		expectedTask := model.MailTaskQueue{UserID: 1}
		expectedJson, _ := json.Marshal(expectedTask)
		mockClient.ExpectBRPop(time.Second, "testQueue").SetVal([]string{"testQueue", string(expectedJson)})

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := taskQueue.SubscribeTask(ctx, 1); err != nil {
				log.Errorf("error: %v", err)
			}
		}()

		t.Run(tc, func(t *testing.T) {
			if task := <-taskCh; task.UserID != expectedTask.UserID {
				t.Errorf("Expected task: %v, got %v", expectedTask, task)
			}
		})
		wg.Wait()
	}
}

func Test_taskQueue_StartConsume(t *testing.T) {
	rdb, mockClient := redismock.NewClientMock()
	taskCh := make(chan model.MailTaskQueue)
	taskQueue := taskqueue.New(
		taskqueue.WithConsumerCount(1),
		taskqueue.WithQueueName("testQueue"),
		taskqueue.WithRedisClient(rdb),
		taskqueue.WithTaskChannel(taskCh),
	)
	{
		tc := "Case 1: Context Cancelled And Return Error"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		errCh := taskQueue.StartConsume(ctx)
		want := "consumer 1 done: context canceled"
		t.Run(tc, func(t *testing.T) {
			if err := <-errCh; err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("Expected error, got nil")
			}
		})
	}
	{
		tc := "Case 2: Redis BRPOP Error And Return Error"
		ctx := context.Background()
		mockClient.ExpectBRPop(time.Second, "testQueue").SetErr(errors.New("error"))
		errCh := taskQueue.StartConsume(ctx)
		t.Run(tc, func(t *testing.T) {
			if err := <-errCh; err == nil || !strings.Contains(err.Error(), "error") {
				t.Errorf("Expected error, got nil")
			}
		})
	}
	{
		tc := "Case 3: Valid Task Model Send Task To Channel"
		ctx := context.Background()

		expectedTask := model.MailTaskQueue{UserID: 1}
		expectedJson, _ := json.Marshal(expectedTask)
		mockClient.ExpectBRPop(time.Second, "testQueue").SetVal([]string{"testQueue", string(expectedJson)})
		wg := sync.WaitGroup{}
		t.Run(tc, func(t *testing.T) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				errCh := taskQueue.StartConsume(ctx)
				if err := <-errCh; err != nil {
					log.Errorf("error: %v", err)
				}
			}()
			if task := <-taskCh; task.UserID != expectedTask.UserID {
				t.Errorf("Expected task: %v, got %v", expectedTask, task)
			}
		})
		wg.Wait()
	}
}
