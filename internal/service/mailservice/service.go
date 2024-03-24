package mailservice

import (
	"crypto/tls"
	"errors"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

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

func throwRandomError() error {
	num := time.Now().Nanosecond()
	if num%2 == 0 {
		return nil
	}
	return errors.New("Random error")
}

func (s *mailService) SendMail(d *gomail.Dialer, m *gomail.Message) error {
	if err := throwRandomError(); err != nil {
		return err
	}
	ch := make(chan *gomail.Message)

	go func() {
		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						log.Print(err)
						return
					}
					open = true
				}
				if err := gomail.Send(s, m); err != nil {
					log.Print(err)
				}
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						log.Print(err)
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
	return nil
}
