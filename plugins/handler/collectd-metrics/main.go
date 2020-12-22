package main

import (
	"math/rand"

	"github.com/infrawatch/sg-core-refactor/pkg/bus"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/infrawatch/sg-core-refactor/pkg/handler"
)

type collectdMetricsHandler struct {
	bus bus.MetricBus
}

func (c *collectdMetricsHandler) Handle(msg []byte) ([]data.Metric, error) {
	return []data.Metric{{
		Name: "collectd_cpu",
		Labels: map[string]string{
			"host": "localhost",
		},
		Type:  data.GAUGE,
		Value: rand.Float64() * 1000,
	}, {
		Name: "collectd_memory",
		Labels: map[string]string{
			"host": "localhost",
		},
		Type:  data.GAUGE,
		Value: rand.Float64() * 1000,
	}}, nil
}

//New create new collectdMetricsHandler object
func New() handler.MetricHandler {
	return &collectdMetricsHandler{}
}
