package pkg

import (
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/jwtutils"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/passutils"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/response"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/validator"
)

var PackagesInstance *Packages

type Packages struct {
	JwtUtils  jwtutils.IJwtUtils
	PassUtils passutils.IPassUtils
	Validator validator.IValidate
	Response  response.IResponse
}

func New() *Packages {
	if PackagesInstance != nil {
		return PackagesInstance
	}
	return &Packages{
		JwtUtils:  jwtutils.New(),
		PassUtils: passutils.New(),
		Validator: validator.New(),
		Response:  response.New(),
	}
}
