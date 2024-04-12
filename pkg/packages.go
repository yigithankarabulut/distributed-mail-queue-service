package pkg

import (
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/jwtutils"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/middleware"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/passutils"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/response"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/validator"
)

type Packages struct {
	JwtUtils   jwtutils.IJwtUtils
	PassUtils  passutils.IPassUtils
	Validator  validator.IValidate
	Response   response.IResponse
	Middleware middleware.IMiddleware
}

type Option func(*Packages)

func WithJwtUtils(jwtUtils jwtutils.IJwtUtils) Option {
	return func(p *Packages) {
		p.JwtUtils = jwtUtils
	}
}

func WithPassUtils(passUtils passutils.IPassUtils) Option {
	return func(p *Packages) {
		p.PassUtils = passUtils
	}
}

func WithValidator(validator validator.IValidate) Option {
	return func(p *Packages) {
		p.Validator = validator
	}
}

func WithResponse(response response.IResponse) Option {
	return func(p *Packages) {
		p.Response = response
	}
}

func WithMiddleware(middleware middleware.IMiddleware) Option {
	return func(p *Packages) {
		p.Middleware = middleware
	}
}

func New(opts ...Option) *Packages {
	p := &Packages{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
