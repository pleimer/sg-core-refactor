package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/infrawatch/sg-core-refactor/cmd/manager"
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
		fmt.Printf("failed opening config file: %s\n", err.Error())
		return
	}

	err = parseConfig(file)
	if err != nil {
		fmt.Printf("failed parsing config file: %s\n", err.Error())
		return
	}

	err = manager.LoadBinaries(*pluginDir)
	if err != nil {
		fmt.Printf("failed loading plugin binaries: %s\n", err.Error())
	}

	for _, pluginConfig := range config.Plugins {
		err = manager.SetPluginConfig(pluginConfig.Name, pluginConfig.Config)
		if err != nil {
			fmt.Printf("failed configuring %s plugin: %s\n", pluginConfig.Name, err)
			continue
		}
	}
}
