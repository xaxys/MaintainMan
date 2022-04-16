package service

import (
	"time"

	"github.com/go-co-op/gocron"
)

var (
	Scheduler = gocron.NewScheduler(time.UTC)
)
