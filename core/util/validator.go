package util

import (
	"regexp"

	"github.com/go-playground/validator"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New()
	orderByRegex := regexp.MustCompile(`^[^@ \t\r\n]+([ ]desc|[ ]asc)?$`)
	Validator.RegisterValidation("order_by", func(fl validator.FieldLevel) bool {
		if orderByRegex.MatchString(fl.Field().String()) {
			return true
		}
		return false
	})
}
