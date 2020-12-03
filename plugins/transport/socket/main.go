package main

import (
	"fmt"

	"github.com/infrawatch/sg-core-refactor/pkg/transport"
)

//Socket basic struct
type Socket struct{}

//Run run the socket!
func (s *Socket) Run() {
	fmt.Println("Running the socket!")
}

func New() transport.Transport {
	return &Socket{}
}
