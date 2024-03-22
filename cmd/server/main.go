package main

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func main() {
	m := gomail.NewMessage()
	m.SetHeader("From", "yigithannkarabulutt@gmail.com")
	m.SetHeader("To", "dobrainmusic@gmail.com")
	m.SetHeader("Subject", "Test Mail")
	m.SetBody("text/plain", "Hello, this is a test email!")

	d := gomail.NewDialer("smtp.gmail.com", 587, "yigithannkarabulutt@gmail.com", "aowjdppvonjvayng")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
