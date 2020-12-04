package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/infrawatch/sg-core-refactor/pkg/transport"
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
		fmt.Printf("failed opening file: %s\n", err.Error())
		return
	}

	err = parseConfig(file)
	if err != nil {
		fmt.Printf("failed parsing config file: %s\n", err.Error())
		return
	}

	// load config binaries
	sharedObjs := map[string]string{}
	err = filepath.Walk(*pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".so" {
			return nil
		}

		baseName := strings.Split(info.Name(), ".")
		sharedObjs[baseName[len(baseName)-2]] = path
		return nil
	})

	if err != nil {
		fmt.Printf("failed loading plugin binaries: %s\n", err.Error())
	}

	plugins := []interface{}{}
	for _, pluginConfig := range config.Plugins {
		if _, ok := sharedObjs[pluginConfig.Name]; !ok {
			fmt.Printf("[WARNING] Could not load plugin '%s': binary does not exist in %s\n", pluginConfig.Name, *pluginDir)
			continue
		}

		p, err := plugin.Open(sharedObjs[pluginConfig.Name])
		if err != nil {
			fmt.Printf("[WARNING] Could not load plugin %s: %s\n", pluginConfig.Name, err)
			continue
		}

		n, err := p.Lookup("New")
		if err != nil {
			fmt.Printf("[WARNING] Could not load plugin 'New' function %s: %s\n", pluginConfig.Name, err)
			continue
		}

		o := n.(func() transport.Transport)()
		plugins = append(plugins, o)
		o.Config(pluginConfig.Config)
	}
}
