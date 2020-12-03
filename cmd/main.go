package main

import (
	"flag"
	"fmt"
	"os"
	"plugin"

	"github.com/infrawatch/sg-core-refactor/pkg/transport"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	configPath := flag.String("config", "/etc/sg-core.conf.yaml", "configuration file path")
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

	fmt.Println(config.Plugins.Transports[0].Name)

	//register plugins

	p, err := plugin.Open("/home/pleimer/go/src/github.com/infrawatch/sg-core-refactor/bin/socket.so")
	if err != nil {
		panic(err)
	}
	s, err := p.Lookup("New")
	if err != nil {
		panic(err)
	}

	newSocket := s.(func() transport.Transport)
	s1 := newSocket()
	s2 := newSocket()
	s1.Run()
	s2.Run()
}
