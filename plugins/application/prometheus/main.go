package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/infrawatch/apputils/logging"
	"github.com/infrawatch/sg-core-refactor/pkg/application"
	"github.com/infrawatch/sg-core-refactor/pkg/config"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type configT struct {
	Host string
	Port int
}

//Prometheus plugin for interfacing with Prometheus
type Prometheus struct {
	configuration configT
	logger        *logging.Logger
	descriptions  map[string]*prometheus.Desc
	metrics       map[string]data.Metric
}

//New constructor
func New(l *logging.Logger) application.Application {
	return &Prometheus{
		configuration: configT{
			Host: "127.0.0.1",
			Port: 3000,
		},
		logger:       l,
		descriptions: map[string]*prometheus.Desc{},
		metrics:      map[string]data.Metric{},
	}
}

//Describe implements prometheus.Collector
func (p *Prometheus) Describe(ch chan<- *prometheus.Desc) {

}

//Collect implements prometheus.Collector
func (p *Prometheus) Collect(ch chan<- prometheus.Metric) {
	p.logger.Info("prometheus attempted a scrape")
	errs := []error{}
	for _, metric := range p.metrics {
		labelValues := []string{}
		for _, val := range metric.Labels { //TODO: optimize this
			labelValues = append(labelValues, val)
		}
		pMetric, err := prometheus.NewConstMetric(p.descriptions[metric.Name], metricTypeToPromValueType(metric.Type), metric.Value, labelValues...)
		if err != nil {
			errs = append(errs, err)
		}
		if metric.Time != nil {
			ch <- prometheus.NewMetricWithTimestamp(*metric.Time, pMetric)
			continue
		}
		ch <- pMetric
	}

	for _, e := range errs {
		p.logger.Metadata(logging.Metadata{"error": e})
		p.logger.Error("prometheus failed scrapping metric")
	}
}

//Run run scrape endpoint
func (p *Prometheus) Run(wg *sync.WaitGroup, eChan chan data.Event, mChan chan []data.Metric) {
	defer wg.Done()
	registry := prometheus.NewRegistry()

	//Set up Metric Exporter
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
                                <head><title>Prometheus Exporter</title></head>
                                <body>cacheutil
                                <h1>Prometheus Exporter</h1>
                                <p><a href='/metrics'>Metrics</a></p>
                                </body>
								</html>`))
		if err != nil {
			p.logger.Metadata(logging.Metadata{"error": err})
			p.logger.Error("HTTP error")
		}
	})

	registry.MustRegister(p)

	//run exporter fro prometheus to scrape
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		metricsURL := fmt.Sprintf("%s:%d", p.configuration.Host, p.configuration.Port)
		p.logger.Info(fmt.Sprintf("Metric server at : %s", metricsURL))

		err := http.ListenAndServe(metricsURL, handler)
		if err != nil {
			p.logger.Metadata(logging.Metadata{"error": err})
			p.logger.Error("Metric server failed")
		}
	}(wg)

	for {
		select {
		case ev := <-eChan:
			fmt.Printf("Prometheus received event with message: %s\n", ev.Message)
		case metrics := <-mChan:
			// update descriptions
			for _, m := range metrics {
				fmt.Printf("Prometheus received metric: %v\n", m)
				p.updateDescs(m.Name, "", m.Labels)
				p.updateMetrics(m)
			}
		}
	}
}

//Config implements application.Application
func (p *Prometheus) Config(c []byte) error {
	err := config.ParseConfig(bytes.NewReader(c), &p.configuration)
	if err != nil {
		return err
	}
	return nil
}

// updateDescs update prometheus descriptions
func (p *Prometheus) updateDescs(name string, description string, labels map[string]string) {
	if _, found := p.descriptions[name]; !found {
		keys := make([]string, 0, len(labels))
		for k := range labels {
			keys = append(keys, k)
		}
		p.descriptions[name] = prometheus.NewDesc(name, description, keys, nil)
	}
}

func (p *Prometheus) updateMetrics(metric data.Metric) {
	p.metrics[metric.Name] = metric
}

// helper functions

func metricTypeToPromValueType(mType data.MetricType) prometheus.ValueType {
	return map[data.MetricType]prometheus.ValueType{
		data.COUNTER: prometheus.CounterValue,
		data.GAUGE:   prometheus.GaugeValue,
		data.UNTYPED: prometheus.UntypedValue,
	}[mType]
}
