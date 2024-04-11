package taskqueue

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"testing"
)

func Test_taskQueue_PublishTask(t *testing.T) {
	rdb, mockClient := redismock.NewClientMock()
	type fields struct {
		consumerCount int
		queueName     string
		rdb           *redis.Client
		taskChannel   chan model.MailTaskQueue
	}
	type args struct {
		task interface{}
	}

	task1 := model.MailTaskQueue{
		UserID: 1,
	}
	task2 := model.MailTaskQueue{
		UserID:         2,
		Body:           "test",
		RecipientEmail: "test@test.com",
		Subject:        "test",
	}
	task3 := "invalid task model"
	expectTask1, _ := json.Marshal(task1)
	expectTask2, _ := json.Marshal(task2)

	TestCase := []struct {
		name    string
		fields  fields
		args    args
		expect  []byte
		wantErr bool
	}{
		{
			name: "Test Case 1 - Valid Task Model",
			fields: fields{
				consumerCount: 1,
				queueName:     "testQueue",
				rdb:           rdb,
				taskChannel:   make(chan model.MailTaskQueue),
			},
			args: args{
				task: task1,
			},
			expect:  expectTask1,
			wantErr: false,
		},
		{
			name: "Test Case 2 - Valid Task Model",
			fields: fields{
				consumerCount: 1,
				queueName:     "testQueue",
				rdb:           rdb,
				taskChannel:   make(chan model.MailTaskQueue),
			},
			args: args{
				task: task2,
			},
			expect:  expectTask2,
			wantErr: false,
		},
		{
			name: "Test Case 3 - Invalid Task Model - JSON Marshalling Error",
			fields: fields{
				consumerCount: 1,
				queueName:     "testQueue",
				rdb:           rdb,
				taskChannel:   make(chan model.MailTaskQueue),
			},
			args: args{
				task: task3,
			},
			expect:  nil,
			wantErr: true,
		},
		{
			name: "Test Case 4 - Redis Error",
			fields: fields{
				consumerCount: 1,
				queueName:     "testQueue",
				rdb:           rdb,
				taskChannel:   make(chan model.MailTaskQueue),
			},
			args: args{
				task: task1,
			},
			expect:  expectTask1,
			wantErr: true,
		},
	}

	for _, tt := range TestCase {
		t.Run(tt.name, func(t *testing.T) {
			r := &taskQueue{
				consumerCount: tt.fields.consumerCount,
				queueName:     tt.fields.queueName,
				rdb:           tt.fields.rdb,
				taskChannel:   tt.fields.taskChannel,
			}
			if tt.wantErr {
				mockClient.ExpectLPush(tt.fields.queueName, tt.expect).SetErr(errors.New("error"))
			} else {
				mockClient.ExpectLPush(tt.fields.queueName, tt.expect).SetVal(1)
			}
			if err := r.PublishTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("PublishTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}

}

func Test_taskQueue_SubscribeTask(t *testing.T) {
	rdb, mockClient := redismock.NewClientMock()
	type fields struct {
		consumerCount int
		queueName     string
		rdb           *redis.Client
		taskChannel   chan model.MailTaskQueue
	}
	type args struct {
		consumerID int
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		routineFunc func(args ...interface{})
		wantErr     bool
	}{
		{
			name: "Test Case 1 - Valid Consumer ID",
			fields: fields{
				consumerCount: 1,
				queueName:     "testQueue",
				rdb:           rdb,
				taskChannel:   make(chan model.MailTaskQueue),
			},
			args: args{
				consumerID: 1,
			},
			routineFunc: func(args ...interface{}) {
				args[0].(func())()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &taskQueue{
				consumerCount: tt.fields.consumerCount,
				queueName:     tt.fields.queueName,
				rdb:           tt.fields.rdb,
				taskChannel:   tt.fields.taskChannel,
			}
			go func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if err := r.SubscribeTask(ctx, tt.args.consumerID); (err != nil) != tt.wantErr {
					t.Errorf("SubscribeTask() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			tt.routineFunc(func() {
				mockClient.ExpectLPop(tt.fields.queueName).SetErr(errors.New("error"))
			})
		})
	}
}

func Test_taskQueue_StartConsume(t *testing.T) {
	rdb, mockClient := redismock.NewClientMock()
	type fields struct {
		consumerCount int
		queueName     string
		rdb           *redis.Client
		taskChannel   chan model.MailTaskQueue
	}
	type args struct {
		consumerID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		expect  []byte
		ctx     int // 0: context.Background(), 1: context.WithCancel(context.Background())
		wantErr bool
	}{
		{
			name: "Test Case 1 - Valid Consumer Count and Queue Name",
			fields: fields{
				consumerCount: 1,
				queueName:     "testQueue",
				rdb:           rdb,
				taskChannel:   make(chan model.MailTaskQueue),
			},
			args: args{
				consumerID: 1,
			},
			ctx:     1,
			expect:  []byte(`{"id":1,"user_id":0,"subject":"","body":"","recipient_email":""}`),
			wantErr: false,
		},
		{
			name: "Test Case 2 - Redis Error",
			fields: fields{
				consumerCount: 1,
				queueName:     "testQueue",
				rdb:           rdb,
				taskChannel:   make(chan model.MailTaskQueue),
			},
			args: args{
				consumerID: 1,
			},
			ctx:     1,
			expect:  []byte(`{"id":1,"user_id":0,"subject":"","body":"","recipient_email":""}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &taskQueue{
				consumerCount: tt.fields.consumerCount,
				queueName:     tt.fields.queueName,
				rdb:           tt.fields.rdb,
				taskChannel:   tt.fields.taskChannel,
			}
			if tt.wantErr {
				mockClient.ExpectBRPop(0, tt.fields.queueName).SetErr(errors.New("error"))
			} else {
				mockClient.ExpectBRPop(0, tt.fields.queueName).SetVal([]string{"test", string(tt.expect)})
			}
			var ctx context.Context
			var cancel context.CancelFunc
			if tt.ctx == 0 {
				ctx = context.Background()
			} else {
				ctx, cancel = context.WithCancel(context.Background())
				cancel()
			}
			errCh := r.StartConsume(ctx)
			if err := <-errCh; err != nil {
				if tt.ctx == 1 && err.Error() == "context canceled" {
					return
				} else if (err != nil) != tt.wantErr {
					t.Errorf("StartConsume() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
