package main

import (
	"fmt"
	"time"

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

	var metric *data.Metric
	metrics := []data.Metric{}
	for index, cdmetric := range *cdmetrics {
		metric, err = createMetric(&cdmetric, index)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, *metric)
	}

	return metrics, nil
}

func createMetric(cdmetric *collectd.Metric, index int) (*data.Metric, error) {
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

	metricName := genMetricName(cdmetric, index)

	//TODO: use cdmetric time in metric
	fmt.Printf("Time is %s\n", cdmetric.Time.Time())
	return &data.Metric{
		Name:  metricName,
		Type:  mt,
		Value: cdmetric.Values[index],
		Time:  time.Now(),
		Labels: map[string]string{
			"host":            cdmetric.Host,
			"plugin_instance": cdmetric.PluginInstance,
			"type_instance":   cdmetric.TypeInstance,
		},
	}, nil
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
