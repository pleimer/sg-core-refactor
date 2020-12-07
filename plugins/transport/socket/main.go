package main

import (
	"fmt"

	"github.com/infrawatch/sg-core-refactor/pkg/transport"
)

// Config holds socket plugin configuration
var Config struct {
	Param string `validator:"required"`
}

//Socket basic struct
type Socket struct{}

//Run implements type Transport
func (s *Socket) Run() {
	fmt.Println("Running the socket!")
}

//Config implements type Transport
func (s *Socket) Config(c interface{}) error {
	fmt.Println(c)
	return nil
}

//New create new socket transport
func New() transport.Transport {
	return &Socket{}
}
