package dtores

import "github.com/yigithankarabulut/distributed-mail-queue-service/model"

type RegisterUserResponse struct {
}

type LoginResponse struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

type UpdateUserResponse struct {
}

type GetUserResponse struct {
	Email        string `json:"email"`
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     int    `json:"smtp_port"`
	SmtpUsername string `json:"smtp-username"`
	SmtpPassword string `json:"smtp-password"`
}

func (r *GetUserResponse) FromUser(user model.User) {
	r.Email = user.Email
	r.SmtpHost = user.SmtpHost
	r.SmtpPort = user.SmtpPort
	r.SmtpUsername = user.SmtpUsername
	r.SmtpPassword = user.SmtpPassword
}
