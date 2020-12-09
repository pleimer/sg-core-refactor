package main

import (
	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/infrawatch/sg-core-refactor/pkg/handler"
)

type collectdMetricsHandler struct{}

func (c *collectdMetricsHandler) Handle(msg []byte) (interface{}, error) {
	return data.Event{Message: string(msg)}, nil
}

func (c *collectdMetricsHandler) Type() data.Type {
	return data.METRIC
}

//New create new collectdMetricsHandler object
func New() handler.Handler {
	return &collectdMetricsHandler{}
}
