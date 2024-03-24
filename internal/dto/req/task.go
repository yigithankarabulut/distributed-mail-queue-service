package dtoreq

import "github.com/yigithankarabulut/distributed-mail-queue-service/model"

type TaskEnqueueRequest struct {
	RecipientEmail string `json:"recipient_email" query:"-" validate:"required,email"`
	Subject        string `json:"subject" query:"-" validate:"required"`
	Body           string `json:"body" query:"-" validate:"required"`
	ScheduledAt    string `json:"scheduled_at" query:"-" validate:"omitempty"`
	UserID         uint   `json:"-" query:"-" validate:"required,numeric"`
}

type GetAllQueuedTasksRequest struct {
	UserID uint `json:"-" query:"-" validate:"required,numeric"`
}

type GetAllFailedTasksRequest struct {
	UserID uint `json:"-" query:"-" validate:"required,numeric"`
}

func (r TaskEnqueueRequest) ConvertToMailTaskQueue() model.MailTaskQueue {
	return model.MailTaskQueue{
		RecipientEmail: r.RecipientEmail,
		Subject:        r.Subject,
		Body:           r.Body,
		UserID:         r.UserID,
	}
}
