package mailservice_test

import (
	"errors"
	"gopkg.in/gomail.v2"
	"io"
)

var (
	ErrorDialer  = errors.New("error dialer")
	ErrorSenders = errors.New("error senders")
)

type mockDialer struct {
	errDial error
	sender  mockSender
}

func (m *mockDialer) Dial() (gomail.SendCloser, error) {
	return &m.sender, m.errDial
}

type mockSender struct {
	errSend  error
	errClose error
}

func (m *mockSender) Send(from string, to []string, msg io.WriterTo) error {
	return m.errSend
}

func (m *mockSender) Close() error {
	return m.errClose
}
