package dtoreq

import "github.com/yigithankarabulut/distributed-mail-queue-service/src/model"

type RegisterRequest struct {
	Email        string `json:"email" query:"-" validate:"required,email"`
	Password     string `json:"password" query:"-" validate:"required"`
	SmtpHost     string `json:"smtp_host" query:"-" validate:"required"`
	SmtpPort     int    `json:"smtp_port" query:"-" validate:"required"`
	SmtpUsername string `json:"smtp-username" query:"-" validate:"required"`
	SmtpPassword string `json:"smtp-password" query:"-" validate:"required"`
}

func (r RegisterRequest) ConvertToUser() model.User {
	return model.User{
		Email:        r.Email,
		Password:     r.Password,
		SmtpHost:     r.SmtpHost,
		SmtpPort:     r.SmtpPort,
		SmtpUsername: r.SmtpUsername,
		SmtpPassword: r.SmtpPassword,
	}
}

type LoginRequest struct {
	Email    string `json:"email" query:"-" validate:"required,email"`
	Password string `json:"password" query:"-" validate:"required"`
}

type UpdateUserRequest struct {
	Email        string `json:"email" validate:"required"`
	Token        string `json:"token" validate:"required"`
	SmtpHost     string `json:"smtp_host" validate:"required"`
	SmtpPort     int    `json:"smtp_port" validate:"required"`
	SmtpUsername string `json:"smtp-username" validate:"required"`
	SmtpPassword string `json:"smtp-password" validate:"required"`
}

type GetUserRequest struct {
	ID uint `json:"-" query:"id" validate:"required,numeric"`
}
