package model

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

type GoMail struct {
	From         string
	To           string
	Subject      string
	Body         string
	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
}

func (g *GoMail) NewMessage() *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", g.From)
	m.SetHeader("To", g.To)
	m.SetHeader("Subject", g.Subject)
	m.SetBody("text/plain", g.Body)
	return m
}

func (g *GoMail) NewDialer() *gomail.Dialer {
	d := gomail.NewDialer(g.SmtpHost, g.SmtpPort, g.SmtpUsername, g.SmtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d
}

func (g *GoMail) Send(d *gomail.Dialer, m *gomail.Message) error {

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

func Example() {
	mail := GoMail{
		From:         "yigithannkarabulutt@gmail.com",
		To:           "dobrainmusic@gmail.com",
		Subject:      "Test Mail",
		Body:         "Hello, this is a test email!",
		SmtpHost:     "smtp.gmail.com",
		SmtpPort:     587,
		SmtpUsername: "yigithannkarabulutt@gmail.com",
		SmtpPassword: "aowjdppvonjvayng",
	}

	if err := mail.Send(mail.NewDialer(), mail.NewMessage()); err != nil {
		panic(err)
	}

	mail2 := GoMail{
		From:         "yigithannkarabulutt@gmail.com",
		To:           "dobrainmusic@gmail.com",
		Subject:      "Test Mail",
		Body:         "Hello, this is a test email!",
		SmtpHost:     "smtp.gmail.com",
		SmtpPort:     587,
		SmtpUsername: "yigithannkarabulutt@gmail.com",
		SmtpPassword: "aowjdppvonjvayng",
	}

	if err := mail2.Send(mail2.NewDialer(), mail2.NewMessage()); err != nil {
		panic(err)
	}
}
