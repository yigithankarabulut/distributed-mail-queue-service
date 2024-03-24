package dtores

import "github.com/yigithankarabulut/distributed-mail-queue-service/model"

type BaseTaskResponse struct {
	TaskID         uint   `json:"task_id"`
	Status         int    `json:"status"`
	TryCount       int    `json:"try_count"`
	RecipientEmail string `json:"recipient_email"`
	Subject        string `json:"subject"`
	Body           string `json:"body"`
}

type TaskEnqueueResponse struct {
	TaskID uint `json:"task_id"`
}

type GetAllQueuedTasksResponse struct {
	Tasks []BaseTaskResponse `json:"tasks"`
}

type GetAllFailedTasksResponse struct {
	Tasks []BaseTaskResponse `json:"tasks"`
}

func (r *GetAllQueuedTasksResponse) ToMailTaskQueue(tasks []model.MailTaskQueue) {
	for _, task := range tasks {
		r.Tasks = append(r.Tasks, BaseTaskResponse{
			TaskID:         task.ID,
			Status:         task.Status,
			TryCount:       task.TryCount,
			RecipientEmail: task.RecipientEmail,
			Subject:        task.Subject,
			Body:           task.Body,
		})
	}
}

func (r *GetAllFailedTasksResponse) ToMailTaskQueue(tasks []model.MailTaskQueue) {
	for _, task := range tasks {
		r.Tasks = append(r.Tasks, BaseTaskResponse{
			TaskID:         task.ID,
			Status:         task.Status,
			TryCount:       task.TryCount,
			RecipientEmail: task.RecipientEmail,
			Subject:        task.Subject,
			Body:           task.Body,
		})
	}
}
