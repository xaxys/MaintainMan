package service

import (
	"github.com/olebedev/emitter"
	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/logger"
)

var (
	Bus *emitter.Emitter
)

func init() {
	size := config.AppConfig.GetUint("bus_buffer")
	Bus = emitter.New(size)
	Bus.On("*", func(e *emitter.Event) {
		logger.Logger.Infof("Event Detected: %s", e.OriginalTopic)
		logger.Logger.Debugf("Event Data: %#v", *e)
	})
}
