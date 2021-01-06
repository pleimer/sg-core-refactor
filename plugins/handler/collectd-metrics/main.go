package main

import (
	"fmt"

	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/infrawatch/sg-core-refactor/pkg/handler"
	"github.com/infrawatch/sg-core-refactor/plugins/handler/collectd-metrics/pkg/collectd"
)

type collectdMetricsHandler struct {
}

func (c *collectdMetricsHandler) Handle(blob []byte) ([]data.Metric, error) {

	cdmetrics, err := collectd.ParseInputByte(blob)
	if err != nil {
		return nil, err
	}

	var ms []data.Metric
	metrics := []data.Metric{}
	for index, cdmetric := range *cdmetrics {
		ms, err = createMetrics(&cdmetric, index)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, ms...)
	}

	return metrics, nil
}

func createMetrics(cdmetric *collectd.Metric, index int) ([]data.Metric, error) {
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
	if mt, err = mt.FromString(cdmetric.Dstypes[index]); err != nil {
		return nil, err
	}

	var metrics []data.Metric
	for index := range cdmetric.Dsnames {
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
