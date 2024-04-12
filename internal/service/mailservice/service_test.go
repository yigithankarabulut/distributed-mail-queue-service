package mailservice_test

import (
	"github.com/yigithankarabulut/distributed-mail-queue-service/internal/service/mailservice"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"testing"
)

func Test_mailService_AddTask(t *testing.T) {
	{
		tc := "Case 1: All Fields Are Valid And Should Return Success"
		mockService := mailservice.New()
		task := model.MailTaskQueue{
			User: model.User{
				Password:     "test",
				Email:        "test@test.com",
				SmtpHost:     "smtp.test.com",
				SmtpPort:     587,
				SmtpUsername: "test",
				SmtpPassword: "test",
			},
			UserID:         1,
			RecipientEmail: "example@ex.com",
			Subject:        "Test",
			Body:           "Test",
		}
		err := mockService.AddTask(task)
		t.Run(tc, func(t *testing.T) {
			if err != nil {
				t.Errorf("Expected error to be nil but got %v", err)
			}
		})
	}
	{
		tc := "Case 2: Missing Task Fields And Should Return Error"
		mockService := mailservice.New()
		task := model.MailTaskQueue{
			User: model.User{
				Password:     "test",
				Email:        "test@test.com",
				SmtpHost:     "smtp.test.com",
				SmtpPort:     587,
				SmtpUsername: "test",
				SmtpPassword: "test",
			},
		}
		err := mockService.AddTask(task)
		want := "Missing fields: RecipientEmail, Subject, Body"
		t.Run(tc, func(t *testing.T) {
			if err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
	}
	{
		tc := "Case 3: Missing User Fields And Should Return Error"
		mockService := mailservice.New()
		task := model.MailTaskQueue{
			User: model.User{
				Password: "test",
			},
			UserID:         1,
			RecipientEmail: "example@ex.com",
			Subject:        "Test",
			Body:           "Test",
		}
		err := mockService.AddTask(task)
		want := "Missing fields: User.Email, User.SmtpHost, User.SmtpPort, User.SmtpUsername, User.SmtpPassword"
		t.Run(tc, func(t *testing.T) {
			if err == nil || err.Error() != want {
				t.Errorf("Expected error to be %s but got %v", want, err)
			}
		})
	}
}

func Test_mailService_NewDialer(t *testing.T) {
	mockService := mailservice.New()
	tc := "Case 1: New dialer should return dialer with tls config"
	dialer := mockService.NewDialer()
	t.Run(tc, func(t *testing.T) {
		if dialer.TLSConfig == nil {
			t.Errorf("expected tls config, got nil")
		}
	})
}

func Test_mailService_NewMessage(t *testing.T) {
	mockService := mailservice.New(
		mailservice.WithTask(model.MailTaskQueue{
			User: model.User{
				Email:        "test@test.com",
				SmtpHost:     "smtp.test.com",
				SmtpPort:     587,
				SmtpUsername: "test",
				SmtpPassword: "test",
			},
			RecipientEmail: "example@ex.com",
			Subject:        "Test",
			Body:           "Test",
		}),
	)
	{
		tc := "Case 1: New message should return message with from"
		message := mockService.NewMessage()
		t.Run(tc, func(t *testing.T) {
			if len(message.GetHeader("From")) == 0 {
				t.Errorf("expected from, got empty")
			}
		})
	}
	{
		tc := "Case 2: New message should return message with to"
		message := mockService.NewMessage()
		t.Run(tc, func(t *testing.T) {
			if len(message.GetHeader("To")) == 0 {
				t.Errorf("expected to, got empty")
			}
		})
	}
	{
		tc := "Case 3: New message should return message with subject"
		message := mockService.NewMessage()
		t.Run(tc, func(t *testing.T) {
			if len(message.GetHeader("Subject")) == 0 {
				t.Errorf("expected subject, got empty")
			}
		})
	}
}

func Test_mailService_SendMail(t *testing.T) {
	{
		mockService := mailservice.New(
			mailservice.WithTask(model.MailTaskQueue{
				User: model.User{
					Email:        "test@test.com",
					SmtpHost:     "smtp.test.com",
					SmtpPort:     587,
					SmtpUsername: "test",
					SmtpPassword: "test",
				},
				RecipientEmail: "example@ex.com",
				Subject:        "Test",
				Body:           "Test",
			}),
		)
		tc := "Case 1: Send mail should return error when dialer returns error"
		dialer := &mockDialer{errDial: ErrorDialer, sender: mockSender{}}
		err := mockService.SendMail(dialer, mockService.NewMessage())
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("Expected error to be not nil but got nil")
			}
		})
	}
	{
		mockService := mailservice.New(
			mailservice.WithTask(model.MailTaskQueue{
				User: model.User{
					Email:        "test@test.com",
					SmtpHost:     "smtp.test.com",
					SmtpPort:     587,
					SmtpUsername: "test",
					SmtpPassword: "test",
				},
				RecipientEmail: "example@ex.com",
				Subject:        "Test",
				Body:           "Test",
			}),
		)
		tc := "Case 2: Send mail should return error when sender returns error"
		dialer := &mockDialer{errDial: nil, sender: mockSender{errSend: ErrorSenders, errClose: ErrorSenders}}
		err := mockService.SendMail(dialer, mockService.NewMessage())
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("Expected error to be not nil but got nil")
			}
		})
	}
	{
		mockService := mailservice.New(
			mailservice.WithTask(model.MailTaskQueue{
				User: model.User{
					Email:        "test@test.com",
					SmtpHost:     "smtp.test.com",
					SmtpPort:     587,
					SmtpUsername: "test",
					SmtpPassword: "test",
				},
				RecipientEmail: "example@ex.com",
				Subject:        "Test",
				Body:           "Test",
			}),
		)
		tc := "Case 3: Send mail should return nil when dialer and sender are successful"
		dialer := &mockDialer{errDial: nil, sender: mockSender{}}
		err := mockService.SendMail(dialer, mockService.NewMessage())
		t.Run(tc, func(t *testing.T) {
			if err != nil && err.Error() != "Random error" {
				t.Errorf("Expected error to be nil but got %v", err)
			}
		})
	}
	{
		mockService := mailservice.New(
			mailservice.WithTask(model.MailTaskQueue{
				User: model.User{
					Email:        "test@test.com",
					SmtpHost:     "smtp.test.com",
					SmtpPort:     587,
					SmtpUsername: "test",
					SmtpPassword: "test",
				},
				RecipientEmail: "example@ex.com",
				Subject:        "Test",
				Body:           "Test",
			}),
		)
		tc := "Case 4: Send mail should return error when Close returns error"
		dialer := &mockDialer{errDial: nil, sender: mockSender{errClose: ErrorSenders, errSend: ErrorSenders}}
		err := mockService.SendMail(dialer, mockService.NewMessage())
		t.Run(tc, func(t *testing.T) {
			if err == nil {
				t.Errorf("Expected error to be not nil but got nil")
			}
		})
	}
}
