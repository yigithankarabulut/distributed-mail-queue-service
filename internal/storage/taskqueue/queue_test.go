package taskqueue

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"testing"
)

//func TestTaskQueue_PublishTask(t *testing.T) {
//	rdb, mock := redismock.NewClientMock()
//	queue := taskqueue.New(
//		taskqueue.WithRedisClient(rdb),
//		taskqueue.WithQueueName("testQueue"),
//		taskqueue.WithConsumerCount(1),
//		taskqueue.WithTaskChannel(make(chan model.MailTaskQueue)),
//	)
//
//	// Beklenen JSON verisi
//	task := model.MailTaskQueue{
//		UserID: 1,
//	}
//	jsonTask, _ := json.Marshal(task)
//
//	// Redis mock'u için LPush beklentisi
//	mock.ExpectLPush("testQueue", jsonTask).SetVal(1)
//
//	// PublishTask metodunu çağırın
//	err := queue.PublishTask(task)
//	assert.NoError(t, err)
//
//	// Beklenen tüm mock çağrılarını doğrulayın
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("mock expectations were not met: %s", err)
//	}
//}

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
	}

	for _, tt := range TestCase {
		t.Run(tt.name, func(t *testing.T) {
			r := &taskQueue{
				consumerCount: tt.fields.consumerCount,
				queueName:     tt.fields.queueName,
				rdb:           tt.fields.rdb,
				taskChannel:   tt.fields.taskChannel,
			}
			mockClient.ExpectLPush(tt.fields.queueName, tt.expect).SetVal(1)
			if err := r.PublishTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("PublishTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}

}

func Test_taskQueue_StartConsume(t *testing.T) {
	rdb, _ := redismock.NewClientMock()
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
			wantErr: false,
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
			if err := r.StartConsume(); (err != nil) != tt.wantErr {
				t.Errorf("StartConsume() error = %v, wantErr %v", err, tt.wantErr)
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
				if err := r.SubscribeTask(tt.args.consumerID); (err != nil) != tt.wantErr {
					t.Errorf("SubscribeTask() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			tt.routineFunc(func() {
				mockClient.ExpectLPop(tt.fields.queueName).SetErr(errors.New("error"))
			})
		})
	}
}
