package main

import (
	"fmt"
	"plugin"

	"github.com/infrawatch/sg-core-refactor/pkg/transport"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

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
