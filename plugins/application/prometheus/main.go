package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/infrawatch/apputils/logging"
	"github.com/infrawatch/sg-core-refactor/pkg/application"
	"github.com/infrawatch/sg-core-refactor/pkg/concurrent"
	"github.com/infrawatch/sg-core-refactor/pkg/config"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type configT struct {
	Host          string
	Port          int
	MetricTimeout int
}

// used to expire stale metrics
type metricExpiry struct {
	sync.RWMutex
	lastArrival time.Time
	interval    float64
	delete      func()
}

func (me *metricExpiry) keepAlive() {
	me.Lock()
	defer me.Unlock()
	me.lastArrival = time.Now()
}

func (me *metricExpiry) Expired() bool {
	me.RLock()
	defer me.RUnlock()
	return (time.Since(me.lastArrival).Seconds() >= me.interval)
}

func (me *metricExpiry) Delete() {
	me.Lock()
	defer me.Unlock()
	me.delete()
}

//Prometheus plugin for interfacing with Prometheus
type Prometheus struct {
	configuration configT
	logger        *logging.Logger
	descriptions  *concurrent.Map
	metrics       *concurrent.Map
	labelKeys     []string
	expiry        *expiryProc
	expirys       map[string]*metricExpiry
}

//New constructor
func New(l *logging.Logger) application.Application {
	return &Prometheus{
		configuration: configT{
			Host:          "127.0.0.1",
			Port:          3000,
			MetricTimeout: 20,
		},
		logger:       l,
		descriptions: concurrent.NewMap(),
		metrics:      concurrent.NewMap(),
		expiry:       newExpiryProc(),
		labelKeys:    []string{},
		expirys:      map[string]*metricExpiry{},
	}
}

//Describe implements prometheus.Collector
func (p *Prometheus) Describe(ch chan<- *prometheus.Desc) {
	for desc := range p.descriptions.Iter() {
		ch <- desc.Value.(*prometheus.Desc)
	}
}

//Collect implements prometheus.Collector
func (p *Prometheus) Collect(ch chan<- prometheus.Metric) {
	errs := []error{}
	for item := range p.metrics.Iter() {
		metric := item.Value.(data.Metric)
		labelValues := make([]string, 0, len(p.labelKeys))
		for _, key := range p.labelKeys {
			labelValues = append(labelValues, metric.Labels[key]) //TODO: optimize this
		}
		desc, _ := p.descriptions.Get(metric.Name)
		pMetric, err := prometheus.NewConstMetric(desc.(*prometheus.Desc), metricTypeToPromValueType(metric.Type), metric.Value, labelValues...)
		if err != nil {
			errs = append(errs, err)
		}
		if !metric.Time.IsZero() {
			ch <- prometheus.NewMetricWithTimestamp(metric.Time, pMetric)
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
func (p *Prometheus) Run(ctx context.Context, wg *sync.WaitGroup, eChan chan data.Event, mChan chan []data.Metric) {
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

	//run exporter for prometheus to scrape
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

	//run metric expiry process
	go p.expiry.run(ctx)

	for {
		select {
		case <-ctx.Done():
			goto done
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
done:
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
	if _, found := p.descriptions.Get(name); !found {
		p.labelKeys = make([]string, 0, len(labels))
		for k := range labels {
			p.labelKeys = append(p.labelKeys, k)
		}

		p.descriptions.Set(name, prometheus.NewDesc(name, description, p.labelKeys, nil))
	}
}

func (p *Prometheus) updateMetrics(metric data.Metric) {
	if _, found := p.expirys[metric.Name]; !found {
		exp := metricExpiry{
			interval: 20,
			delete: func() {
				p.metrics.Delete(metric.Name)
				p.descriptions.Delete(metric.Name)
				delete(p.expirys, metric.Name)
			},
		}
		p.expirys[metric.Name] = &exp
		p.expiry.register(&exp)
	}
	p.metrics.Set(metric.Name, metric)
	p.expirys[metric.Name].keepAlive()
}

// helper functions

func metricTypeToPromValueType(mType data.MetricType) prometheus.ValueType {
	return map[data.MetricType]prometheus.ValueType{
		data.COUNTER: prometheus.CounterValue,
		data.DERIVE:  prometheus.CounterValue,
		data.GAUGE:   prometheus.GaugeValue,
		data.UNTYPED: prometheus.UntypedValue,
	}[mType]
}
