package main

import (
	"github.com/infrawatch/sg-core-refactor/pkg/bus"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/infrawatch/sg-core-refactor/pkg/handler"
)

type collectdMetricsHandler struct {
	bus bus.MetricBus
}

func (c *collectdMetricsHandler) Handle(msg []byte) (data.Metric, error) {
	return data.Metric{
		Name: "collectd",
		Labels: map[string]string{
			"host": "localhost",
		},
		Type:  data.GAUGE,
		Value: 123421.2,
	}, nil
}

//New create new collectdMetricsHandler object
func New() handler.MetricHandler {
	return &collectdMetricsHandler{}
}
