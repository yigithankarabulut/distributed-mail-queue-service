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
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     int    `json:"smtp_port"`
	SmtpUsername string `json:"smtp-username"`
}

func (r *GetUserResponse) FromUser(user model.User) {
	r.ID = user.ID
	r.Email = user.Email
	r.SmtpHost = user.SmtpHost
	r.SmtpPort = user.SmtpPort
	r.SmtpUsername = user.SmtpUsername
}
