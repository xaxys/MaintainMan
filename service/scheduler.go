package service

import (
	"maintainman/config"
	"time"

	"github.com/go-co-op/gocron"
)

var (
	Scheduler = gocron.NewScheduler(time.UTC)
)

func init() {
	duration := config.AppConfig.GetDuration("app.appraise.purge")
	Scheduler.Every(duration).SingletonMode().Do(AutoAppraiseOrder)
}
