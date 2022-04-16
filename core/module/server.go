package module

import (
	"github.com/go-co-op/gocron"
	"github.com/go-playground/validator"
	"github.com/kataras/golog"
	"gorm.io/gorm"
)

type Server struct {
	Validator *validator.Validate
	Logger    *golog.Logger
	Scheduler *gocron.Scheduler
	Database  *gorm.DB

	Registry *Registry
}
