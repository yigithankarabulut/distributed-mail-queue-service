package model

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
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
	return d.DialAndSend(m)
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
