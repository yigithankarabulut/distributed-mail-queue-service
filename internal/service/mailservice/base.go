package mailservice

import (
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gopkg.in/gomail.v2"
)

type MailService interface {
	NewDialer() *gomail.Dialer
	NewMessage() *gomail.Message
	SendMail(d *gomail.Dialer, m *gomail.Message) error
}

type mailService struct {
	From         string
	To           string
	Subject      string
	Body         string
	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
}

type Option func(*mailService)

func WithTask(task model.MailTaskQueue) Option {
	return func(m *mailService) {
		m.From = task.User.Email
		m.To = task.RecipientEmail
		m.Subject = task.Subject
		m.Body = task.Body
		m.SmtpHost = task.User.SmtpHost
		m.SmtpPort = task.User.SmtpPort
		m.SmtpUsername = task.User.SmtpUsername
		m.SmtpPassword = task.User.SmtpPassword
	}
}

func New(opts ...Option) MailService {
	mail := &mailService{}
	for _, opt := range opts {
		opt(mail)
	}
	return mail
}
