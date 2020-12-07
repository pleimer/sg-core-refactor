package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/infrawatch/sg-core-refactor/cmd/manager"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
)

func main() {
	configPath := flag.String("config", "/etc/sg-core.conf.yaml", "configuration file path")
	pluginDir := flag.String("pluginDir", "/usr/lib64/sg-core/", "path to plugin binaries")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS]\n\nAvailable options:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		fmt.Printf("failed opening config file: %s\n", err)
		return
	}

	err = parseConfig(file)
	if err != nil {
		fmt.Printf("failed parsing config file: %s\n", err)
		return
	}

	manager.SetPluginDir(*pluginDir)

	for _, tConfig := range config.Transports {
		err = manager.SetTransport(tConfig.Name, tConfig.Mode, tConfig.Handlers, tConfig.Config)
		if err != nil {
			fmt.Printf("failed configuring %s plugin: %s\n", tConfig.Name, err)
			continue
		}
	}
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
