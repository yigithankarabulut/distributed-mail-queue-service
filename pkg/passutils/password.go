package passutils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type IPassUtils interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash, pass string) error
}

type PassUtils struct{}

func New() *PassUtils {
	return &PassUtils{}
}

func (p *PassUtils) HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password is empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func (p *PassUtils) ComparePassword(hash, password string) error {
	if len(hash) == 0 {
		return errors.New("password hash is empty")
	}
	if len(password) == 0 {
		return errors.New("password is empty")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("invalid password")
		}
		return err
	}
	return nil
}
