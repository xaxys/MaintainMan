package module

import (
	"github.com/go-co-op/gocron"
	"github.com/go-playground/validator"
	"github.com/kataras/golog"
	"github.com/olebedev/emitter"
	"gorm.io/gorm"
)

type Server struct {
	Validator *validator.Validate
	Logger    *golog.Logger
	Scheduler *gocron.Scheduler
	Database  *gorm.DB
	EventBus  *emitter.Emitter
	Registry  *Registry
}
