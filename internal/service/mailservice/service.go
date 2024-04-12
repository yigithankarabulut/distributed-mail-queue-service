package mailservice

import (
	"crypto/tls"
	"errors"
	"github.com/yigithankarabulut/distributed-mail-queue-service/model"
	"gopkg.in/gomail.v2"
	"log"
	"reflect"
	"strings"
	"time"
)

func (s *mailService) AddTask(task model.MailTaskQueue) error {
	v := reflect.ValueOf(task)
	t := reflect.TypeOf(task)
	var missingFields []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			if t.Field(i).Name == "Model" || t.Field(i).Name == "ScheduledAt" {
				continue
			}
			for j := 0; j < field.NumField(); j++ {
				if t.Field(j).Name == "Model" {
					continue
				}
				if field.Field(j).IsZero() {
					missingFields = append(missingFields, t.Field(i).Name+"."+field.Type().Field(j).Name)
				}
			}
		} else {
			if t.Field(i).Name == "Status" || t.Field(i).Name == "TryCount" || t.Field(i).Name == "CreatedAt" || t.Field(i).Name == "UpdatedAt" || t.Field(i).Name == "UserID" {
				continue
			}
			if field.IsZero() {
				missingFields = append(missingFields, t.Field(i).Name)
			}
		}
	}
	if len(missingFields) > 0 {
		return errors.New("Missing fields: " + strings.Join(missingFields, ", "))
	}
	s.From = task.User.Email
	s.To = task.RecipientEmail
	s.Subject = task.Subject
	s.Body = task.Body
	s.SmtpHost = task.User.SmtpHost
	s.SmtpPort = task.User.SmtpPort
	s.SmtpUsername = task.User.SmtpUsername
	s.SmtpPassword = task.User.SmtpPassword
	return nil
}

func (s *mailService) NewDialer() *gomail.Dialer {
	d := gomail.NewDialer(s.SmtpHost, s.SmtpPort, s.SmtpUsername, s.SmtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d
}

func (s *mailService) NewMessage() *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", s.To)
	m.SetHeader("Subject", s.Subject)
	m.SetBody("text/plain", s.Body)
	return m
}

func (s *mailService) throwRandomError() error {
	num := time.Now().Nanosecond()
	if num%2 == 0 {
		return nil
	}
	return errors.New("Random error")
}

func (s *mailService) SendMail(d Dialer, m *gomail.Message) error {
	if err := s.throwRandomError(); err != nil {
		return err
	}
	ch := make(chan *gomail.Message)
	errChan := make(chan error, 1)
	go func() {
		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					close(errChan)
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						log.Print(err)
						errChan <- err
						close(errChan)
						return
					}
					open = true
				}
				if err := gomail.Send(s, m); err != nil {
					log.Print(err)
					errChan <- err
				}
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						log.Print(err)
						errChan <- err
						close(errChan)
						return
					}
					open = false
				}
			}
		}
	}()

	// Use the channel in your program to send emails.
	ch <- m
	// Close the channel to stop the mail daemon.
	close(ch)
	err := <-errChan
	return err
}
