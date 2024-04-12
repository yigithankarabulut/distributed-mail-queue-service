package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

type IValidate interface {
	BindAndValidate(c *fiber.Ctx, data interface{}) error
}

type Validator struct {
	Validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		Validator: validator.New(),
	}
}

func (v *Validator) RegisterValidation(data interface{}) error {
	v.Validator.RegisterTagNameFunc(v.getTagNameFunc("json"))
	v.Validator.RegisterTagNameFunc(v.getTagNameFunc("query"))
	return v.Validator.Struct(data)
}

func (v *Validator) getTagNameFunc(tag string) func(fld reflect.StructField) string {
	return func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	}
}

func (v *Validator) parseAndValidate(c *fiber.Ctx, data interface{}) error {
	var (
		err  error
		err2 error
	)
	err = c.BodyParser(data)
	err2 = c.QueryParser(data)
	if err != nil && err2 != nil {
		return err
	}
	if err := v.RegisterValidation(data); err != nil {
		return err
	}
	return nil
}

func (v *Validator) BindAndValidate(c *fiber.Ctx, data interface{}) error {
	if err := v.parseAndValidate(c, data); err != nil {
		return err
	}
	return nil
}
