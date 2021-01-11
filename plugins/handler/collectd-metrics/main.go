package main

import (
	"fmt"
	"time"

	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/infrawatch/sg-core-refactor/pkg/handler"
	"github.com/infrawatch/sg-core-refactor/plugins/handler/collectd-metrics/pkg/collectd"
)

type collectdMetricsHandler struct {
	totalMetricsReceived uint64
	totalDecodeErrors    uint64
}

func (c *collectdMetricsHandler) Handle(blob []byte) []data.Metric {

	var err error
	var cdmetrics *[]collectd.Metric

	cdmetrics, err = collectd.ParseInputByte(blob)
	if err != nil {
		c.totalDecodeErrors++
		return nil
	}

	var ms []data.Metric
	metrics := []data.Metric{}
	for _, cdmetric := range *cdmetrics {
		ms, err = c.createMetrics(&cdmetric)
		if err != nil {
			c.totalDecodeErrors++
			continue
		}

		metrics = append(metrics, ms...)
	}

	return metrics
}

func (c *collectdMetricsHandler) createMetrics(cdmetric *collectd.Metric) ([]data.Metric, error) {
	if cdmetric.Host == "" {
		return nil, fmt.Errorf("missing host: %v ", cdmetric)
	}

	pluginInstance := cdmetric.PluginInstance
	if pluginInstance == "" {
		pluginInstance = "base"
	}
	typeInstance := cdmetric.TypeInstance
	if typeInstance == "" {
		typeInstance = "base"
	}

	var mt data.MetricType
	var err error

	var metrics []data.Metric
	for index := range cdmetric.Dsnames {
		if mt, err = mt.FromString(cdmetric.Dstypes[index]); err != nil {
			return nil, err
		}

		metrics = append(metrics,
			data.Metric{
				Name:  genMetricName(cdmetric, index),
				Type:  mt,
				Value: cdmetric.Values[index],
				Time:  cdmetric.Time.Time(),
				Labels: map[string]string{
					"host":            cdmetric.Host,
					"plugin_instance": cdmetric.PluginInstance,
					"type_instance":   cdmetric.TypeInstance,
				}})
	}
	metrics = append(metrics, []data.Metric{{
		Name:  "sg_total_metric_rcv_count",
		Type:  data.COUNTER,
		Value: float64(c.totalMetricsReceived),
		Time:  time.Now(),
		Labels: map[string]string{
			"source": "SG",
		},
	}, {
		Name:  "sg_total_metric_decode_error_count",
		Type:  data.COUNTER,
		Value: float64(c.totalDecodeErrors),
		Time:  time.Now(),
		Labels: map[string]string{
			"source": "SG",
		},
	},
	}...)
	c.totalMetricsReceived++
	return metrics, nil
}

func genMetricName(cdmetric *collectd.Metric, index int) (name string) {

	name = "collectd_" + cdmetric.Plugin + "_" + cdmetric.Type
	if cdmetric.Type == cdmetric.Plugin {
		name = "collectd_" + cdmetric.Plugin
	}

	if dsname := cdmetric.Dsnames[index]; dsname != "value" {
		name += "_" + dsname
	}

	switch cdmetric.Dstypes[index] {
	case "counter", "derive":
		name += "_total"
	}

	return
}

//New create new collectdMetricsHandler object
func New() handler.MetricHandler {
	return &collectdMetricsHandler{}
}
