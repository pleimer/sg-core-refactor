package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	log "github.com/infrawatch/apputils/logging"
	"github.com/infrawatch/sg-core-refactor/cmd/manager"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
)

func main() {
	configPath := flag.String("config", "/etc/sg-core.conf.yaml", "configuration file path")
	pluginDir := flag.String("pluginDir", "/usr/lib64/sg-core/", "path to plugin binaries")
	//logLevel := flag.String("logLevel", "ERROR", "log level")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS]\n\nAvailable options:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	logger, err := log.NewLogger(log.DEBUG, "console")
	if err != nil {
		fmt.Printf("failed initializing logger: %s", err)
	}
	logger.Timestamp = true

	file, err := os.Open(*configPath)
	if err != nil {
		logger.Metadata(log.Metadata{"error": err})
		logger.Error("failed opening config file")
		return
	}

	err = parseConfig(file)
	if err != nil {
		logger.Metadata(log.Metadata{"error": err})
		logger.Error("failed parsing config file")
		return
	}

	manager.SetLogger(logger)
	manager.SetPluginDir(*pluginDir)

	for _, tConfig := range config.Transports {
		err = manager.InitTransport(tConfig.Name, tConfig.Mode, tConfig.Config)
		if err != nil {
			logger.Metadata(log.Metadata{"transport": tConfig.Name, "error": err})
			logger.Warn("failed configuring transport")
			continue
		}
		err = manager.SetTransportHandlers(tConfig.Name, tConfig.Handlers)
		if err != nil {
			logger.Metadata(log.Metadata{"transport": tConfig.Name, "error": err})
			logger.Error("failed loading handlers for transport")
			continue
		}
		logger.Metadata(log.Metadata{"transport": tConfig.Name})
		logger.Info("loaded transport")
	}

	wg := new(sync.WaitGroup)
	manager.RunTransports(wg)

	wg.Wait()
}

func run(ctx context.Context) {
	metricChannel := make(chan data.Metric)
	eventChannel := make(chan data.Metric)

	for {
		select {
		case <-eventChannel:
			fmt.Printf("Recieved event")
		case <-metricChannel:
			fmt.Printf("Recieved metric")
		case <-ctx.Done():
			fmt.Printf("Exiting")
		}
	}
}
