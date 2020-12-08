package main

import (
	"fmt"
	"sync"

	"github.com/infrawatch/sg-core-refactor/pkg/transport"
)

// Config holds socket plugin configuration
var Config struct {
	Param string `validator:"required"`
}

//Socket basic struct
type Socket struct{}

//Run implements type Transport
func (s *Socket) Run(wg *sync.WaitGroup, t chan []byte) {
	defer wg.Done()
	t <- []byte("hello!")
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
